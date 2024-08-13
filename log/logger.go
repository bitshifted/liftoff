package log

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

func Init(enableDebug bool) {
	if enableDebug {
		Logger = zerolog.New(zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
		}).Level(zerolog.DebugLevel).With().Timestamp().Caller().Logger()
		Logger.Debug().Msg("Debug logging enabled")
	} else {
		Logger = zerolog.New(zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
		}).Level(zerolog.InfoLevel).With().Timestamp().Logger()
	}

}
