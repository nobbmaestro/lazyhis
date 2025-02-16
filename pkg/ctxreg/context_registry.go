package ctxreg

import (
	"context"

	"github.com/nobbmaestro/lazyhis/domain/service"
	"github.com/nobbmaestro/lazyhis/pkg/config"
)

type serviceCtxKey struct{}
type configCtxKey struct{}

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
