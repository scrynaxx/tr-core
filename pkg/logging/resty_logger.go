package logging

import (
	"fmt"
	"log/slog"
)

type RestyLogger struct {
}

func NewRestyLogger() *RestyLogger {
	return &RestyLogger{}
}

func (l *RestyLogger) Errorf(format string, v ...interface{}) {
	slog.Error("[resty] error: " + fmt.Sprintf(format, v...))
}

func (l *RestyLogger) Warnf(format string, v ...interface{}) {
	slog.Warn("[resty] warning: " + fmt.Sprintf(format, v...))
}

func (l *RestyLogger) Debugf(format string, v ...interface{}) {
	slog.Debug("[resty] debug: " + fmt.Sprintf(format, v...))
}
