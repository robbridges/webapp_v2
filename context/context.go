package context

import (
	"context"
	"github.com/robbridges/webapp_v2/models"
)

type key string

const (
	userKey key = "user"
)

func WithUser(ctx context.Context, user *models.User) context.Context {
	ctx = context.WithValue(ctx, userKey, user)
	return ctx
}

func User(ctx context.Context) *models.User {
	val := ctx.Value(userKey)

	user, ok := val.(*models.User)
	if !ok {
		// we should only hit this if there is an invalid user, or nothing was stored in the first place because we've
		// told go we're returning a models.User pointer it will know the type of this nil
		return nil
	}

	return user
}
