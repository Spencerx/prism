package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"

	"go.dalton.dog/prism/internal"
)

var (
	Version       = "1.3b"
	configLoadErr error
)

var rootCmd = &cobra.Command{
	Use:   "prism",
	Short: "Prism is a wrapper around go test to make it simple and beautiful",
	Long: `Prism is a wrapper around Go's built in test command that aims to make it beautiful and organized. 

Issues? Requests? Feedback? Let me know! -- github.com/DaltonSW/prism`,
	Args: cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		internal.Execute(args)
	},
}

func Execute() {
	if configLoadErr != nil {
		fmt.Fprintf(os.Stderr, "warning: unable to load persisted config: %v\n", configLoadErr)
	}
	if err := fang.Execute(context.Background(), rootCmd, fang.WithoutCompletions(), fang.WithVersion(Version)); err != nil {
		os.Exit(1)
	}
}

func init() {
	var cfg internal.Config
	cfg, configLoadErr = internal.LoadConfig()
	internal.GlobalConfig = cfg

	rootCmd.PersistentFlags().BoolVarP(&internal.GlobalConfig.Verbose, "verbose", "v", internal.GlobalConfig.Verbose, "Include test sub-output")
	rootCmd.PersistentFlags().BoolVarP(&internal.GlobalConfig.OnlyFails, "only-fails", "f", internal.GlobalConfig.OnlyFails, "Only run failing tests")
	rootCmd.PersistentFlags().BoolVar(&internal.GlobalConfig.NoLogo, "no-logo", internal.GlobalConfig.NoLogo, "Hide Prism logo header")
}
