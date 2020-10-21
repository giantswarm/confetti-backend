package cmd

import (
	"github.com/spf13/cobra"
)

const (
	flagPort = "port"
)

type flag struct {
	Port int
}

func (f *flag) Init(cmd *cobra.Command) {
	cmd.PersistentFlags().IntVar(&f.Port, flagPort, 7777, "Set the port that the application will run on.")
}

func (f *flag) Validate() error {
	return nil
}
