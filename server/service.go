package server

import (
	"os"
	"path"
	"sync"
)

// ServiceEnv is a list of environment variables to lookup the service name.
type ServiceEnv []Name

// Name is used to return the service name.
type Name string

// String returns the name as a string.
func (n Name) String() string {
	return string(n)
}

// DefaultEnv is the default environment variables to lookup the service name.
var DefaultEnv = ServiceEnv{Name("SERVICE_NAME")}

func init() {
	Service.lookup(DefaultEnv...)
}

// Service is used to configure the
type service struct {
	name string

	sync.Once
}

// Service is used to configure the service.
var Service = &service{}

// Name returns the service name.
func (s *service) Name() string {
	return s.name
}

// Loopkup is used to lookup the service name.
func (s *service) Lookup(names ...Name) string {
	s.Do(func() {
		s.lookup(names...)
	})

	return s.Name()
}

func (s *service) lookup(names ...Name) {
	for _, name := range names {
		v, ok := os.LookupEnv(name.String())
		if ok {
			s.name = v
			break
		}
	}

	if s.name == "" {
		s.name = path.Base(os.Args[0])
	}
}
