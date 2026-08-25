package main

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/zeiss/pkg/cast"
	"github.com/zeiss/pkg/filex"
	"github.com/zeiss/pkg/logx"

	"github.com/google/go-github/v90/github"
	"github.com/spf13/pflag"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/provenance"
	"helm.sh/helm/v3/pkg/repo"
)

type flags struct {
	Repo  string
	Index string
}

func main() {
	log.SetFlags(0)
	log.SetOutput(os.Stderr)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := logx.RedirectStdLog(logx.LogSink)
	if err != nil {
		panic(err)
	}

	f := &flags{}

	pflag.StringVar(&f.Index, "index", f.Index, "index (default: index.yaml)")
	pflag.StringVar(&f.Repo, "repo", f.Repo, "repo")
	pflag.Parse()

	client, err := github.NewClient()
	if err != nil {
		panic(err)
	}

	indexFile := repo.NewIndexFile()

	opts := &github.ListOptions{
		PerPage: 100,
	}

	slices := strings.Split(f.Repo, "/")
	owner := slices[0]
	repo := slices[1]

	releases, _, err := client.Repositories.ListReleases(ctx, owner, repo, opts)
	if err != nil {
		panic(err)
	}

	for _, release := range releases {
		for _, asset := range release.Assets {
			ext := filepath.Ext(cast.Value(asset.Name))
			name := cast.Value(asset.Name)

			u, err := url.Parse(cast.Value(asset.BrowserDownloadURL))
			if err != nil {
				panic(err)
			}

			if strings.HasSuffix(ext, ".tgz") {
				req := &http.Request{
					Method: "GET",
					URL:    u,
				}

				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					panic(err)
				}

				buf := &bytes.Buffer{}
				if _, err := io.Copy(buf, resp.Body); err != nil {
					log.Printf("Failed to send out response: %v", err)
				}
				defer resp.Body.Close()

				hashReader := bytes.NewReader(buf.Bytes())
				hash, err := provenance.Digest(hashReader)
				if err != nil {
					panic(err)
				}

				bodyReader := bytes.NewReader(buf.Bytes())
				chart, err := loader.LoadArchive(bodyReader)
				if err != nil {
					panic(err)
				}

				url := strings.TrimRight(u.String(), name)
				err = indexFile.MustAdd(chart.Metadata, name, url, hash)
				if err != nil {
					panic(err)
				}
			}
		}
	}

	base := filepath.Dir(f.Index)

	err = filex.MkdirAll(base, 0o755)
	if err != nil {
		panic(err)
	}

	indexFile.SortEntries()
	indexFile.Generated = time.Now()

	err = indexFile.WriteFile(f.Index, 0o644)
	if err != nil {
		panic(err)
	}
}
