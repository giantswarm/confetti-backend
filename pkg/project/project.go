package project

var (
	description = ""
	gitSHA      = "n/a"
	name        = "confetti-backend"
	source      = "https://github.com/giantswarm/confetti-backend"
	version     = "0.1.0-dev"
)

func Description() string {
	return description
}

func GitSHA() string {
	return gitSHA
}

func Name() string {
	return name
}

func Source() string {
	return source
}

func Version() string {
	return version
}
