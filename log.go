package godoc2api

import "log"

// Possible log levels
const (
	LOG_DEBUG   = iota // log debug + following
	LOG_WARN           // log warnings + following
	LOG_ERR            // log errors + following
	LOG_NOTHING        // log nothing
)

// Current log level
var LogLevel uint = LOG_NOTHING

// Log with the level `debug`
func debug(str string, param ...interface{}) {
	if LogLevel <= LOG_DEBUG {
		log.Printf("DEBUG "+str, param...)
	}
}

// Log with the level `warning`
func warn(str string, param ...interface{}) {
	if LogLevel <= LOG_WARN {
		log.Printf("WARN "+str, param...)
	}
}

// Log with the level `error`
func problem(str string, param ...interface{}) {
	if LogLevel <= LOG_ERR {
		log.Printf("PROBLEM "+str, param...)
	}
}
