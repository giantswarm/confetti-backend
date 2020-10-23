package user

import (
	"github.com/savsgio/atreugo/v11"
)

const (
	key = "user"
)

type User struct {
	Token string
}

func SaveContext(ctx *atreugo.RequestCtx, u *User) {
	ctx.SetUserValue(key, u)
}

func FromContext(ctx *atreugo.RequestCtx) (*User, bool) {
	user := ctx.UserValue(key)
	if u, ok := user.(*User); ok {
		return u, ok
	}

	return nil, false
}
