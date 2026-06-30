package internal

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/yarlson/pin"
)

func Execute(args []string) {
	benchMode := containsBenchmarkFlag(args)

	cmdArgs := []string{"test", "-json"}

	if len(args) == 0 {
		cmdArgs = append(cmdArgs, "./...")
	} else {
		cmdArgs = append(cmdArgs, args...)
	}

	if benchMode {
		p := pin.New(" Running benchmarks...",
			pin.WithPosition(pin.PositionRight),
			pin.WithTextColor(pin.ColorCyan),
			pin.WithSpinnerColor(pin.ColorMagenta),
		)

		// Match the top margin of the summary so the spinner doesn't sit flush
		fmt.Print(spinnerMarginStyle.Render(""))
		cancel := p.Start(context.Background())
		defer cancel()
		stopTimer := startElapsedTimer(p, "Running benchmarks...")

		summary, err := runBenchmarks(cmdArgs)

		stopTimer()
		p.Stop()

		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"%s\n",
				errorStyle.Render(fmt.Sprintf("Error running benchmarks: %v", err)),
			)
			os.Exit(1)
		}

		if summary.Total == 0 {
			displayZeroBenchmarks()
		} else {
			displayBenchmarkResults(summary)
		}
	} else {
		p := pin.New(" Running tests...",
			pin.WithPosition(pin.PositionRight),
			pin.WithTextColor(pin.ColorCyan),
			pin.WithSpinnerColor(pin.ColorMagenta),
		)

		// Match the top margin of the summary so the spinner doesn't sit flush
		fmt.Print(spinnerMarginStyle.Render(""))
		cancel := p.Start(context.Background())
		defer cancel()
		stopTimer := startElapsedTimer(p, "Running tests...")

		start := time.Now()
		summary, err := runTests(cmdArgs)
		elapsed := time.Since(start)

		stopTimer()
		p.Stop()

		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"%s\n",
				errorStyle.Render(fmt.Sprintf("Error running tests: %v", err)),
			)
			os.Exit(1)
		}

		summary.Elapsed = elapsed

		// Capture all display output as a single string and wrap it
		if summary.Total == 0 {
			displayZeroState()
		} else {
			displayResults(summary)
		}
	}

}

// startElapsedTimer updates the spinner label once a second with the elapsed
// runtime, e.g. "Running tests... (3s)". The returned func stops the ticker and
// must be called before p.Stop().
func startElapsedTimer(p *pin.Pin, label string) func() {
	start := time.Now()
	ticker := time.NewTicker(time.Second)
	done := make(chan struct{})

	render := func() {
		elapsed := time.Since(start).Round(time.Second)
		p.UpdateMessage(fmt.Sprintf(" %s (%s)", label, elapsed))
	}
	render() // show "(0s)" immediately so the timer doesn't jump in at 1s

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				render()
			}
		}
	}()

	return func() {
		ticker.Stop()
		close(done)
	}
}

func containsBenchmarkFlag(args []string) bool {
	for i := range len(args) {
		arg := args[i]
		if strings.HasPrefix(arg, "-bench") || strings.HasPrefix(arg, "-test.bench") {
			return true
		}
	}
	return false
}

