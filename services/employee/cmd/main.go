package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/scrynaxx/tr-core/config"
	"github.com/scrynaxx/tr-core/pkg/logging"
	"github.com/scrynaxx/tr-core/services/employee/internal/app"
)

func main() {
	logging.InitSlog()

	cfg, err := config.Load()
	if err != nil {
		slog.Error("[cmd] failed to load config", slog.Any("error", err))
		os.Exit(1)
	}

	if err = app.Run(context.Background(), cfg); err != nil {
		slog.Error("[cmd] failed to run app", slog.Any("error", err))
		os.Exit(1)
	}
}
