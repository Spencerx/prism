package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/lipgloss/v2/table"
	"github.com/spf13/cobra"

	"go.dalton.dog/prism/internal"
)

var configCmd = &cobra.Command{
	Use:          "config",
	Short:        "Manage Prism configuration",
	SilenceUsage: true,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Display current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		out := cmd.OutOrStdout()

		table := table.New().
			Rows(
				[]string{"no_logo", fmt.Sprintf("%t", internal.GlobalConfig.NoLogo)},
				[]string{"only_fails", fmt.Sprintf("%t", internal.GlobalConfig.OnlyFails)},
				[]string{"verbose", fmt.Sprintf("%t", internal.GlobalConfig.Verbose)},
			).
			Border(lipgloss.HiddenBorder())

		fmt.Fprintln(out, table.String())
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

		switch key {
		case "no-logo", "no_logo":
			parsed, err := strconv.ParseBool(value)
			if err != nil {
				return fmt.Errorf("invalid boolean value %q: %w", value, err)
			}
			internal.GlobalConfig.NoLogo = parsed

			if err := internal.SaveConfig(internal.GlobalConfig); err != nil {
				return fmt.Errorf("save config: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "no_logo set to %t\n", parsed)
			return nil
		case "only-fails", "only_fails":
			parsed, err := strconv.ParseBool(value)
			if err != nil {
				return fmt.Errorf("invalid boolean value %q: %w", value, err)
			}
			internal.GlobalConfig.OnlyFails = parsed

			if err := internal.SaveConfig(internal.GlobalConfig); err != nil {
				return fmt.Errorf("save config: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "only_fails set to %t\n", parsed)
			return nil
		case "verbose":
			parsed, err := strconv.ParseBool(value)
			if err != nil {
				return fmt.Errorf("invalid boolean value %q: %w", value, err)
			}
			internal.GlobalConfig.Verbose = parsed

			if err := internal.SaveConfig(internal.GlobalConfig); err != nil {
				return fmt.Errorf("save config: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "verbose set to %t\n", parsed)
			return nil
		default:
			return fmt.Errorf("unknown configuration key %q", args[0])
		}
	},
}

func init() {
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)
	rootCmd.AddCommand(configCmd)
}
