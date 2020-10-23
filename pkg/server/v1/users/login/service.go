package login

import (
	"github.com/giantswarm/confetti-backend/flag"
)

type ServiceConfig struct {
	Flags *flag.Flag
}

type Service struct {
	flags *flag.Flag
}

func NewService(c ServiceConfig) (*Service, error) {
	service := &Service{
		flags: c.Flags,
	}

	return service, nil
}
