package models

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/confetti-backend/internal/flags"
	"github.com/giantswarm/confetti-backend/pkg/server/models/events"
)

type Config struct {
	Flags  *flags.Flags
	Logger micrologger.Logger
}

type Model struct {
	Events *events.Repository

	flags  *flags.Flags
	logger micrologger.Logger
}

func New(c Config) (*Model, error) {
	if c.Flags == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Flags must not be empty", c)
	}
	if c.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", c)
	}

	var err error

	var eventsModel *events.Repository
	{
		c := events.RepositoryConfig{
			Flags:  c.Flags,
			Logger: c.Logger,
		}

		eventsModel, err = events.NewRepository(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	m := &Model{
		Events: eventsModel,

		flags:  c.Flags,
		logger: c.Logger,
	}

	return m, nil
}
