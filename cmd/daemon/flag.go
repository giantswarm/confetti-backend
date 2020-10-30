package daemon

import (
	"github.com/spf13/cobra"
)

const (
	flagAddress       = "address"
	flagAllowedOrigin = "allowed-origin"
)

type flag struct {
	Address       string
	AllowedOrigin string
}

func (f *flag) Init(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&f.Address, flagAddress, "0.0.0.0:7777", "Set the address that the application will run on.")
	cmd.PersistentFlags().StringVar(&f.AllowedOrigin, flagAllowedOrigin, "*", "Set the allowed origin for connections.")
}

func (f *flag) Validate() error {
	return nil
}
