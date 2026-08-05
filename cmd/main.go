package main

import (
	"log/slog"
	"os"

	"github.com/scrynaxx/tr-core/config"
	"github.com/scrynaxx/tr-core/internal/app"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	cfg, err := config.Load()
	if err != nil {
		slog.Error("[cmd] load config", slog.Any("error", err))
		os.Exit(1)
	}

	if err = app.Run(cfg); err != nil {
		slog.Error("[cmd] app run", slog.Any("error", err))
		os.Exit(1)
	}
}
