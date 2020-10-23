package login

import (
	"testing"
	"github.com/giantswarm/confetti-backend/flag"
)

func TestService_generateToken(t *testing.T) {
	testCases := []struct{
		name string
		tries int
	}{
		{
			name: "tokens should not repeat themselves",
			tries: 64,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func (t *testing.T) {
			flags := flag.New()

			c := ServiceConfig {
				Flags: flags,
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