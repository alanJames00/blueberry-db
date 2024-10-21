// manages server wide logging
package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// initialize logger with loaded configs
func InitLogger(logLevel string) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix;
	
	// set the log level	
	setLogLevel(logLevel)

	// output log to stdout
	log.Logger = log.Output(os.Stdout);
}

// sets log level: debug, info and error
func setLogLevel(level string) {
	switch level {
	case "debug" :
		zerolog.SetGlobalLevel(zerolog.DebugLevel);
	case "info" : 
		zerolog.SetGlobalLevel(zerolog.InfoLevel);
	case "error" :
		zerolog.SetGlobalLevel(zerolog.ErrorLevel);
	}
}

// logging functions with different log levels
func Debug(message string, fields ...func(e *zerolog.Event)) {
	event := log.Debug();

	for _, field := range fields {
		field(event)
	}

	event.Msg(message);
}

func Info(message string, fields ...func(e *zerolog.Event)) {
    event := log.Info()
    for _, field := range fields {
        field(event)
    }
    event.Msg(message)
}

func Error(message string, fields ...func(e *zerolog.Event)) {
    event := log.Error()
    for _, field := range fields {
        field(event)
    }
    event.Msg(message)
}
