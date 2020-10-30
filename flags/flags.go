package flags

type Flags struct {
	Address       string
	AllowedOrigin string
}

func New() *Flags {
	f := &Flags{}

	return f
}
