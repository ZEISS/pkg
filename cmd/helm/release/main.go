package main

import (
	"bytes"
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"text/template"
	"time"

	"github.com/zeiss/pkg/cast"
	"github.com/zeiss/pkg/logx"

	"github.com/Songmu/retry"
	"github.com/google/go-github/v90/github"
	"github.com/spf13/pflag"
	"helm.sh/helm/v3/pkg/chart/loader"
)

type flags struct {
	PackagePath string
	Repo        string
	Owner       string
	Token       string
	ReleaseName string
	Commit      string
}

func main() {
	log.SetFlags(0)
	log.SetOutput(os.Stderr)

	_, err := logx.RedirectStdLog(logx.LogSink)
	if err != nil {
		panic(err)
	}

	f := &flags{}

	pflag.StringVar(&f.PackagePath, "package-path", f.PackagePath, "package path")
	pflag.StringVar(&f.Repo, "repo", f.Repo, "repo")
	pflag.StringVar(&f.Owner, "owner", f.Owner, "owner")
	pflag.StringVar(&f.Token, "token", f.Token, "token")
	pflag.StringVar(&f.Commit, "commit", f.Commit, "commit")
	pflag.StringVar(&f.ReleaseName, "release-name", "{{ .Name }}-{{ .Version }}", "release name")
	pflag.Parse()

	client, err := github.NewClient(github.WithAuthToken(f.Token))
	if err != nil {
		panic(err)
	}

	chartPackages, err := filepath.Glob(f.PackagePath + "/*.tgz")
	if err != nil {
		panic(err)
	}

	for _, pkg := range chartPackages {
		a, err := os.Open(pkg)
		if err != nil {
			panic(err)
		}
		defer a.Close()

		chart, err := loader.LoadArchive(a)
		if err != nil {
			panic(err)
		}

		tmpl, err := template.New("gotpl").Parse(f.ReleaseName)
		if err != nil {
			panic(err)
		}

		var buffer bytes.Buffer
		if err := tmpl.Execute(&buffer, chart.Metadata); err != nil {
			panic(err)
		}

		releaseName := buffer.String()

		req := github.CreateReleaseRequest{
			Name:            cast.Ptr(releaseName),
			Body:            cast.Ptr(chart.Metadata.Description),
			TagName:         releaseName,
			TargetCommitish: cast.Ptr(f.Commit),
			MakeLatest:      cast.Ptr(strconv.FormatBool(true)),
		}

		release, _, err := client.Repositories.CreateRelease(context.TODO(), f.Owner, f.Repo, req)
		if err != nil {
			panic(err)
		}

		opts := &github.UploadOptions{
			// Use base name by default
			Name: filepath.Base(pkg),
		}

		_, err = a.Seek(0, io.SeekStart)
		if err != nil {
			panic(err)
		}

		err = retry.Retry(3, 3*time.Second, func() error { //nolint: revive
			if _, _, err = client.Repositories.UploadReleaseAsset(context.TODO(), f.Owner, f.Repo, release.ID, opts, a); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			panic(err)
		}
	}
}
