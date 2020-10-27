package events

import (
	"github.com/giantswarm/confetti-backend/flags"
)

type ServiceConfig struct {
	Flags *flags.Flags
}

type Service struct {
	flags *flags.Flags
}

func NewService(c ServiceConfig) (*Service, error) {
	service := &Service{
		flags: c.Flags,
	}

	return service, nil
}