func runTests(args []string) (*TestSummary, error) {
	cmd := exec.CommandContext(context.Background(), "go", args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start command: %w", err)
	}

	summary := &TestSummary{
		Results:        make([]TestResult, 0),
		CachedPackages: make(map[string]bool),
	}
	testMap := make(map[string]*TestResult)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			var event TestEvent
			if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
				fmt.Fprintf(
					os.Stderr,
					"%s\n",
					errorStyle.Render(fmt.Sprintf(
						"Warning: Failed to unmarshal JSON event: %v (line: %s)",
						err,
						scanner.Text(),
					)),
				)
				continue
			}
			processEvent(&event, testMap, summary)
		}
		if err := scanner.Err(); err != nil && err != io.EOF {
			fmt.Fprintf(
				os.Stderr,
				"%s\n",
				errorStyle.Render(fmt.Sprintf("Error reading stdout: %v", err)),
			)
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			fmt.Fprintf(os.Stderr, "%s\n", scanner.Text())
		}
		if err := scanner.Err(); err != nil && err != io.EOF {
			fmt.Fprintf(
				os.Stderr,
				"%s\n",
				errorStyle.Render(fmt.Sprintf("Error reading stderr: %v", err)),
			)
		}
	}()

	cmdErr := cmd.Wait()

	wg.Wait()

	for _, result := range testMap {
		summary.Results = append(summary.Results, *result)
	}

	if cmdErr != nil {
		if exitErr, ok := cmdErr.(*exec.ExitError); ok {
			if exitErr.ExitCode() == 1 {
				return summary, nil
			}
			return nil, fmt.Errorf(
				"command exited with non-zero status %d: %w",
				exitErr.ExitCode(),
				cmdErr,
			)
		}
		return nil, fmt.Errorf("command execution failed: %w", cmdErr)
	}

	return summary, nil
}

func runBenchmarks(args []string) (*BenchmarkSummary, error) {
	cmd := exec.CommandContext(context.Background(), "go", args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start command: %w", err)
	}

	summary := &BenchmarkSummary{
		Results:    make([]*BenchmarkResult, 0),
		PackageEnv: make(map[string]*BenchmarkPackageEnv),
	}
	benchmarkMap := make(map[string]*BenchmarkResult)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			var event TestEvent
			if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
				fmt.Fprintf(
					os.Stderr,
					"%s\n",
					errorStyle.Render(fmt.Sprintf(
						"Warning: Failed to unmarshal JSON event: %v (line: %s)",
						err,
						scanner.Text(),
					)),
				)
				continue
			}
			processBenchmarkEvent(&event, benchmarkMap, summary)
		}
		if err := scanner.Err(); err != nil && err != io.EOF {
			fmt.Fprintf(
				os.Stderr,
				"%s\n",
				errorStyle.Render(fmt.Sprintf("Error reading stdout: %v", err)),
			)
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			fmt.Fprintf(os.Stderr, "%s\n", scanner.Text())
		}
		if err := scanner.Err(); err != nil && err != io.EOF {
			fmt.Fprintf(
				os.Stderr,
				"%s\n",
				errorStyle.Render(fmt.Sprintf("Error reading stderr: %v", err)),
			)
		}
	}()

	cmdErr := cmd.Wait()

	wg.Wait()

	if cmdErr != nil {
		if exitErr, ok := cmdErr.(*exec.ExitError); ok {
			if exitErr.ExitCode() == 1 {
				return summary, nil
			}
			return nil, fmt.Errorf(
				"command exited with non-zero status %d: %w",
				exitErr.ExitCode(),
				cmdErr,
			)
		}
		return nil, fmt.Errorf("command execution failed: %w", cmdErr)
	}

	return summary, nil
}

func processEvent(event *TestEvent, testMap map[string]*TestResult, summary *TestSummary) {
	if event.Test == "" {
		if Status(event.Action) == StatusOutput && strings.Contains(event.Output, "(cached)") {
			summary.Lock()
			summary.CachedPackages[event.Package] = true
			summary.Unlock()
		}
		return
	}

	key := event.Package + "/" + event.Test

	summary.Lock()
	defer summary.Unlock()

	result, exists := testMap[key]
	if !exists {
		result = &TestResult{
			Name:    event.Test,
			Package: event.Package,
			Status:  StatusRunning,
			Output:  make([]string, 0),
		}
		testMap[key] = result
		summary.Total++
	}

	switch action := Status(event.Action); action {
	case StatusOutput:
		output := strings.TrimSpace(event.Output)
		if output != "" {
			result.Output = append(result.Output, output)
		}

	case StatusPass, StatusFail, StatusSkip:
		result.Status = action
		result.Duration = time.Duration(event.Elapsed * float64(time.Second))

		switch action {
		case StatusPass:
			summary.Passed++
		case StatusFail:
			summary.Failed++
		case StatusSkip:
			summary.Skipped++
		}
	}
}
