package cmd

import (
	"github.com/spf13/cobra"
)

var (
	// Used for flags.
	cfgFile     string
	userLicense string

	rootCmd = &cobra.Command{
		Use:   "skaffer",
		Short: "A scaffolding application for similiar micro service based applications",
		Long:  ``,
	}
)

func Execute() error {
	return rootCmd.Execute()
}
