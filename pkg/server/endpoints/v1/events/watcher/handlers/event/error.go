package event

import "github.com/giantswarm/microerror"

var invalidConfigError = &microerror.Error{
	Kind: "invalidConfigError",
}

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var invalidEventTypeError = &microerror.Error{
	Kind: "invalidEventTypeError",
}

// IsInvalidEventType asserts invalidEventTypeError.
func IsInvalidEventType(err error) bool {
	return microerror.Cause(err) == invalidEventTypeError
}
