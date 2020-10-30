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
	user := ctx.UserValue(key)
	if u, ok := user.(*usersModelTypes.User); ok {
		return u, ok
	}

	return nil, false
}
