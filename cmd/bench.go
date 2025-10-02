package cmd

import "github.com/spf13/cobra"

var benchCmd = &cobra.Command{
	Use:     "bench <regex>",
	Args:    cobra.ExactArgs(1),
	Short:   "Run benchmarks for the given package",
	Example: "prism bench .",
}

func init() {
	rootCmd.AddCommand(benchCmd)
}
