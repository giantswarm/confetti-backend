package types

func MakeInitialData() []Event {
	var christmasOnsite2020 *OnsiteEvent
	{
		christmasOnsite2020 = NewOnsiteEvent()
		christmasOnsite2020.BaseEvent.active = true
		christmasOnsite2020.BaseEvent.id = "christmas-onsite-2020"
		christmasOnsite2020.BaseEvent.name = "Christmas Onsite 2020"
		christmasOnsite2020.Rooms = append(christmasOnsite2020.Rooms, []OnsiteEventRoom{
			{
				ID: "some-room",
				Name: "Some room",
				ConferenceURL: "https://www.giantswarm.io",
			},
		}...)
	}

	events := []Event{
		christmasOnsite2020,
	}

	return events
}