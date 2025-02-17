package task

type ProcessRunner interface {
	Start()
	Stop()
	Kill()
	Done() bool
	Wait() bool
}
