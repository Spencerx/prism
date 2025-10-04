package cmd

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"go.dalton.dog/prism/internal"
)

var (
	benchTime  string
	benchNum   int
	benchNoMem bool
	benchTests bool
)

var benchCmd = &cobra.Command{
	Use:   "bench [regex] [path]",
	Short: "Run benchmarks for the given packages",
	Example: `
# No args will use '.' for the regex and ./... for path
prism bench 

# Run for a certain duration or number of cycles. Mutually exclusive
prism bench --time 1m30s
prism bench --num 500

# A single argument will supply the regex and still use ./... for the path
prism bench SimpleBench* -v

# If you want to specify a package, you MUST pass a regex match as the first arg
prism bench . ./pkg1 --no-mem
`,
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Debug("Bench called", "args", args)
		if benchTime != "" && benchNum > 0 {
			return fmt.Errorf("flags --time and --num cannot be used together")
		}
		if benchNum < 0 {
			return fmt.Errorf("--num must be zero or greater")
		}

		benchArgs := []string{}
		if len(args) == 0 {
			benchArgs = append(benchArgs, "-bench=.")
			benchArgs = append(benchArgs, "./...")
		} else {
			benchArgs = append(benchArgs, fmt.Sprintf("-bench=%v", args[0]))
			if len(args) > 1 {
				benchArgs = append(benchArgs, "./...")

			} else {
				benchArgs = append(benchArgs, args[1:]...)
			}
		}

		if !benchTests {
			benchArgs = append(benchArgs, "-run=^$")
		}
		if !benchNoMem {
			benchArgs = append(benchArgs, "-benchmem")
		}

		if benchTime != "" {
			benchArgs = append(benchArgs, fmt.Sprintf("-benchtime=%s", benchTime))
		} else if benchNum > 0 {
			benchArgs = append(benchArgs, fmt.Sprintf("-benchtime=%dx", benchNum))
		}

		internal.Execute(benchArgs)
		return nil
	},
}

func init() {
	benchCmd.Flags().StringVarP(&benchTime, "time", "t", "", "Run benchmarks for the specified duration (e.g. 1s)")
	benchCmd.Flags().IntVarP(&benchNum, "num", "n", 0, "Run benchmarks for the specified number of iterations")
	benchCmd.Flags().BoolVar(&benchNoMem, "no-mem", false, "Disable benchmark memory allocation statistics")
	benchCmd.Flags().BoolVar(&benchTests, "tests", false, "Include the tests in packages prior to benchmarking")
	benchCmd.MarkFlagsMutuallyExclusive("time", "num")
	rootCmd.AddCommand(benchCmd)
}
