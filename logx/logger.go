package logx

import (
	"sync"

	"go.uber.org/zap"
)

// LogSink is the logger sink.
var LogSink Logger

func init() {
	l, err := NewLogSink()
	if err != nil {
		panic(err)
	}

	log := NewLogger(WithLogger(l))

	LogSink = log
}

// NewLogSink returns a new logger sink.
func NewLogSink() (*zap.Logger, error) {
	return zap.NewProduction()
}

// Logger represents a standard logging interface.
type Logger interface {
	// Log a notice statement
	Noticef(format string, v ...interface{})
	// Infof is logging an info statement.
	Infof(format string, v ...interface{})
	// Log a warning statement
	Warnf(format string, v ...interface{})
	// Log a fatal error
	Fatalf(format string, v ...interface{})
	// Log an error
	Errorf(format string, v ...interface{})
	// Log a debug statement
	Debugf(format string, v ...interface{})
	// Log a trace statement
	Tracef(format string, v ...interface{})
	// Panicf is logging a panic statement.
	Panicf(format string, v ...interface{})
	// Printf is logging a printf statement.
	Printf(format string, v ...interface{})
	// Debugw is logging a debug statement with context.
	Debugw(msg string, keysAndValues ...interface{})
	// Infow is logging an info statement with context.
	Infow(msg string, keysAndValues ...interface{})
	// Warnw is logging a warning statement with context.
	Warnw(msg string, keysAndValues ...interface{})
	// Errorw is logging an error statement with context.
	Errorw(msg string, keysAndValues ...interface{})
	// DPanicw is logging a debug panic statement with context.
	DPanicw(msg string, keysAndValues ...interface{})
	// Panicw is logging a panic statement with context.
	Panicw(msg string, keysAndValues ...interface{})
	// Fatalw is logging a fatal statement with context.
	Fatalw(msg string, keysAndValues ...interface{})
}

var _ Logger = (*logger)(nil)

// LogFunc is a bridge between Logger and any third party logger.
type LogFunc func(string, ...interface{})

// Printf is a bridge between Logger and any third party logger.
func (f LogFunc) Printf(msg string, args ...interface{}) { f(msg, args...) }

type logger struct {
	opts *Opts
	sync.RWMutex
}

// Opt is a logger option.
type Opt func(*Opts)

// Opts are the options for the logger.
type Opts struct {
	Logger *zap.Logger
}

// Configure is configuring the logger.
func (o *Opts) Configure(opts ...Opt) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithLogger is setting the logger.
func WithLogger(l *zap.Logger) Opt {
	return func(o *Opts) {
		o.Logger = l
	}
}

// NewLogger is creating a new logger.
func NewLogger(o ...Opt) Logger {
	options := new(Opts)
	options.Configure(o...)

	l := new(logger)
	l.opts = options

	return l
}

// Errorf is logging an error.
func (l *logger) Errorf(format string, v ...interface{}) {
	l.logFunc(func(log *zap.Logger, format string, v ...interface{}) {
		log.Sugar().Errorf(format, v...)
	}, format, v...)
}

// Debugf is logging a debug statement.
func (l *logger) Debugf(format string, v ...interface{}) {
	l.logFunc(func(log *zap.Logger, format string, v ...interface{}) {
		log.Sugar().Debugf(format, v...)
	}, format, v...)
}

// Fatalf is logging a fatal error.
func (l *logger) Fatalf(format string, v ...interface{}) {
	l.logFunc(func(log *zap.Logger, format string, v ...interface{}) {
		log.Sugar().Fatalf(format, v...)
	}, format, v...)
}

// Noticef is logging a notice statement.
func (l *logger) Noticef(format string, v ...interface{}) {
	l.logFunc(func(log *zap.Logger, format string, v ...interface{}) {
		log.Sugar().Infof(format, v...)
	}, format, v...)
}

