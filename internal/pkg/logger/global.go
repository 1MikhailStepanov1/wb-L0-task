package logger

//nolint:noctx
func Debug(msg string, args ...any) {
	globalLogger.Debug(msg, args...)
}

//nolint:noctx
func Info(msg string, args ...any) {
	globalLogger.Info(msg, args...)
}

//nolint:noctx
func Warn(msg string, args ...any) {
	globalLogger.Warn(msg, args...)
}

//nolint:noctx
func Error(msg string, args ...any) {
	globalLogger.Error(msg, args...)
}
