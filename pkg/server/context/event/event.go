package event

import (
	"github.com/savsgio/atreugo/v11"

	eventsModelTypes "github.com/giantswarm/confetti-backend/pkg/server/models/events/types"
)

const (
	key = "event"
)

func SaveContext(ctx *atreugo.RequestCtx, e eventsModelTypes.Event) {
	ctx.SetUserValue(key, e)
}

func FromContext(ctx *atreugo.RequestCtx) (eventsModelTypes.Event, bool) {
	return FromValueGetter(ctx.UserValue)
}

func FromValueGetter(getter func(key string) interface{}) (eventsModelTypes.Event, bool) {
	event := getter(key)
	if u, ok := event.(eventsModelTypes.Event); ok {
		return u, ok
	}

	return nil, false
}
