package internal

import (
	"fmt"
	"maps"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/lipgloss/v2/table"
)

type BenchmarkSummary struct {
	sync.Mutex
	Results   []*BenchmarkResult
	Total     int
	Succeeded int
	Failed    int
	Skipped   int
	Duration  time.Duration
}

type BenchmarkResult struct {
	Name       string
	Package    string
	Threads    int
	Iterations int
	Metrics    map[string]string
	Output     []string
	RawLine    string
	Completed  bool
	Status     Status
	Duration   time.Duration
	StartedAt  time.Time
}

type BenchmarkPackageResults struct {
	Name       string
	Benchmarks []*BenchmarkResult
	Total      int
	Passed     int
	Failed     int
	Skipped    int
	Duration   time.Duration
}

type benchmarkLine struct {
	BaseName   string
	Threads    int
	Iterations int
	Metrics    map[string]string
	Raw        string
}

func displayBenchmarkResults(summary *BenchmarkSummary) {
	grouped := make(map[string]*BenchmarkPackageResults)
	for _, result := range summary.Results {
		pkgName := result.Package
		if _, exists := grouped[pkgName]; !exists {
			grouped[pkgName] = &BenchmarkPackageResults{
				Name:       pkgName,
				Benchmarks: make([]*BenchmarkResult, 0),
			}
		}

		pkg := grouped[pkgName]
		pkg.Benchmarks = append(pkg.Benchmarks, result)
		pkg.Total++
		switch result.Status {
		case StatusPass:
			pkg.Passed++
		case StatusFail:
			pkg.Failed++
		case StatusSkip:
			pkg.Skipped++
		}
		pkg.Duration += result.Duration
	}

	packageNames := make([]string, 0, len(grouped))
	for pkgName := range grouped {
		packageNames = append(packageNames, pkgName)
	}
	sort.Strings(packageNames)

	renderBlocks := make([]string, 0, len(packageNames))
	for _, pkgName := range packageNames {
		pkg := grouped[pkgName]
		sort.Slice(pkg.Benchmarks, func(i, j int) bool {
			nameI := strings.TrimPrefix(pkg.Benchmarks[i].Name, "Benchmark")
			nameJ := strings.TrimPrefix(pkg.Benchmarks[j].Name, "Benchmark")
			return nameI < nameJ
		})

		block := displayBenchmarkPackageBlock(pkg)
		if strings.TrimSpace(block) != "" {
			renderBlocks = append(renderBlocks, block)
		}
	}

	if len(renderBlocks) == 0 {
		displayZeroBenchmarks()
		return
	}

	if len(renderBlocks) > 1 {
		renderBlocks = append(renderBlocks, displayBenchmarkOverallSummary(summary))
	}

	mainChunk := lipgloss.JoinVertical(lipgloss.Left, renderBlocks...)

	fmt.Println(AppOverallOutputStyle.Render(mainChunk))
}

