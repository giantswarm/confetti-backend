package login

import (
	"testing"

	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/confetti-backend/internal/flags"
	"github.com/giantswarm/confetti-backend/pkg/server/models"
)

func TestService_generateToken(t *testing.T) {
	testCases := []struct {
		name  string
		tries int
	}{
		{
			name:  "tokens should not repeat themselves",
			tries: 64,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var err error

			flags := flags.New()

			var logger micrologger.Logger
			{
				c := micrologger.Config{}
				logger, err = micrologger.New(c)
				if err != nil {
					t.Fatalf("unexpected error: %s", err.Error())
				}
			}

			var model *models.Model
			{
				c := models.Config{
					Flags:  flags,
					Logger: logger,
				}
				model, err = models.New(c)
				if err != nil {
					t.Fatalf("unexpected error: %s", err.Error())
				}
			}

			var service *Service
			{
				c := ServiceConfig{
					Flags:  flags,
					Models: model,
					Logger: logger,
				}
				service, err = NewService(c)
				if err != nil {
					t.Fatalf("unexpected error: %s", err.Error())
				}
			}

			results := make(map[string]bool)

			var result string
			for i := 0; i < tc.tries; i++ {
				result, err = service.Authenticate()
				if err != nil {
					t.Fatalf("unexpected error: %s", err.Error())
				}

				if _, exists := results[result]; exists {
					t.Fatalf("expected tokens to be unique, got: %s", result)
				}
			}
		})
	}
}
