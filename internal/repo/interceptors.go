package repo

import (
	"context"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

const ctxRepoKey = "repo"

// NewContext returns new child context with a Repo stored as a Value
func NewContext(ctx context.Context, r Repo) context.Context {
	return context.WithValue(ctx, ctxRepoKey, r)
}

// FromContext Extracts Repo instance from a given context.
// Will panic if context didn't contain database at `ctxRepoKey` key.
func FromContext(ctx context.Context) Repo {
	r, ok := ctx.Value(ctxRepoKey).(Repo)
	if !ok {
		log.Panic().
			Msg("Db unexpectedly is not presented in context")
	}
	return r
}

// NewInterceptorWithRepo builds a grpc interceptor instance that puts Repo instance to context
func NewInterceptorWithRepo(r Repo) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		return handler(NewContext(ctx, r), req)
	}
}
