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

const (
	// ContextKeyLogID is the key to store log id in context.
	// Set this key to context when you request is coming.
	ContextKeyLogID = "LOG_ID"
)
