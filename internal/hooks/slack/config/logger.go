package config

import (
	"log/slog"
)

type (
	logger struct {
		*slog.Logger
	}
)

func (l *logger) Output(calldepth int, s string) error {
	l.Info(s)

	return nil
}
