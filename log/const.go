package log

type LoggerLevel string

const (
	// LoggerLevelDebug is the debug level
	LoggerLevelDebug LoggerLevel = "debug"
	// LoggerLevelInfo is the info level
	LoggerLevelInfo LoggerLevel = "info"
	// LoggerLevelWarn is the warn level
	LoggerLevelWarn LoggerLevel = "warn"
	// LoggerLevelError is the error level
	LoggerLevelError LoggerLevel = "error"
)
