package model

import (
	"github.com/giantswarm/confetti-backend/pkg/server/endpoints/v1/events/model"
	"github.com/giantswarm/confetti-backend/pkg/server/endpoints/v1/events/model/eventtypes"
)

func makeData() []model.Event {
	var christmasOnsite2020 *eventtypes.OnsiteEvent
	{
		christmasOnsite2020 = eventtypes.NewOnsiteEvent()
		christmasOnsite2020.BaseEvent.Active = true
		christmasOnsite2020.BaseEvent.ID = "christmas-onsite-2020"
		christmasOnsite2020.BaseEvent.Name = "Christmas Onsite 2020"
		christmasOnsite2020.Rooms = append(christmasOnsite2020.Rooms, []eventtypes.OnsiteEventRoom{
			{
				ID: "some-room",
				Name: "Some room",
				ConferenceURL: "https://www.giantswarm.io",
			},
		}...)
	}

	events := []model.Event{
		christmasOnsite2020,
	}

	return events
}