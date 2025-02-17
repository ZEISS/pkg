package task

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"golang.org/x/sync/errgroup"
)

type Runner struct {
	procs []ProcessRunner

	doneCh chan struct{}
	stopCh chan struct{}

	stopOnce sync.Once
}

func NewRunner(tasks []Task) *Runner {
	procs := make([]ProcessRunner, 0, len(tasks))
	logName := "RUNPROC"
	maxLen := len(logName)
	for _, t := range tasks {
		if len(t.Name) > maxLen {
			maxLen = len(t.Name)
		}
	}

	if maxLen > len(logName) {
		logName += strings.Repeat(" ", maxLen-len(logName))
	}
	log.SetOutput(&prefixer{
		out:    log.Default().Writer(),
		prefix: color.New(color.BgRed, color.FgWhite).Sprint(logName),
	})
	for _, t := range tasks {
		if len(t.WatchFiles) > 0 {
			procs = append(procs, NewProcessWatcher(t, maxLen))
			continue
		}

		procs = append(procs, NewProcess(t, maxLen))
	}
	return &Runner{
		procs:  procs,
		stopCh: make(chan struct{}),
		doneCh: make(chan struct{}),
	}
}

func (r *Runner) Run(ctx context.Context) error {
	g, _ := errgroup.WithContext(ctx)
	result := make(chan bool, len(r.procs))

	for _, proc := range r.procs {
		g.Go(func() error {
			proc.Start()
			result <- proc.Wait()
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func (r *Runner) Stop() {
	r.stopOnce.Do(r._stop)
}

func (r *Runner) _stop() {
	close(r.stopCh)

	for _, proc := range r.procs {
		go proc.Stop()
	}

	t := time.NewTimer(time.Second)
	defer t.Stop()
	select {
	case <-r.doneCh:
	case <-t.C:
		for _, proc := range r.procs {
			if !proc.Done() {
				proc.Kill()
			}
		}
	}

	<-r.doneCh
}
