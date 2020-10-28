package watcher

import (
	"net/http"

	"github.com/giantswarm/microerror"
	"github.com/savsgio/atreugo/v11"

	"github.com/giantswarm/confetti-backend/pkg/server/endpoints/v1/events/model"
)

func ValidateIDMiddleware(service *Service) atreugo.Middleware {
	return func(ctx *atreugo.RequestCtx) error {
		id, ok := ctx.UserValue("id").(string)
		if !ok {
			ctx.SetStatusCode(http.StatusBadRequest)

			return microerror.Mask(invalidParamsError)
		}

		_, err := service.GetEventByID(id)
		if model.IsNotFoundError(err) {
			ctx.SetStatusCode(http.StatusNotFound)

			return microerror.Mask(err)
		} else if err != nil {
			return microerror.Mask(err)
		}

		return ctx.Next()
	}
}
