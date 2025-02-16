package context

import (
	"context"

	"github.com/nobbmaestro/lazyhis/pkg/domain/service"
)

type serviceCtxKey struct{}

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
