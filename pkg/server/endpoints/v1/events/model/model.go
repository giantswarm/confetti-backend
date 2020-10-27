package model

import (
	"fmt"

	"github.com/giantswarm/confetti-backend/flags"
	"github.com/giantswarm/confetti-backend/pkg/server/endpoints/v1/events/model/types"
	"github.com/giantswarm/microerror"
)

type RepositoryConfig struct {
	Flags *flags.Flags
}

type Repository struct {
	flags *flags.Flags

	events []types.Event
}

func NewRepository(c RepositoryConfig) (*Repository, error) {
	if c.Flags == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Flags must not be empty", c)
	}

	repository := &Repository{
		flags: c.Flags,

		events: types.MakeInitialData(),
	}

	return repository, nil
}

func (r *Repository) FindAll() ([]types.Event, error) {
	return r.events, nil
}

func (r *Repository) FindOneByID(id string) (types.Event, error) {
	_, e, exists := r.findByID(id)
	if exists {
		return e, nil
	}

	return nil, microerror.Maskf(notFoundError, fmt.Sprintf("couldn't find any event with ID %s", id))
}

func (r *Repository) Update(event types.Event) (types.Event, error) {
	id := event.ID()
	i, _, exists := r.findByID(id)
	if exists {
		r.events[i] = event
	}

	return nil, microerror.Maskf(notFoundError, fmt.Sprintf("couldn't find any event with ID %s", id))
}

func (r *Repository) findByID(id string) (int, types.Event, bool) {
	for i, e := range r.events {
		if e.ID() == id {
			return i, e, true
		}
	}

	return -1, nil, false
}