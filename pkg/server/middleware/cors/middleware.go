package cors

import (
	"github.com/giantswarm/microerror"
	"github.com/savsgio/atreugo/v11"

	"github.com/giantswarm/confetti-backend/internal/flags"
	"github.com/giantswarm/confetti-backend/pkg/server/models"
)

type Config struct {
	Flags  *flags.Flags
	Models *models.Model
}

type Middleware struct {
	flags  *flags.Flags
	models *models.Model
}

func New(c Config) (*Middleware, error) {
	if c.Flags == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Flags must not be empty", c)
	}
	if c.Models == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Models must not be empty", c)
	}

	m := &Middleware{
		flags:  c.Flags,
		models: c.Models,
	}

	return m, nil
}

func (m *Middleware) Middleware() atreugo.Middleware {
	return func(ctx *atreugo.RequestCtx) error {
		origin := string(ctx.Request.Header.Peek("Origin"))
		if !m.isOriginAllowed(m.flags.AllowedOrigins, origin) {
			return ctx.Next()
		}

		ctx.Response.Header.Set("Access-Control-Allow-Origin", origin)

		ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		ctx.Response.Header.Set("Access-Control-Allow-Headers", "Accept, Authorization, Cache-Control, Content-Type, DNT, If-Modified-Since, Keep-Alive, User-Agent, X-Request-ID, X-Requested-With")
		ctx.Response.Header.Set("Access-Control-Expose-Headers", "Location")
		ctx.Response.Header.Set("Access-Control-Max-Age", "86400")
		ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")

		return ctx.Next()
	}
}

func (m *Middleware) isOriginAllowed(allowed []string, origin string) bool {
	for _, v := range allowed {
		if v == origin || v == "*" {
			return true
		}
	}

	return false
}