// Warnf is logging a warning statement.
func (l *logger) Warnf(format string, v ...interface{}) {
	l.logFunc(func(log *zap.Logger, format string, v ...interface{}) {
		log.Sugar().Warnf(format, v...)
	}, format, v...)
}

// Tracef is logging a trace statement.
func (l *logger) Tracef(format string, v ...interface{}) {
	l.logFunc(func(log *zap.Logger, format string, v ...interface{}) {
		log.Sugar().Debugf(format, v...)
	}, format, v...)
}

// Infof is logging an info statement.
func (l *logger) Infof(format string, v ...interface{}) {
	l.logFunc(func(log *zap.Logger, format string, v ...interface{}) {
		log.Sugar().Infof(format, v...)
	}, format, v...)
}

// Panicf is logging a panic statement.
func (l *logger) Panicf(format string, v ...interface{}) {
	l.logFunc(func(log *zap.Logger, format string, v ...interface{}) {
		log.Sugar().Panicf(format, v...)
	}, format, v...)
}

// Printf is logging a printf statement.
func (l *logger) Printf(format string, v ...interface{}) {
	l.logFunc(func(log *zap.Logger, format string, v ...interface{}) {
		log.Sugar().Infof(format, v...)
	}, format, v...)
}

// Debugw is logging a debug statement with context.
func (l *logger) Debugw(msg string, keysAndValues ...interface{}) {
	l.logFunc(func(log *zap.Logger, msg string, keysAndValues ...interface{}) {
		log.Sugar().Debugw(msg, keysAndValues...)
	}, msg, keysAndValues...)
}

// Infow is logging an info statement with context.
func (l *logger) Infow(msg string, keysAndValues ...interface{}) {
	l.logFunc(func(log *zap.Logger, msg string, keysAndValues ...interface{}) {
		log.Sugar().Infow(msg, keysAndValues...)
	}, msg, keysAndValues...)
}

// Warnw is logging a warning statement with context.
func (l *logger) Warnw(msg string, keysAndValues ...interface{}) {
	l.logFunc(func(log *zap.Logger, msg string, keysAndValues ...interface{}) {
		log.Sugar().Warnw(msg, keysAndValues...)
	}, msg, keysAndValues...)
}

// Errorw is logging an error statement with context.
func (l *logger) Errorw(msg string, keysAndValues ...interface{}) {
	l.logFunc(func(log *zap.Logger, msg string, keysAndValues ...interface{}) {
		log.Sugar().Errorw(msg, keysAndValues...)
	}, msg, keysAndValues...)
}

// DPanicw is logging a debug panic statement with context.
func (l *logger) DPanicw(msg string, keysAndValues ...interface{}) {
	l.logFunc(func(log *zap.Logger, msg string, keysAndValues ...interface{}) {
		log.Sugar().DPanicw(msg, keysAndValues...)
	}, msg, keysAndValues...)
}

// Panicw is logging a panic statement with context.
func (l *logger) Panicw(msg string, keysAndValues ...interface{}) {
	l.logFunc(func(log *zap.Logger, msg string, keysAndValues ...interface{}) {
		log.Sugar().Panicw(msg, keysAndValues...)
	}, msg, keysAndValues...)
}

// Fatalw is logging a fatal statement with context.
func (l *logger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.logFunc(func(log *zap.Logger, msg string, keysAndValues ...interface{}) {
		log.Sugar().Fatalw(msg, keysAndValues...)
	}, msg, keysAndValues...)
}

func (l *logger) logFunc(f func(log *zap.Logger, format string, v ...interface{}), format string, args ...interface{}) {
	l.Lock()
	defer l.Unlock()

	if l.opts.Logger == nil {
		return
	}

	f(l.opts.Logger, format, args...)
}

// RedirectStdLog is redirecting the standard logger to the logger.
func RedirectStdLog(l Logger) (func(), error) {
	return zap.RedirectStdLogAt(l.(*logger).opts.Logger, zap.DebugLevel)
}
