package login

import (
	"testing"

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
			flags := flags.New()

			modelsConfig := models.Config{
				Flags: flags,
			}
			models, err := models.New(modelsConfig)
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}

			c := ServiceConfig{
				Flags:  flags,
				Models: models,
			}
			s, err := NewService(c)
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}

			results := make(map[string]bool)

			var result string
			for i := 0; i < tc.tries; i++ {
				result, err = s.Authenticate()
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
