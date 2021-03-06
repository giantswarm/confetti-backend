package daemon

import (
	"context"
	"io"

	"github.com/atreugo/websocket"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/savsgio/atreugo/v11"
	"github.com/spf13/cobra"

	globalFlags "github.com/giantswarm/confetti-backend/internal/flags"
	wrappedLogger "github.com/giantswarm/confetti-backend/pkg/logger"
	"github.com/giantswarm/confetti-backend/pkg/project"
	"github.com/giantswarm/confetti-backend/pkg/server"
)

type runner struct {
	flag   *flag
	stdout io.Writer
	stderr io.Writer
}

func (r *runner) Run(cmd *cobra.Command, args []string) error {
	err := r.flag.Validate()
	if err != nil {
		return microerror.Mask(err)
	}

	err = r.run(cmd.Context(), cmd, args)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (r *runner) run(ctx context.Context, cmd *cobra.Command, args []string) error {
	var err error

	f := globalFlags.New()
	{
		f.Address = r.flag.Address
		f.AllowedOrigins = r.flag.AllowedOrigins
	}

	var logger micrologger.Logger
	{
		c := micrologger.Config{}

		logger, err = micrologger.New(c)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	var atreugoLogger *wrappedLogger.Logger
	{
		c := wrappedLogger.Config{
			WrappedLogger: logger,
		}

		atreugoLogger, err = wrappedLogger.New(c)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	var atreugoServer *atreugo.Atreugo
	{
		c := atreugo.Config{
			Name:             project.Name(),
			Addr:             f.Address,
			Logger:           atreugoLogger,
			GracefulShutdown: true,
			ErrorView: func(ctx *atreugo.RequestCtx, err error, statusCode int) {
				_ = ctx.JSONResponse(atreugo.JSON{"code": statusCode, "msg": err.Error()}, statusCode)
			},
		}
		atreugoServer = atreugo.New(c)
	}

	var websocketUpgrader *websocket.Upgrader
	{
		c := websocket.Config{
			AllowedOrigins:  f.AllowedOrigins,
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
		websocketUpgrader = websocket.New(c)
	}

	var s *server.Server
	{
		c := server.Config{
			Atreugo:           atreugoServer,
			Flags:             f,
			WebsocketUpgrader: websocketUpgrader,
			Logger:            logger,
		}
		s, err = server.New(c)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	err = s.Boot()
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
