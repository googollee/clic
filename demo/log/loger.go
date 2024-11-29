package log

import (
	"context"
	"log/slog"
	"os"

	"github.com/googollee/clic"
)

var logger *slog.Logger

type config struct {
	Level slog.Level `clic:"level,info,log level [debug,info,warn,error]"`
}

func initLogger(ctx context.Context, cfg *config) error {
	logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: cfg.Level,
	}))
	return nil
}

func init() {
	clic.RegisterWithCallback("log", initLogger)
}

func Debug(msg string, args ...any) {
	logger.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	logger.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	logger.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	logger.Error(msg, args...)
}
