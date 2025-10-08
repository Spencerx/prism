package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"go.dalton.dog/prism/internal"
)

var configCmd = &cobra.Command{
	Use:          "config",
	Short:        "Manage Prism configuration",
	SilenceUsage: true,
}

var configClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Remove any prism config files/dirs",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return internal.ClearConfig()
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Display current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		internal.PrintConfig(internal.GlobalConfig)
		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Persist a configuration value",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := strings.ToLower(args[0])
		value := args[1]

		// This is fine at this level since all config settings are bools
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("Invalid boolean value %q: %w", value, err)
		}

		return internal.SetConfig(internal.GlobalConfig, key, parsed)
	},
}

func init() {
	configCmd.AddCommand(configClearCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)
	rootCmd.AddCommand(configCmd)
}
