package jobs

import (
	"fmt"
	"log/slog"
)

type logger struct{}

func (l *logger) buildMetadata(keysAndValues ...interface{}) map[string]any {
	md := make(map[string]any, len(keysAndValues)/2)

	for i := 0; i+1 < len(keysAndValues); i += 2 {
		key := fmt.Sprint(keysAndValues[i])
		md[key] = keysAndValues[i+1]
	}

	return md
}

func (l *logger) Info(msg string, keysAndValues ...interface{}) {
	msg = fmt.Sprintf("[scheduler]: %s", msg)
	md := l.buildMetadata(keysAndValues...)
	if len(md) > 0 {
		slog.Info(msg, slog.Any("metadata", l.buildMetadata(keysAndValues...)))
	}

	slog.Info(msg)
}

func (l *logger) Error(err error, msg string, keysAndValues ...interface{}) {
	msg = fmt.Sprintf("[scheduler]: %s", msg)
	md := l.buildMetadata(keysAndValues...)
	if len(md) > 0 {
		slog.Error(msg,
			slog.Any("metadata", l.buildMetadata(keysAndValues...)),
			slog.Any("error", err),
		)
	}

	slog.Error(msg, slog.Any("error", err))
}
