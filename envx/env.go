package envx

import (
	"strings"
)

// Env is the environment variables.
type Env map[string]string

// Copy returns a copy of the environment.
func (e Env) Copy() Env {
	out := Env{}
	for k, v := range e {
		out[k] = v
	}

	return out
}

// Strings returns the current environment as a list of strings, suitable for
// os executions.
func (e Env) Strings() []string {
	result := make([]string, 0, len(e))
	for k, v := range e {
		result = append(result, k+"="+v)
	}

	return result
}

// ToEnv converts the environment to a list of strings.
func ToEnv(env []string) Env {
	ee := Env{}

	for _, e := range env {
		k, v, ok := strings.Cut(e, "=")
		if !ok || k == "" {
			continue
		}

		ee[k] = v
	}

	return ee
}
