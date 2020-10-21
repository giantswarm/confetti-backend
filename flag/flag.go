package flag

type Flag struct {
	Port int
}

func New() *Flag {
	f := &Flag{}

	return f
}
