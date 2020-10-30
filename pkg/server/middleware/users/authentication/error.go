package authentication

import "github.com/giantswarm/microerror"

var invalidConfigError = &microerror.Error{
	Kind: "invalidConfigError",
}

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var unauthorizedError = &microerror.Error{
	Kind: "unauthorizedError",
}

// IsUnauthorized asserts unauthorizedError.
func IsUnauthorized(err error) bool {
	return microerror.Cause(err) == unauthorizedError
}
