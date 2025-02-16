package task

import (
	"time"
)

type ProcessWatcher struct {
	p *Process

	t       Task
	padding int

	state chan ProcessState

	exitCh chan struct{}
	result bool
}

func NewProcessWatcher(t Task, padding int) *ProcessWatcher {
	stateCh := make(chan ProcessState, 1)
	stateCh <- ProcessStateIdle
	return &ProcessWatcher{
		t:       t,
		padding: padding,
		state:   stateCh,
		exitCh:  make(chan struct{}),
	}
}

func (w *ProcessWatcher) Start() {
	s := <-w.state
	if s != ProcessStateIdle {
		return
	}

	go w.watch()

	w.state <- ProcessStateStarting
	w._start()
}

func (w *ProcessWatcher) Stop() {
	s := <-w.state
	// nolint:exhaustive
	switch s {
	case ProcessStateRestarting, ProcessStateRunning:
	case ProcessStateIdle:
		w.state <- ProcessStateExited
		close(w.exitCh)
		return
	default:
		w.state <- s
		return
	}

	w.state <- ProcessStateStopping
	w.p.Stop()
}

func (w *ProcessWatcher) Kill() {
	s := <-w.state
	// nolint:exhaustive
	switch s {
	case ProcessStateRestarting, ProcessStateRunning, ProcessStateStopping:
	case ProcessStateIdle:
		w.state <- ProcessStateExited
		close(w.exitCh)
		return
	default:
		w.state <- s
		return
	}

	w.state <- ProcessStateKilling
	w.p.Kill()
}

func (w *ProcessWatcher) _start() {
	s := <-w.state
	// nolint:exhaustive
	switch s {
	case ProcessStateStarting, ProcessStateRestarting:
	default:
		w.state <- s
		return
	}

	proc := NewProcess(w.t, w.padding)
	proc.Start()
	w.p = proc

	w.state <- ProcessStateRunning
	go func() {
		result := proc.Wait()
		s := <-w.state
		if s == ProcessStateRestarting {
			w.state <- ProcessStateRestarting
			go w._start()
			return
		}

		w.result = result
		w.state <- ProcessStateExited
		close(w.exitCh)
	}()
}

func (w *ProcessWatcher) watch() {
	wt := Watch(w.t.WatchFiles)

	t := time.NewTicker(100 * time.Millisecond)
	defer t.Stop()

	for {

		select {
		case <-w.exitCh:
			return
		case <-t.C:
		}

		if !wt.Changed() {
			continue
		}

		w._restart()
	}
}

func (w *ProcessWatcher) _restart() {
	s := <-w.state
	if s != ProcessStateRunning {
		w.state <- s
		return
	}

	w.p.logAction("Change detected, restarting...")

	w.state <- ProcessStateRestarting
	w.p.Stop()
}

func (w *ProcessWatcher) Done() bool {
	s := <-w.state
	w.state <- s
	return s == ProcessStateExited
}

func (w *ProcessWatcher) Wait() bool {
	<-w.exitCh
	return w.result
}
