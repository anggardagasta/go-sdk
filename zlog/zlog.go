package zlog

import (
	"context"
	"os"

	"github.com/rs/zerolog"
)

var log *Zlog

type ZlogLevel int

const (
	InfoLevel  ZlogLevel = ZlogLevel(zerolog.InfoLevel)
	ErrorLevel ZlogLevel = ZlogLevel(zerolog.ErrorLevel)
	FatalLevel ZlogLevel = ZlogLevel(zerolog.FatalLevel)
	DebugLevel ZlogLevel = ZlogLevel(zerolog.DebugLevel)
	WarnLevel  ZlogLevel = ZlogLevel(zerolog.WarnLevel)
	PanicLevel ZlogLevel = ZlogLevel(zerolog.PanicLevel)
	TraceLevel ZlogLevel = ZlogLevel(zerolog.TraceLevel)
)

type Zlog struct {
	log zerolog.Logger
}

type ZlogProperties struct {
	Level      ZlogLevel
	AppName    string
	AppVersion string
	AppEnv     string
}

func New(logz zerolog.Logger) *Zlog {
	log = &Zlog{
		log: logz,
	}

	return log
}

// Init initializes the global logger with a default configuration.
func Init() {
	if log == nil {
		New(zerolog.New(os.Stderr).With().Timestamp().Logger())
	}
}

func checkAndInit() {
	if log == nil {
		Init()
	}
}

func Info(ctx context.Context, data interface{}, msg string) {
	checkAndInit()
	log.log.Info().Fields(data).Msg(msg)
}

func Error(ctx context.Context, data interface{}, msg string) {
	checkAndInit()
	log.log.Error().Fields(data).Msg(msg)
}

func Fatal(ctx context.Context, data interface{}, msg string) {
	checkAndInit()
	log.log.Fatal().Fields(data).Msg(msg)
}

func Debug(ctx context.Context, data interface{}, msg string) {
	checkAndInit()
	log.log.Debug().Fields(data).Msg(msg)
}

func Warn(ctx context.Context, data interface{}, msg string) {
	checkAndInit()
	log.log.Warn().Fields(data).Msg(msg)
}

func Panic(ctx context.Context, data interface{}, msg string) {
	checkAndInit()
	log.log.Panic().Fields(data).Msg(msg)
}

func Trace(ctx context.Context, data interface{}, msg string) {
	checkAndInit()
	log.log.Trace().Fields(data).Msg(msg)
}
