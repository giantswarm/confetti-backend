package login

import (
	"crypto/rand"
	"fmt"

	"github.com/giantswarm/microerror"

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

func (s *Service) Authenticate() (string, error) {
	token, err := s.generateToken()
	if err != nil {
		return "", microerror.Mask(err)
	}

	return token, nil
}

func (s *Service) generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", microerror.Mask(err)
	}

	return fmt.Sprintf("%x", b), nil
}
