package events

import (
	"fmt"
	"sync"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/confetti-backend/internal/flags"
	types "github.com/giantswarm/confetti-backend/pkg/server/models/events/types"
)

type RepositoryConfig struct {
	Flags *flags.Flags
}

type Repository struct {
	flags *flags.Flags

	mu     sync.Mutex
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
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.events, nil
}

func (r *Repository) FindOneByID(id string) (types.Event, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, e, exists := r.findByID(id)
	if exists {
		return e, nil
	}

	return nil, microerror.Maskf(notFoundError, fmt.Sprintf("couldn't find any event with ID %s", id))
}

func (r *Repository) Update(event types.Event) (types.Event, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

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
