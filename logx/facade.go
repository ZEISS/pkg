package logx

// Printf ...
func Printf(format string, args ...interface{}) {
	LogSink.Infof(format, args...)
}

// Debugf ...
func Debugf(format string, args ...interface{}) {
	LogSink.Debugf(format, args...)
}

// Infof ...
func Infof(format string, args ...interface{}) {
	LogSink.Infof(format, args...)
}

// Errorf ...
func Errorf(format string, args ...interface{}) {
	LogSink.Errorf(format, args...)
}

// Warnf ...
func Warnf(format string, args ...interface{}) {
	LogSink.Warnf(format, args...)
}

// Panicf ...
func Panicf(format string, args ...interface{}) {
	LogSink.Panicf(format, args...)
}

// Fatalf ...
func Fatalf(format string, args ...interface{}) {
	LogSink.Fatalf(format, args...)
}

// Debugw ...
func Debugw(msg string, keysAndValues ...interface{}) {
	LogSink.Debugw(msg, keysAndValues...)
}

// Infow ...
func Infow(msg string, keysAndValues ...interface{}) {
	LogSink.Infow(msg, keysAndValues...)
}

// Warnw ...
func Warnw(msg string, keysAndValues ...interface{}) {
	LogSink.Warnw(msg, keysAndValues...)
}

// Errorw ...
func Errorw(msg string, keysAndValues ...interface{}) {
	LogSink.Errorw(msg, keysAndValues...)
}

// DPanicw ...
func DPanicw(msg string, keysAndValues ...interface{}) {
	LogSink.DPanicw(msg, keysAndValues...)
}

// Panicw ...
func Panicw(msg string, keysAndValues ...interface{}) {
	LogSink.Panicw(msg, keysAndValues...)
}

// Fatalw ...
func Fatalw(msg string, keysAndValues ...interface{}) {
	LogSink.Fatalw(msg, keysAndValues...)
}
