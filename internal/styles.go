package internal

import (
	"fmt"

	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/lipgloss/v2/compat"
)

// --- Lipgloss Styles ---
var (
	redColor     = lipgloss.Color("1")
	greenColor   = lipgloss.Color("2")
	yellowColor  = lipgloss.Color("3")
	blueColor    = lipgloss.Color("4")
	magentaColor = lipgloss.Color("5")

	cyanColor = lipgloss.Color("14")

	// mainTextColor = compat.AdaptiveColor{Light: lipgloss.Color("0"), Dark: lipgloss.Color("15")}
	subTextColor = compat.AdaptiveColor{Light: lipgloss.Color("8"), Dark: lipgloss.Color("8")}

	passStyle = lipgloss.NewStyle().Foreground(greenColor).Bold(true)  // Light Green
	failStyle = lipgloss.NewStyle().Foreground(redColor).Bold(true)    // Light Red/Coral
	skipStyle = lipgloss.NewStyle().Foreground(yellowColor).Bold(true) // Pale Yellow

	packageStyle  = lipgloss.NewStyle().Foreground(magentaColor).Bold(true) // Light Aqua
	testNameStyle = lipgloss.NewStyle()                                     // Off-white
	durationStyle = lipgloss.NewStyle().Foreground(subTextColor)            // Medium Gray

	benchmarkLabelStyle  = lipgloss.NewStyle().Foreground(cyanColor).Bold(true)
	benchmarkHeaderStyle = lipgloss.NewStyle().Foreground(subTextColor).Bold(true)
	benchmarkMetricStyle = lipgloss.NewStyle().Foreground(subTextColor)
	benchmarkNoticeStyle = lipgloss.NewStyle().Padding(1).Bold(true).Foreground(cyanColor)

	outputStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Italic(true).MarginLeft(3)

	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true) // Orange-Red

	packageSeparatorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("0")) // Dark Gray

	// AppOverallOutputStyle is the top-level style that wraps all the display output.
	AppOverallOutputStyle = lipgloss.NewStyle().
				AlignHorizontal(lipgloss.Center).
				MarginLeft(1)

	// Style for the package test table
	pkgTableStyle = lipgloss.NewStyle().
			Align(lipgloss.Center)
)

var HeaderStr = fmt.Sprintf("━━━━━▲%v%v%v%v%v\n",
	lipgloss.NewStyle().Foreground(redColor).Render("P"),
	lipgloss.NewStyle().Foreground(yellowColor).Render("R"),
	lipgloss.NewStyle().Foreground(greenColor).Render("I"),
	lipgloss.NewStyle().Foreground(blueColor).Render("S"),
	lipgloss.NewStyle().Foreground(magentaColor).Render("M"),
)

var FigletHeaderOne = lipgloss.NewStyle().Foreground(redColor).Render(` ____  ____  ____  ___  __  __ `)
var FigletHeaderTwo = lipgloss.NewStyle().Foreground(yellowColor).Render(`(  _ \(  _ \(_  _)/ __)(  \/  )`)
var FigletHeaderThr = lipgloss.NewStyle().Foreground(greenColor).Render(` )___/ )   / _)(_ \__ \ )    ( `)
var FigletHeaderFou = lipgloss.NewStyle().Foreground(blueColor).Render(`(__)  (_)\_)(____)(___/(_/\/\_)`)

func Header() string {
	return lipgloss.JoinVertical(lipgloss.Center, FigletHeaderOne, FigletHeaderTwo, FigletHeaderThr, FigletHeaderFou, "")
}

var PrismHeader = `         
	 /\
        /  \ #########
       /    \ ########
 #### /      \ #######
     /        \ ######
    /          \ #####
   /            \ ####
   ---------------`

func displayZeroState() {
	fmt.Println(lipgloss.NewStyle().
		Padding(1).
		Bold(true).
		Foreground(redColor).
		Render("No tests found. Get to writing!"),
	)
}
