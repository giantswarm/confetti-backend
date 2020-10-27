package model

import (
	"github.com/giantswarm/confetti-backend/flags"
	"github.com/giantswarm/microerror"
)

type RepositoryConfig struct {
	Flags *flags.Flags
}

type Repository struct {
	flags *flags.Flags
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