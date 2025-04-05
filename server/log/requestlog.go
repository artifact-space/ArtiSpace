package log

import (
	"log/slog"
	"time"

	"github.com/go-chi/httplog/v2"
)

func HttpLogger() *httplog.Logger {
	return httplog.NewLogger("arti-space", httplog.Options{
		LogLevel:         slog.LevelDebug,
		Concise:          true,
		RequestHeaders:   true,
		MessageFieldName: "message",
		Tags:             map[string]string{},
		QuietDownRoutes:  []string{},
		QuietDownPeriod:  10 * time.Second,
	})
}