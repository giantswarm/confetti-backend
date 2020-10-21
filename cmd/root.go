package cmd

import (
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/giantswarm/confetti-backend/pkg/project"
)

type Config struct {
	Stderr io.Writer
	Stdout io.Writer
}

func New(config Config) (*cobra.Command, error) {
	if config.Stderr == nil {
		config.Stderr = os.Stderr
	}
	if config.Stdout == nil {
		config.Stdout = os.Stdout
	}

	f := &flag{}

	r := &runner{
		flag:   f,
		stderr: config.Stderr,
		stdout: config.Stdout,
	}

	c := &cobra.Command{
		Use:           project.Name(),
		Short:         project.Description(),
		Long:          project.Description(),
		RunE:          r.Run,
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       project.Version(),
	}

	f.Init(c)

	return c, nil
}
