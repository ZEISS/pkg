package main

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/zeiss/pkg/filex"
	"github.com/zeiss/pkg/logx"

	"github.com/spf13/pflag"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/registry"
)

type flags struct {
	PackagePath string
	Paths       []string
}

func main() {
	log.SetFlags(0)
	log.SetOutput(os.Stderr)

	_, err := logx.RedirectStdLog(logx.LogSink)
	if err != nil {
		panic(err)
	}

	f := &flags{}

	pflag.StringSliceVar(&f.Paths, "path", f.Paths, "chart path")
	pflag.StringVar(&f.PackagePath, "package-path", f.PackagePath, "package path")
	pflag.Parse()

	helmClient := action.NewPackage()
	helmClient.DependencyUpdate = true
	helmClient.Destination = f.PackagePath

	settings := cli.New()
	getters := getter.All(settings)
	registryClient, err := registry.NewClient()
	if err != nil {
		panic(err)
	}

	err = filex.MkdirAll(f.PackagePath, 0o755)
	if err != nil {
		panic(err)
	}

	for _, chartPath := range f.Paths {
		path, err := filepath.Abs(chartPath)
		if err != nil {
			panic(err)
		}
		if _, err := os.Stat(chartPath); err != nil {
			panic(err)
		}

		downloadManager := &downloader.Manager{
			Out:              io.Discard,
			ChartPath:        path,
			Keyring:          helmClient.Keyring,
			Getters:          getters,
			Debug:            settings.Debug,
			RepositoryConfig: settings.RepositoryConfig,
			RepositoryCache:  settings.RepositoryCache,
			RegistryClient:   registryClient,
		}
		if err := downloadManager.Build(); err != nil {
			panic(err)
		}

		packageRun, err := helmClient.Run(path, nil)
		if err != nil {
			panic(err)
		}

		logx.Printf("Successfully packaged chart in %s\n", packageRun)
	}
}
