package main

import (
	"bytes"
	"context"
	"log"
	"os"

	"github.com/spf13/cobra"
)

type config struct {
	file  string
	local string
}

var cfg = &config{}

var rootCmd = &cobra.Command{
	Use:   "nctl",
	Short: "nctl is a tool for managing operator resources",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRoot(cmd.Context())
	},
}

func init() {
	rootCmd.Flags().StringP("file", "f", cfg.file, "Procfile to run.")
	rootCmd.Flags().StringP("local", "l", cfg.local, "Local Procfile to append.")
}

func runRoot(ctx context.Context) error {
	data, err := os.ReadFile(cfg.file)
	if err != nil {
		log.Fatal(err)
	}

	envData, err := os.ReadFile(cfg.local)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(data)
	buf.WriteString("\n")
	buf.Write(envData)

	tasks, err := Parse(buf)
	if err != nil {
		log.Fatalln(err)
	}
	log.SetFlags(log.Lshortfile)

	run := NewRunner(tasks)

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
