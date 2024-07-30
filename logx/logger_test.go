package logx_test

import (
	"testing"

	"github.com/zeiss/pkg/logx"

	"github.com/stretchr/testify/assert"
)

func TestFacade(t *testing.T) {
	logx.Printf("test %q", "print")

	logx.Debugf("test %q", "debug")
	logx.Infof("test %q", "info")
	logx.Warnf("test %q", "warn")
	assert.Panics(t, func() { logx.Panicf("test %q", "panic") })
	logx.Errorf("test %q", "error")

	logx.Debugw("test", "some", "debug")
	logx.Infow("test", "some", "info")
	logx.Warnw("test", "some", "warn")
	assert.Panics(t, func() { logx.Panicw("test", "some", "panic") })
	logx.Errorw("test", "some", "error")
}
