package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/scrynaxx/tr-core/config"
	"github.com/scrynaxx/tr-core/pkg/logging"
	"github.com/scrynaxx/tr-core/services/gateway/internal/app"
)

func main() {
	logging.InitSlog()

	cfg, err := config.Load()
	if err != nil {
		slog.Error("[cmd] load config failed", slog.Any("error", err))
		os.Exit(1)
	}

	if err = app.Run(context.Background(), cfg); err != nil {
		slog.Error("[cmd] stopped with error", slog.Any("error", err))
		os.Exit(1)
	}
}
