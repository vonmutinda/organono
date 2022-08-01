package ctxhelper

import (
	"context"

	"github.com/vonmutinda/organono/app/entities"
)

func TokenInfo(ctx context.Context) *entities.TokenInfo {
	existing := ctx.Value(entities.ContextKeyTokenInfo)
	if existing == nil {
		return &entities.TokenInfo{}
	}

	tokenInfo, ok := existing.(*entities.TokenInfo)
	if !ok {
		return &entities.TokenInfo{}
	}

	return tokenInfo
}

func WithTokenInfo(ctx context.Context, tokenInfo *entities.TokenInfo) context.Context {
	return context.WithValue(ctx, entities.ContextKeyTokenInfo, tokenInfo)
}

func WithUserID(ctx context.Context, userID int64) context.Context {
	tokenInfo := TokenInfo(ctx)
	tokenInfo.UserID = userID
	return WithTokenInfo(ctx, tokenInfo)
}

func UserID(ctx context.Context) int64 {
	return TokenInfo(ctx).UserID
}
