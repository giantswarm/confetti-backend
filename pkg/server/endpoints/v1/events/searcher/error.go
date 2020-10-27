package searcher

import "github.com/giantswarm/microerror"

var invalidConfigError = &microerror.Error{
	Kind: "invalidConfigError",
}

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var invalidParamsError = &microerror.Error{
	Kind: "invalidParamsError",
}

// IsInvalidParamsError asserts invalidParamsError.
func IsInvalidParamsError(err error) bool {
	return microerror.Cause(err) == invalidParamsError
}

var notFoundError = &microerror.Error{
	Kind: "notFoundError",
}

// IsInvalidParamsError asserts notFoundError.
func IsNotFoundError(err error) bool {
	return microerror.Cause(err) == notFoundError
}
