package main

import (
	"bytes"
	"context"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/zeiss/pkg/cmd/runproc/task"
)

type config struct {
	File  string
	Local string
}

var cfg = &config{}

var rootCmd = &cobra.Command{
	Use:   "runproc",
	Short: `runproc is a simple process runner.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRoot(cmd.Context())
	},
}

func init() {
	rootCmd.Flags().StringVarP(&cfg.File, "file", "f", cfg.File, "Procfile to run.")
	rootCmd.Flags().StringVarP(&cfg.Local, "local", "l", cfg.Local, "Local Procfile to append.")

	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true
}

func runRoot(ctx context.Context) error {
	data, err := os.ReadFile(cfg.File)
	if err != nil {
		log.Fatal(err)
	}

	envData, err := os.ReadFile(cfg.Local)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(data)
	buf.WriteString("\n")
	buf.Write(envData)

	tasks, err := task.Parse(buf)
	if err != nil {
		log.Fatalln(err)
	}
	log.SetFlags(log.Lshortfile)

	run := task.NewRunner(tasks)

	err = run.Run(ctx)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	log.SetFlags(0)
	log.SetOutput(os.Stderr)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