func displayBenchmarkPackageBlock(pkg *BenchmarkPackageResults) string {
	displayResults := filterBenchmarkResults(pkg.Benchmarks)
	if len(displayResults) == 0 {
		return ""
	}

	metricKeys := collectBenchmarkMetricKeys(displayResults)
	includeThreads := benchmarksHaveThreads(displayResults)

	headers := []string{"RESULT", "BENCHMARK"}
	if includeThreads {
		headers = append(headers, "P")
	}
	headers = append(headers, "ITER", "TIME")
	headers = append(headers, metricKeys...)

	metricOffset := len(headers) - len(metricKeys)
	timeColumn := metricOffset - 1

	rows := make([][]string, 0, len(displayResults))
	for _, bench := range displayResults {
		row := []string{
			bench.Status.String(),
			testNameStyle.Render(strings.TrimPrefix(bench.Name, "Benchmark")),
		}
		if includeThreads {
			row = append(row, benchmarkMetricStyle.Render(formatThreads(bench.Threads)))
		}
		row = append(row, benchmarkMetricStyle.Render(formatIterations(bench.Iterations)))
		row = append(row, durationStyle.Render(formatBenchmarkDuration(bench.Duration)))
		for _, key := range metricKeys {
			value := bench.Metrics[key]
			if value == "" {
				value = "—"
			}
			row = append(row, benchmarkMetricStyle.Render(value))
		}
		rows = append(rows, row)

		if len(bench.Output) > 0 && GlobalConfig.Verbose {
			for _, line := range bench.Output {
				outputRow := make([]string, len(headers))
				if len(headers) > 1 {
					outputRow[1] = outputStyle.Render(line)
				} else {
					outputRow[0] = outputStyle.Render(line)
				}
				rows = append(rows, outputRow)
			}
		}
	}

	t := table.New().
		Border(lipgloss.HiddenBorder()).
		BorderTop(false).BorderLeft(false).BorderBottom(false).BorderRight(false).
		Headers(headers...).
		Rows(rows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch row {
			case table.HeaderRow:
				return benchmarkHeaderStyle
			default:
				if col == 0 {
					return lipgloss.NewStyle()
				}
				if col == 1 {
					return testNameStyle
				}
				if col == timeColumn {
					return durationStyle
				}
				return benchmarkMetricStyle
			}
		})

	tableStr := t.Render()

	pkgHeader := fmt.Sprintf(
		"%s %s %s",
		benchmarkLabelStyle.Render("BENCH"),
		packageStyle.Render(pkg.Name),
		durationStyle.Render(fmt.Sprintf("(%v)", pkg.Duration.Truncate(time.Millisecond))),
	)

	summaryParts := []string{
		passStyle.Render(fmt.Sprintf("%d succeeded", pkg.Passed)),
		failStyle.Render(fmt.Sprintf("%d failed", pkg.Failed)),
	}
	if pkg.Skipped > 0 {
		summaryParts = append(summaryParts, skipStyle.Render(fmt.Sprintf("%d skipped", pkg.Skipped)))
	}
	pending := pkg.Total - pkg.Passed - pkg.Failed - pkg.Skipped
	if pending > 0 {
		summaryParts = append(summaryParts, benchmarkMetricStyle.Render(fmt.Sprintf("%d pending", pending)))
	}
	summaryParts = append(summaryParts, durationStyle.Render(fmt.Sprintf("%v total", pkg.Duration.Truncate(time.Millisecond))))

	pkgSummary := fmt.Sprintf("%d benchmarks • %s", pkg.Total, strings.Join(summaryParts, " • "))

	separatorLine := packageSeparatorStyle.Render(strings.Repeat("─", max(lipgloss.Width(tableStr), lipgloss.Width(pkgHeader))))

	return lipgloss.JoinVertical(lipgloss.Left,
		pkgHeader,
		benchmarkMetricStyle.Render(pkgSummary),
		separatorLine,
		pkgTableStyle.Render(tableStr),
		" ",
		" ",
	)
}

func displayBenchmarkOverallSummary(summary *BenchmarkSummary) string {
	summaryParts := []string{
		passStyle.Render(fmt.Sprintf("%d succeeded", summary.Succeeded)),
		failStyle.Render(fmt.Sprintf("%d failed", summary.Failed)),
	}
	if summary.Skipped > 0 {
		summaryParts = append(summaryParts, skipStyle.Render(fmt.Sprintf("%d skipped", summary.Skipped)))
	}
	pending := summary.Total - summary.Succeeded - summary.Failed - summary.Skipped
	if pending > 0 {
		summaryParts = append(summaryParts, benchmarkMetricStyle.Render(fmt.Sprintf("%d pending", pending)))
	}
	summaryParts = append(summaryParts, durationStyle.Render(fmt.Sprintf("%v total", summary.Duration.Truncate(time.Millisecond))))

	out := fmt.Sprintf("%d benchmarks • %s", summary.Total, strings.Join(summaryParts, " • "))
	return pkgTableStyle.AlignHorizontal(lipgloss.Left).MarginBottom(0).Render(out)
}

