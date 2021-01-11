package logger

import (
	"fmt"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
)

type Config struct {
	WrappedLogger micrologger.Logger
}

// Logger wraps micrologger.Logger to create a new logger that
// implements the atreugo.Logger interface, for usage in the
// web server.
type Logger struct {
	wrappedLogger micrologger.Logger
}

func New(c Config) (*Logger, error) {
	if c.WrappedLogger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.WrappedLogger must not be empty", c)
	}

	l := &Logger{
		wrappedLogger: c.WrappedLogger,
	}

	return l, nil
}

func (l *Logger) Print(v ...interface{}) {
	l.wrappedLogger.Log("level", "debug", "message", v)
}

func (l *Logger) Printf(format string, v ...interface{}) {
	l.wrappedLogger.Log("level", "debug", "message", fmt.Sprintf(format, v...))
}
