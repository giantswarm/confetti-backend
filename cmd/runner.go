package cmd

import (
	"context"
	"fmt"
	"io"

	"github.com/atreugo/websocket"
	"github.com/giantswarm/microerror"
	"github.com/savsgio/atreugo/v11"
	"github.com/spf13/cobra"

	globalFlags "github.com/giantswarm/confetti-backend/flags"
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
		f.Port = r.flag.Port
	}

	var atreugoServer *atreugo.Atreugo
	{
		c := atreugo.Config{
			LogName:          project.Name(),
			Name:             project.Name(),
			Addr:             fmt.Sprintf("0.0.0.0:%d", f.Port),
			LogOutput:        r.stdout,
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
			AllowedOrigins: []string{"*"},
		}
		websocketUpgrader = websocket.New(c)
	}

	var s *server.Server
	{
		c := server.Config{
			Atreugo:           atreugoServer,
			Flags:             f,
			WebsocketUpgrader: websocketUpgrader,
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
