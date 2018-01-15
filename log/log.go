// Package log provides a normalized interface for logging.
package log // import "github.com/motki/core/log"

import (
	"io"
	"io/ioutil"
	stdlog "log"
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Logger is the standard interface services should make use of.
type Logger logrus.FieldLogger

type outputType int

const (
	OutputStdout outputType = iota
	OutputStderr
	OutputNull
)

// Config contains information on how to configure a logger.
type Config struct {
	Level      string     `toml:"level"`
	OutputType outputType `toml:"output_type"`
}

// New creates and configures a new Logger using the given Config.
func New(c Config) Logger {
	l, err := logrus.ParseLevel(c.Level)
	if err != nil {
		l = logrus.DebugLevel
	}
	logger := logrus.New()
	switch c.OutputType {
	case OutputStderr:
		logger.Out = os.Stderr
	case OutputNull:
		logger.Out = ioutil.Discard
	default:
		// do nothing.
	}
	logger.Level = l
	logger.Formatter = &logrus.TextFormatter{}
	// Re-check for the above error and log it as a warning if it exist
	if err != nil {
		logger.Warnf("invalid log level '%s', defaulting to '%s'", c.Level, l.String())
	}
	return logger
}

func StdLogger(l Logger, level string) (*stdlog.Logger, io.Closer, error) {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return nil, nil, err
	}
	var wc io.WriteCloser
	switch logger := l.(type) {
	case *logrus.Logger:
		wc = logger.WriterLevel(lvl)
	case *logrus.Entry:
		wc = logger.WriterLevel(lvl)
	default:
		return nil, nil, errors.New("unsupported logger type")
	}
	return stdlog.New(wc, "", 0), wc, nil
}
