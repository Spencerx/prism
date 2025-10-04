package cmd

import (
	"context"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"

	"go.dalton.dog/prism/internal"
)

var Version = "1.2b"

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
	if err := fang.Execute(context.Background(), rootCmd, fang.WithoutCompletions(), fang.WithVersion(Version)); err != nil {
		os.Exit(1)
	}
}

func init() {
	internal.GlobalConfig = internal.Config{}
	rootCmd.PersistentFlags().BoolVarP(&internal.GlobalConfig.Verbose, "verbose", "v", false, "Include test sub-output")
	rootCmd.PersistentFlags().BoolVarP(&internal.GlobalConfig.OnlyFails, "only-fails", "f", false, "Only run failing tests")
}
