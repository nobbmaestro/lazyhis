package registry

import (
	"context"
	"log/slog"

	"github.com/nobbmaestro/lazyhis/pkg/app"
	"github.com/nobbmaestro/lazyhis/pkg/config"
	"github.com/nobbmaestro/lazyhis/pkg/domain/service"
)

type ContextKey int

const (
	AppKey ContextKey = iota
	ServiceKey
	ConfigKey
	ConfigPathKey
	LoggerKey
)

type Registry struct {
	Context context.Context
}

type Option func(*Registry)

func NewRegistry(opts ...Option) Registry {
	r := Registry{context.Background()}
	for _, opt := range opts {
		opt(&r)
	}
	return r
}

func WithContext(context context.Context) Option {
	return func(r *Registry) {
		r.Context = context
	}
}

func WithApp(app *app.App) Option {
	return func(r *Registry) {
		r.Context = context.WithValue(r.Context, AppKey, app)
	}
}

func WithService(historyService *service.HistoryService) Option {
	return func(r *Registry) {
		r.Context = context.WithValue(r.Context, ServiceKey, historyService)
	}
}

func WithConfig(cfg *config.UserConfig) Option {
	return func(r *Registry) {
		r.Context = context.WithValue(r.Context, ConfigKey, cfg)
	}
}

func WithConfigPath(path string) Option {
	return func(r *Registry) {
		r.Context = context.WithValue(r.Context, ConfigPathKey, path)
	}
}

func WithLogger(logger *slog.Logger) Option {
	return func(r *Registry) {
		r.Context = context.WithValue(r.Context, LoggerKey, logger)
	}
}

func (r Registry) GetApp() *app.App {
	if val, ok := r.Context.Value(AppKey).(*app.App); ok {
		return val
	}
	return nil
}

func (r Registry) GetService() *service.HistoryService {
	if val, ok := r.Context.Value(ServiceKey).(*service.HistoryService); ok {
		return val
	}
	return nil
}

func (r Registry) GetConfig() *config.UserConfig {
	if val, ok := r.Context.Value(ConfigKey).(*config.UserConfig); ok {
		return val
	}
	return nil
}

func (r Registry) GetConfigPath() *config.UserConfig {
	if val, ok := r.Context.Value(ConfigPathKey).(*config.UserConfig); ok {
		return val
	}
	return nil
}

func (r Registry) GetLogger() *slog.Logger {
	if val, ok := r.Context.Value(LoggerKey).(*slog.Logger); ok {
		return val
	}
	return nil
}
