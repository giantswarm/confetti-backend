package watcher

import (
	"net/http"

	"github.com/giantswarm/microerror"
	"github.com/savsgio/atreugo/v11"

	eventsModel "github.com/giantswarm/confetti-backend/pkg/server/models/events"
)

func ValidateIDMiddleware(service *Service) atreugo.Middleware {
	return func(ctx *atreugo.RequestCtx) error {
		id, ok := ctx.UserValue("id").(string)
		if !ok {
			ctx.SetStatusCode(http.StatusBadRequest)

			return microerror.Mask(invalidParamsError)
		}

		_, err := service.GetEventByID(id)
		if eventsModel.IsNotFoundError(err) {
			ctx.SetStatusCode(http.StatusNotFound)

			return microerror.Mask(err)
		} else if err != nil {
			return microerror.Mask(err)
		}

		return ctx.Next()
	}
}
