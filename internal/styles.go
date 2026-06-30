package internal

import (
	"fmt"

	"github.com/charmbracelet/lipgloss/v2"
)

// --- Lipgloss Styles ---
var (
	passStyle = lipgloss.NewStyle().Foreground(lipgloss.Green).Bold(true)
	failStyle = lipgloss.NewStyle().Foreground(lipgloss.Red).Bold(true)
	skipStyle = lipgloss.NewStyle().Foreground(lipgloss.Yellow).Bold(true)

	packageStyle  = lipgloss.NewStyle().Foreground(lipgloss.Magenta).Bold(true)
	testNameStyle = lipgloss.NewStyle().Foreground(lipgloss.NoColor{})
	durationStyle = lipgloss.NewStyle().Foreground(lipgloss.BrightBlack)

	benchmarkLabelStyle  = lipgloss.NewStyle().Foreground(lipgloss.BrightCyan).Bold(true)
	benchmarkHeaderStyle = lipgloss.NewStyle().Foreground(lipgloss.BrightBlack).Bold(true)
	benchmarkMetricStyle = lipgloss.NewStyle().Foreground(lipgloss.BrightBlack)
	benchmarkNoticeStyle = lipgloss.NewStyle().Foreground(lipgloss.BrightCyan).Padding(0, 1, 1, 1).Bold(true)

	outputStyle = lipgloss.NewStyle().Foreground(lipgloss.White).Italic(true).MarginLeft(3)
	errorStyle  = lipgloss.NewStyle().Foreground(lipgloss.BrightRed).Bold(true)

	packageSeparatorStyle = lipgloss.NewStyle().Foreground(lipgloss.BrightBlack)

	// AppOverallOutputStyle is the top-level style that wraps all the display output.
	// The single top margin is emitted by spinnerMarginStyle before the spinner; the
	// summary then renders flush onto the line the spinner reclaims when it stops.
	AppOverallOutputStyle = lipgloss.NewStyle().AlignHorizontal(lipgloss.Center).MarginLeft(1)

	// spinnerMarginStyle gives the running spinner a top margin. The summary and the
	// zero-state notices render flush below it, so this is the only top margin shown.
	spinnerMarginStyle = lipgloss.NewStyle().MarginTop(1)

	// Style for the package test table
	pkgTableStyle = lipgloss.NewStyle().
			Align(lipgloss.Center)
)

func UnsetColors() {
	passStyle = passStyle.Foreground(lipgloss.NoColor{})
	failStyle = failStyle.Foreground(lipgloss.NoColor{})
	skipStyle = skipStyle.Foreground(lipgloss.NoColor{})

	packageStyle = packageStyle.Foreground(lipgloss.NoColor{})
	durationStyle = durationStyle.Foreground(lipgloss.NoColor{})

	benchmarkLabelStyle = benchmarkLabelStyle.Foreground(lipgloss.NoColor{})
	benchmarkHeaderStyle = benchmarkHeaderStyle.Foreground(lipgloss.NoColor{})
	benchmarkMetricStyle = benchmarkMetricStyle.Foreground(lipgloss.NoColor{})
	benchmarkNoticeStyle = benchmarkNoticeStyle.Foreground(lipgloss.NoColor{})

	outputStyle = outputStyle.Foreground(lipgloss.NoColor{})
	errorStyle = errorStyle.Foreground(lipgloss.NoColor{})

	packageSeparatorStyle = packageSeparatorStyle.Foreground(lipgloss.NoColor{})
}

var FigletHeaderOne = lipgloss.NewStyle().Foreground(lipgloss.Red).Render(` ____  ____  ____  ___  __  __ `)
var FigletHeaderTwo = lipgloss.NewStyle().Foreground(lipgloss.Yellow).Render(`(  _ \(  _ \(_  _)/ __)(  \/  )`)
var FigletHeaderThr = lipgloss.NewStyle().Foreground(lipgloss.Green).Render(` )___/ )   / _)(_ \__ \ )    ( `)
var FigletHeaderFou = lipgloss.NewStyle().Foreground(lipgloss.Blue).Render(`(__)  (_)\_)(____)(___/(_/\/\_)`)

func Header() string {
	return lipgloss.JoinVertical(lipgloss.Center, FigletHeaderOne, FigletHeaderTwo, FigletHeaderThr, FigletHeaderFou, "")
}

func displayZeroState() {
	fmt.Println(lipgloss.NewStyle().
		Padding(0, 1, 1, 1).
		Bold(true).
		Foreground(lipgloss.Red).
		Render("No tests found. Get to writing!"),
	)
}
