package authentication

import (
	"bytes"
	"net/http"

	"github.com/giantswarm/microerror"
	"github.com/savsgio/atreugo/v11"

	"github.com/giantswarm/confetti-backend/internal/flags"
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

func (m *Middleware) Middleware() atreugo.Middleware {
	return func(ctx *atreugo.RequestCtx) error {
		token := m.getAuthorizationToken(ctx)
		if len(token) < 1 {
			ctx.SetStatusCode(http.StatusUnauthorized)

			return microerror.Maskf(unauthorizedError, "you are not authenticated")
		}

		u := &user.User{
			Token: token,
		}
		user.SaveContext(ctx, u)

		return ctx.Next()
	}
}

func (m *Middleware) getAuthorizationToken(ctx *atreugo.RequestCtx) string {
	var token string
	{
		// Authorization header.
		auth := ctx.Request.Header.Peek("Authorization")
		if bytes.HasPrefix(auth, tokenPrefix) {
			token = string(auth[len(tokenPrefix):])
		} else {
			// Parameter in URL.
			token = string(ctx.URI().QueryArgs().Peek("token"))
		}
	}

	return token
}
