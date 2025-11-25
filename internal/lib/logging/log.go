package logging

import (
	"context"
	"github.com/rs/zerolog/pkgerrors"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	fatalLevel = "fatal"
	errorLevel = "error"
	warnLevel  = "warn"
	infoLevel  = "info"
	debugLevel = "debug"
)

func Init(logLevel string) {
	var zeroLogLevel zerolog.Level

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	level := strings.ToLower(logLevel)

	switch level {
	case fatalLevel:
		zeroLogLevel = zerolog.FatalLevel
	case errorLevel:
		zeroLogLevel = zerolog.ErrorLevel
	case warnLevel:
		zeroLogLevel = zerolog.WarnLevel
	case debugLevel:
		zeroLogLevel = zerolog.DebugLevel
	default:
		zeroLogLevel = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(zeroLogLevel)
}

type LogFieldsKey struct{}

func Info(ctx context.Context) *zerolog.Event {
	fields := ctx.Value(LogFieldsKey{})
	return log.Ctx(log.Logger.WithContext(ctx)).Info().Fields(fields)
}

func Debug(ctx context.Context) *zerolog.Event {
	fields := ctx.Value(LogFieldsKey{})
	return log.Ctx(log.Logger.WithContext(ctx)).Debug().Fields(fields)
}

func Warn(ctx context.Context) *zerolog.Event {
	fields := ctx.Value(LogFieldsKey{})
	return log.Ctx(log.Logger.WithContext(ctx)).Warn().Fields(fields)
}

func Error(ctx context.Context, err error) *zerolog.Event {
	fields := ctx.Value(LogFieldsKey{})
	return log.Ctx(log.Logger.WithContext(ctx)).Error().Fields(fields).Stack().Err(err)
}

func Panic(ctx context.Context, err error) *zerolog.Event {
	fields := ctx.Value(LogFieldsKey{})
	return log.Ctx(log.Logger.WithContext(ctx)).Panic().Fields(fields).Stack().Err(err)
}

func Fatal(ctx context.Context, err error) *zerolog.Event {
	fields := ctx.Value(LogFieldsKey{})
	return log.Ctx(log.Logger.WithContext(ctx)).Fatal().Fields(fields).Stack().Err(err)
}
