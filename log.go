package doc2raml

import "log"

// Log level
const (
	LOG_DEBUG = iota // log debug
	LOG_WARN         // log warnings
	LOG_ERR          // log errors
	LOG_PROD         // log nothing
)

var LogLevel int = LOG_DEBUG

func debug(str string, param ...interface{}) {
	if LogLevel <= LOG_DEBUG {
		log.Printf(str, param...)
	}
}

func warn(str string, param ...interface{}) {
	if LogLevel <= LOG_WARN {
		log.Printf(str, param...)
	}
}

func problem(str string, param ...interface{}) {
	if LogLevel <= LOG_ERR {
		log.Printf(str, param...)
	}
}
