package app

import (
	"log/slog"
	"os"
	"strings"
)

func SetupLog(appEnv, nodeEnv string) {
	env := strings.TrimSpace(appEnv)
	if env == "" {
		env = strings.TrimSpace(nodeEnv)
	}
	var handler slog.Handler
	if env == "prod" || env == "production" {
		handler = slog.NewJSONHandler(os.Stdout, nil)
	} else {
		handler = slog.NewTextHandler(os.Stdout, nil)
	}
	slog.SetDefault(slog.New(handler))
}
