package user

import (
	"github.com/savsgio/atreugo/v11"

	usersModelTypes "github.com/giantswarm/confetti-backend/pkg/server/models/users/types"
)

const (
	key = "user"
)

func SaveContext(ctx *atreugo.RequestCtx, u *usersModelTypes.User) {
	ctx.SetUserValue(key, u)
}

func FromContext(ctx *atreugo.RequestCtx) (*usersModelTypes.User, bool) {
	return FromValueGetter(ctx.UserValue)
}

func FromValueGetter(getter func(key string) interface{}) (*usersModelTypes.User, bool) {
	user := getter(key)
	if u, ok := user.(*usersModelTypes.User); ok {
		return u, ok
	}

	return nil, false
}
