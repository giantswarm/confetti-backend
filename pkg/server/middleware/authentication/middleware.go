package authentication

import (
	"bytes"
	"net/http"

	"github.com/giantswarm/microerror"
	"github.com/savsgio/atreugo/v11"

	"github.com/giantswarm/confetti-backend/flags"
	"github.com/giantswarm/confetti-backend/pkg/server/context/user"
)

var (
	tokenPrefix = []byte("Bearer")
)

type Config struct {
	Flags *flags.Flags
}

type Middleware struct {
	flags *flags.Flags
}

func New(c Config) (*Middleware, error) {
	if c.Flags == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Flags must not be empty", c)
	}

	m := &Middleware{
		flags: c.Flags,
	}

	return m, nil
}

func (m *Middleware) Middleware(ctx *atreugo.RequestCtx) error {
	token := m.getAuthorizationToken(ctx)
	if len(token) < 1 {
		ctx.SetStatusCode(http.StatusUnauthorized)

		return microerror.Maskf(unauthorizedError, "you are not authenticated")
	}

	u := &user.User{
		Token: token,
	}
	user.SaveContext(ctx, u)

	err := ctx.Next()
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (m *Middleware) getAuthorizationToken(ctx *atreugo.RequestCtx) string {
	var token string
	{
		auth := ctx.Request.Header.Peek("Authorization")
		if bytes.HasPrefix(auth, tokenPrefix) {
			token = string(auth[len(tokenPrefix):])
		}
	}

	return token
}
