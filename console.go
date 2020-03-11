package mlog

//日志级别判断
func (l *Logger) enable(loglevel LogLevel) bool {
	return loglevel >= l.level
}

// Debug ...
func (l *Logger) Debug(format string, args ...interface{}) {
	l.writeLog(DEBUG, format, args...)
}

// Info ...
func (l *Logger) Info(format string, args ...interface{}) {
	l.writeLog(INFO, format, args...)
}

// Warn ...
func (l *Logger) Warn(format string, args ...interface{}) {
	l.writeLog(WARN, format, args...)
}

// Error ...
func (l *Logger) Error(format string, args ...interface{}) {
	l.writeLog(ERROR, format, args...)
}

// Fatal ...
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.writeLog(FATAL, format, args...)
}