func displayZeroBenchmarks() {
	fmt.Println(benchmarkNoticeStyle.Render("No benchmarks found. Add Benchmark functions to your tests!"))
}

func processBenchmarkEvent(event *TestEvent, benchmarkMap map[string]*BenchmarkResult, summary *BenchmarkSummary) {
	if event.Test == "" || !strings.HasPrefix(event.Test, "Benchmark") {
		return
	}

	key := event.Package + "/" + event.Test

	summary.Lock()
	result, exists := benchmarkMap[key]
	if !exists {
		result = &BenchmarkResult{
			Name:    event.Test,
			Package: event.Package,
			Metrics: make(map[string]string),
			Output:  make([]string, 0),
			Status:  StatusRunning,
		}
		benchmarkMap[key] = result
		summary.Results = append(summary.Results, result)
		summary.Total++
	}
	summary.Unlock()

	action := Status(event.Action)
	switch action {
	case StatusOutput:
		line := strings.TrimSpace(event.Output)
		if line == "" {
			return
		}

		parsed, ok := parseBenchmarkLine(line)
		if !ok {
			summary.Lock()
			result.Output = append(result.Output, line)
			summary.Unlock()
			return
		}

		summary.Lock()
		if parsed.BaseName != "" {
			result.Name = parsed.BaseName
		}
		result.Threads = parsed.Threads
		result.Iterations = parsed.Iterations
		if result.Metrics == nil {
			result.Metrics = make(map[string]string)
		}

		maps.Copy(result.Metrics, parsed.Metrics)

		result.RawLine = parsed.Raw
		if result.StartedAt.IsZero() {
			result.StartedAt = event.Time
		}
		if !event.Time.IsZero() && !result.StartedAt.IsZero() {
			duration := event.Time.Sub(result.StartedAt)
			setBenchmarkDuration(summary, result, duration)
		}
		if result.Status != StatusFail {
			setBenchmarkStatus(summary, result, StatusPass)
			result.Completed = true
		}
		summary.Unlock()

	case StatusFail:
		summary.Lock()
		setBenchmarkStatus(summary, result, StatusFail)
		setBenchmarkDuration(summary, result, durationFromEvent(result, event.Time, event.Elapsed))
		result.Completed = true
		summary.Unlock()

	case StatusPass:
		summary.Lock()
		if result.Status != StatusFail {
			setBenchmarkStatus(summary, result, StatusPass)
			setBenchmarkDuration(summary, result, durationFromEvent(result, event.Time, event.Elapsed))
			result.Completed = true
		}
		summary.Unlock()

	case StatusSkip:
		summary.Lock()
		setBenchmarkStatus(summary, result, StatusSkip)
		setBenchmarkDuration(summary, result, durationFromEvent(result, event.Time, event.Elapsed))
		result.Completed = true
		summary.Unlock()

	case StatusRun:
		summary.Lock()
		if result.StartedAt.IsZero() {
			result.StartedAt = event.Time
		}
		summary.Unlock()
	}
}

func setBenchmarkStatus(summary *BenchmarkSummary, result *BenchmarkResult, newStatus Status) {
	if result.Status == newStatus {
		return
	}

	switch result.Status {
	case StatusPass:
		if summary.Succeeded > 0 {
			summary.Succeeded--
		}
	case StatusFail:
		if summary.Failed > 0 {
			summary.Failed--
		}
	case StatusSkip:
		if summary.Skipped > 0 {
			summary.Skipped--
		}
	}

	switch newStatus {
	case StatusPass:
		summary.Succeeded++
	case StatusFail:
		summary.Failed++
	case StatusSkip:
		summary.Skipped++
	}

	result.Status = newStatus
}

