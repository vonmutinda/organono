package ctxhelper

import (
	"context"

	"github.com/vonmutinda/organono/app/entities"
)

func IPAddress(ctx context.Context) string {
	existing := ctx.Value(entities.ContextKeyIpAddress)
	if existing == nil {
		return ""
	}

	return existing.(string)
}

func WithIpAddress(ctx context.Context, ipAddress string) context.Context {
	return context.WithValue(ctx, entities.ContextKeyIpAddress, ipAddress)
}
