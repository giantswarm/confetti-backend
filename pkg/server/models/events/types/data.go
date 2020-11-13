package events

func MakeInitialData() []Event {
	var christmasOnsite2020 *OnsiteEvent
	{
		christmasOnsite2020 = NewOnsiteEvent()
		christmasOnsite2020.BaseEvent.active = true
		christmasOnsite2020.BaseEvent.id = "christmas-onsite-2020"
		christmasOnsite2020.BaseEvent.name = "Christmas Onsite 2020"
		christmasOnsite2020.Rooms = []OnsiteEventRoom{
			{
				ID:            "main-stage",
				Name:          "Main Stage",
				ConferenceURL: "https://www.giantswarm.io",
				Description:   "Join Oliver to talk about your Kubernetes journey, your experience.",
			},
			{
				ID:            "mulled-wine",
				Name:          "Mulled wine",
				ConferenceURL: "https://www.giantswarm.io",
				Description:   "The Cocktailkunst team has been shaping bar culture in Germany for over 10 years and has already trained over 10,000 participants in cocktail courses and workshops. Some of the best bartenders in Germany are behind the project. Founder Stephan Hinz has been awarded several world championship titles and has published various specialist books on the topics of beverages and enjoyment. Look forward to perfect craftsmanship, real expert knowledge and a relaxed atmosphere!",
			},
			{
				ID:            "remote-work",
				Name:          "Remote work / Distributed teams / etc.",
				ConferenceURL: "https://www.giantswarm.io",
				Description: `REMOTE FTW!
Giant Swarm has been fully remote for 6 years now.
If you want to dive deeper into a remote company organization, join us at our stall. We are curious about your challenges.
We are not limited to remote only discussions. Other "new work buzzwords" like transparency, self organization, agile, ...and how we fill them with life may also be on the plate. We are happy to share and discuss <3.`,
			},
			{
				ID:            "monitoring",
				Name:          "Monitoring / logging (\"Observability\" sounds cooler - Joe)",
				ConferenceURL: "https://www.giantswarm.io",
				Description:   "Prometheus? Loki? Jaeger? Feeling lost? Monitoring and observability is core to running modern Cloud Native infrastructure - find out the latest and greatest about the pillars of metrics, logging, tracing, and all things observability, from Giant Swarm engineers in the field and the larger Giant Swarm community.",
			},
			{
				ID:            "security",
				Name:          "Security / Auditing",
				ConferenceURL: "https://www.giantswarm.io",
				Description: `--- CLASSIFIED ---

MESSAGE VIA SIGNAL W/ BASE64 ENC TO LEARN / SHARE / DISCUSS LATEST SECURITY AND AUDITING MEASURES DELETE AFTER READING

--- CLASSIFIED ---`,
			},
			{
				ID:            "release-engineering",
				Name:          "Release engineering",
				ConferenceURL: "https://www.giantswarm.io",
				Description:   "Cloud Native Infrastructure is powered by Continuous Integration and Continuous Deployment, and Release Engineering as a whole. Given that, the landscape is hard to navigate, with hundreds of tools, and many best practices. Talk with Giant Swarm Release Engineering experts, as well as other members of the larger Giant Swarm community.",
			},
			{
				ID:            "kubernetes",
				Name:          "Kubernetes (incl. operators)",
				ConferenceURL: "https://www.giantswarm.io",
				Description:   "Kubernetes has already placed itself as the market leader orchestration system, with Giant Swarm running production-grade Kubernetes for the last five years. It's fair to say we've picked up a few tips and tricks, including how to build world-class operators - discover Giant Swarm's thoughts on the past, present, and future of Kubernetes.",
			},
			{
				ID:            "managed-apps",
				Name:          "Managed apps",
				ConferenceURL: "https://www.giantswarm.io",
				Description:   "With Kubernetes being just one building block in your digital transformation, there's a growing need for all the other projects in the Cloud Native ecosystem, and beyond. With Giant Swarm's App Platform, we've learnt all about running and supporting applications on Kubernetes clusters, and enabling the future of the Cloud Native stack. Let's have a chat about managing Prometheus, Loki, Istio, and more.",
			},
			{
				ID:            "devops",
				Name:          "DevOps / Operations",
				ConferenceURL: "https://www.giantswarm.io",
				Description:   "We all wish it was as simple as pushing code, but Operations forms a large part of critical infrastructure. With modern improvements like DevOps and SecOps entering the space, room for innovation is at an all-time high. Chat with both your team members from other companies, and Giant Swarm engineers on everything that keeps the lights on.",
			},
			{
				ID:            "direkt-gruppe",
				Name:          "Partner - Direkt Gruppe",
				ConferenceURL: "https://www.giantswarm.io",
				Description:   "",
			},
			{
				ID:            "container-solutions",
				Name:          "Partner - Container Solutions",
				ConferenceURL: "https://www.giantswarm.io",
				Description:   "",
			},
			{
				ID:            "viadee",
				Name:          "Partner - Viadee",
				ConferenceURL: "https://www.giantswarm.io",
				Description:   "",
			},
			{
				ID:            "spare",
				Name:          "A spare stall",
				ConferenceURL: "https://www.giantswarm.io",
				Description:   "Just in case you get cold and you want to relax for a bit.",
			},
		}
	}

	events := []Event{
		christmasOnsite2020,
	}

	return events
}
