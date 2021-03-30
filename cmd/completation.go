package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

func init() {
	completion.Flags().StringP("shell", "s", "", "Shell type for auto completion")
	rootCmd.AddCommand(completion)
}

var completion = &cobra.Command{
	Use:   "completion",
	Short: "Generate shell completion scripts",

	RunE: func(cmd *cobra.Command, args []string) error {
		shell, err := cmd.Flags().GetString("shell")
		if err != nil {
			return err
		}

		if shell == "" {
			return errors.New("error: value `--shell` is required")
		}

		switch shell {
		case "fish":
			return rootCmd.GenFishCompletion(rootCmd.OutOrStdout(), true)
		case "bash":
			return rootCmd.GenBashCompletion(cmd.OutOrStdout())
		case "zsh":
			return rootCmd.GenZshCompletion(cmd.OutOrStdout())
		case "powershell":
			return rootCmd.GenPowerShellCompletion(cmd.OutOrStdout())
		}

		return nil
	},
}
