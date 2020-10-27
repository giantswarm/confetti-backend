package model

import (
	"fmt"

	"github.com/giantswarm/confetti-backend/flags"
	"github.com/giantswarm/microerror"
)

type RepositoryConfig struct {
	Flags *flags.Flags
}

type Repository struct {
	flags *flags.Flags

	events []Event
}

func NewRepository(c RepositoryConfig) (*Repository, error) {
	if c.Flags == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Flags must not be empty", c)
	}

	repository := &Repository{
		flags: c.Flags,
	}

	return repository, nil
}

func (r *Repository) FindAll() ([]Event, error) {
	return r.events, nil
}

func (r *Repository) FindOneByID(id EventID) (Event, error) {
	_, e, exists := r.findByID(id)
	if exists {
		return e, nil
	}

	return nil, microerror.Maskf(notFoundError, fmt.Sprintf("couldn't find any event with ID %s", id))
}

func (r *Repository) Update(event Event) (Event, error) {
	id := event.(*BaseEvent).ID
	i, _, exists := r.findByID(id)
	if exists {
		r.events[i] = event
	}

	return nil, microerror.Maskf(notFoundError, fmt.Sprintf("couldn't find any event with ID %s", id))
}

func (r *Repository) findByID(id EventID) (int, Event, bool) {
	for i, e := range r.events {
		if e.(*BaseEvent).ID == id {
			return i, e, true
		}
	}

	return -1, nil, false
}