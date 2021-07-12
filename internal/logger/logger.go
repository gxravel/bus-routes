package logger

import (
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog"
)

type Logger interface {
	// Formatted messaging with different severnity levels.
	Debugf(format string, data ...interface{})
	Infof(format string, data ...interface{})
	Warnf(format string, data ...interface{})
	Errorf(format string, data ...interface{})
	Fatalf(format string, data ...interface{})

	// Messaging with different severnity levels.
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)

	WithField(name string, value interface{}) Logger
	WithFields(pairs ...interface{}) Logger

	WithErr(err error) Logger
	WithStr(key string, val string) Logger
	WithModule(val string) Logger
}

const (
	DefaultLevel      = "debug"
	DefaultTimeFormat = zerolog.TimeFormatUnix

	logFieldErr = "err"
)

var DefaultOutput = os.Stdout

func Default() Logger {
	var log, err = New(DefaultLevel, DefaultOutput)
	if err != nil {
		panic(err)
	}

	return log
}

func New(level string, output io.Writer) (Logger, error) {
	var lvl, err = zerolog.ParseLevel(level)
	if err != nil {
		return nil, fmt.Errorf("unknown log level %s", level)
	}

	var zl = zerolog.New(output).Level(lvl)

	return &logger{
		Logger: &zl,
	}, nil
}

type logger struct {
	*zerolog.Logger
}

// Debug implements Debug method for logger.
func (zl *logger) Debug(msg string) {
	zl.Logger.Debug().Msg(msg)
}

// Info implements Info method for logger.
func (zl *logger) Info(msg string) {
	zl.Logger.Info().Msg(msg)
}

// Warn implements Warn method for logger.
func (zl *logger) Warn(msg string) {
	zl.Logger.Warn().Msg(msg)
}

// Error implements Error method for logger.
func (zl *logger) Error(msg string) {
	zl.Logger.Error().Msg(msg)
}

// Fatal implements Fatal method for logger.
func (zl *logger) Fatal(msg string) {
	zl.Logger.Fatal().Msg(msg)
}

// Debugf implements Debugf method for logger.
func (zl *logger) Debugf(format string, args ...interface{}) {
	zl.Logger.Debug().Msgf(format, args...)
}

// Infof implements Infof method for logger.
func (zl *logger) Infof(format string, args ...interface{}) {
	zl.Logger.Info().Msgf(format, args...)
}

// Warnf implements Warnf method for logger.
func (zl *logger) Warnf(format string, args ...interface{}) {
	zl.Logger.Warn().Msgf(format, args...)
}

// Errorf implements Errorf method for logger.
func (zl *logger) Errorf(format string, args ...interface{}) {
	zl.Logger.Error().Msgf(format, args...)
}

// Fatalf implements Fatalf method for logger.
func (zl *logger) Fatalf(format string, args ...interface{}) {
	zl.Logger.Fatal().Msgf(format, args...)
}

// WithField implements WithField method for logger.
func (zl *logger) WithField(name string, value interface{}) Logger {
	var outzl = zl.With().Interface(name, value).Logger()
	return &logger{Logger: &outzl}
}

// WithFields implements WithFields method for logger.
func (zl *logger) WithFields(pairs ...interface{}) Logger {
	var n = len(pairs)
	if n%2 != 0 {
		pairs = append(pairs, "")
	}

	var outzl = *zl.Logger
	for i := 0; i < n; i += 2 {
		outzl = outzl.With().Interface(fmt.Sprint(pairs[i]), pairs[i+1]).Logger()
	}

	return &logger{Logger: &outzl}
}

// WithErr implements WithErr method for logger.
func (zl *logger) WithErr(err error) Logger {
	var outzl = zl.With().Str(logFieldErr, err.Error()).Logger()
	return &logger{Logger: &outzl}
}

// WithStr implements WithStr method for logger.
func (zl *logger) WithStr(key string, val string) Logger {
	var outzl = zl.With().Str(key, val).Logger()
	return &logger{Logger: &outzl}
}

// WithModule implements WithModule method for logger.
func (zl *logger) WithModule(val string) Logger {
	return zl.WithStr("module", val)
}
