package context

import (
	"context"
	"log/slog"

	"github.com/nobbmaestro/lazyhis/pkg/config"
	"github.com/nobbmaestro/lazyhis/pkg/domain/service"
)

type serviceCtxKey struct{}
type configCtxKey struct{}
type loggerCtxKey struct{}

func NewContext() context.Context {
	return context.Background()
}

func WithService(
	ctx context.Context,
	historyService *service.HistoryService,
) context.Context {
	return context.WithValue(ctx, serviceCtxKey{}, historyService)
}

func GetService(ctx context.Context) *service.HistoryService {
	if val, ok := ctx.Value(serviceCtxKey{}).(*service.HistoryService); ok {
		return val
	}
	return nil
}

func WithConfig(
	ctx context.Context,
	cfg *config.UserConfig,
) context.Context {
	return context.WithValue(ctx, configCtxKey{}, cfg)
}

func GetConfig(ctx context.Context) *config.UserConfig {
	if val, ok := ctx.Value(configCtxKey{}).(*config.UserConfig); ok {
		return val
	}
	return nil
}

func WithLogger(
	ctx context.Context,
	logger *slog.Logger,
) context.Context {
	return context.WithValue(ctx, loggerCtxKey{}, logger)
}

func GetLogger(ctx context.Context) *slog.Logger {
	if val, ok := ctx.Value(loggerCtxKey{}).(*slog.Logger); ok {
		return val
	}
	return nil
}