func parseBenchmarkLine(line string) (*benchmarkLine, bool) {
	trimmed := strings.TrimSpace(line)
	if !strings.HasPrefix(trimmed, "Benchmark") {
		return nil, false
	}
	if !strings.Contains(trimmed, "\t") {
		return nil, false
	}

	parts := strings.Split(trimmed, "\t")
	cleaned := make([]string, 0, len(parts))
	for _, part := range parts {
		value := strings.TrimSpace(part)
		if value != "" {
			cleaned = append(cleaned, value)
		}
	}

	if len(cleaned) < 2 {
		return nil, false
	}

	baseName, threads := splitBenchmarkName(cleaned[0])

	iterations, err := strconv.Atoi(cleaned[1])
	if err != nil {
		iterations = 0
	}

	metrics := make(map[string]string)
	for _, metricPart := range cleaned[2:] {
		fields := strings.Fields(metricPart)
		if len(fields) < 2 {
			continue
		}
		key := strings.Join(fields[1:], " ")
		metrics[key] = fields[0]
	}

	return &benchmarkLine{
		BaseName:   baseName,
		Threads:    threads,
		Iterations: iterations,
		Metrics:    metrics,
		Raw:        trimmed,
	}, true
}

func splitBenchmarkName(name string) (string, int) {
	idx := strings.LastIndex(name, "-")
	if idx == -1 {
		return name, 0
	}

	threads, err := strconv.Atoi(name[idx+1:])
	if err != nil {
		return name, 0
	}

	return name[:idx], threads
}

func filterBenchmarkResults(results []*BenchmarkResult) []*BenchmarkResult {
	filtered := make([]*BenchmarkResult, 0, len(results))
	for _, bench := range results {
		if GlobalConfig.OnlyFails && bench.Status != StatusFail {
			continue
		}
		filtered = append(filtered, bench)
	}
	return filtered
}

func collectBenchmarkMetricKeys(results []*BenchmarkResult) []string {
	seen := make(map[string]struct{})
	for _, bench := range results {
		for key := range bench.Metrics {
			if key != "" {
				seen[key] = struct{}{}
			}
		}
	}

	if len(seen) == 0 {
		return []string{}
	}

	ordered := make([]string, 0, len(seen))
	preferred := []string{"ns/op", "B/op", "allocs/op", "MB/s"}
	for _, key := range preferred {
		if _, ok := seen[key]; ok {
			ordered = append(ordered, key)
			delete(seen, key)
		}
	}

	remaining := make([]string, 0, len(seen))
	for key := range seen {
		remaining = append(remaining, key)
	}
	sort.Strings(remaining)
	ordered = append(ordered, remaining...)
	return ordered
}

func benchmarksHaveThreads(results []*BenchmarkResult) bool {
	for _, bench := range results {
		if bench.Threads > 0 {
			return true
		}
	}
	return false
}

func formatThreads(value int) string {
	if value <= 0 {
		return "—"
	}
	return fmt.Sprintf("%d", value)
}

func formatIterations(iterations int) string {
	if iterations <= 0 {
		return "—"
	}
	return fmt.Sprintf("%d", iterations)
}

func formatBenchmarkDuration(d time.Duration) string {
	if d <= 0 {
		return "—"
	}
	return fmt.Sprintf("%s", d.Truncate(time.Millisecond))
}

func setBenchmarkDuration(summary *BenchmarkSummary, result *BenchmarkResult, duration time.Duration) {
	if duration <= 0 {
		return
	}
	if result.Duration > 0 {
		summary.Duration -= result.Duration
	}
	result.Duration = duration
	summary.Duration += duration
}

func durationFromEvent(result *BenchmarkResult, eventTime time.Time, elapsed float64) time.Duration {
	if elapsed > 0 {
		return time.Duration(elapsed * float64(time.Second))
	}
	if !eventTime.IsZero() && !result.StartedAt.IsZero() {
		return eventTime.Sub(result.StartedAt)
	}
	return 0
}
