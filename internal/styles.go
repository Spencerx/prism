package internal

import (
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/lipgloss/v2/compat"
)

// --- Lipgloss Styles ---
var (
	redColor    = lipgloss.Color("1")
	greenColor  = lipgloss.Color("2")
	yellowColor = lipgloss.Color("3")

	// mainTextColor = compat.AdaptiveColor{Light: lipgloss.Color("0"), Dark: lipgloss.Color("15")}
	subTextColor = compat.AdaptiveColor{Light: lipgloss.Color("8"), Dark: lipgloss.Color("8")}

	passStyle = lipgloss.NewStyle().Foreground(greenColor).Bold(true)  // Light Green
	failStyle = lipgloss.NewStyle().Foreground(redColor).Bold(true)    // Light Red/Coral
	skipStyle = lipgloss.NewStyle().Foreground(yellowColor).Bold(true) // Pale Yellow

	packageStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Bold(true) // Light Aqua
	testNameStyle = lipgloss.NewStyle()                                             // Off-white
	durationStyle = lipgloss.NewStyle().Foreground(subTextColor)                    // Medium Gray

	outputStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Italic(true).MarginLeft(3)

	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true) // Orange-Red

	packageSeparatorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("0")) // Dark Gray

	// AppOverallOutputStyle is the top-level style that wraps all the display output.
	AppOverallOutputStyle = lipgloss.NewStyle().
				AlignHorizontal(lipgloss.Center).
				Margin(1)

	// Style for the package test table
	pkgTableStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Align(lipgloss.Center).
			MarginBottom(1)
)

var PrismHeader = `         /\
        /  \ #########
       /    \ ########
 #### /      \ #######
     /        \ ######
    /          \ #####
   /            \ ####
   ---------------`
