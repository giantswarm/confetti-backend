package flags

type Flags struct {
	Port int
}

func New() *Flags {
	f := &Flags{}

	return f
}
