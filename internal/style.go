package internal

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// Color palette — soft, modern tones
var (
	colorGreen  = lipgloss.Color("#04B575") // success
	colorCyan   = lipgloss.Color("#06B6D4") // accent / alias
	colorRed    = lipgloss.Color("#FF6B6B") // error
	colorGray   = lipgloss.Color("#6C7086") // dim / secondary
	colorText   = lipgloss.Color("#CDD6F4") // list item text
)

// Predefined styles
var (
	successStyle = lipgloss.NewStyle().Bold(true).Foreground(colorGreen)
	aliasStyle   = lipgloss.NewStyle().Bold(true).Foreground(colorCyan)
	errorStyle   = lipgloss.NewStyle().Foreground(colorRed)
	dimStyle     = lipgloss.NewStyle().Foreground(colorGray)
	bulletStyle  = lipgloss.NewStyle().Foreground(colorCyan).Bold(true)
	listItemStyle = lipgloss.NewStyle().Foreground(colorText)
)

// Successf renders a success-prefixed message.
func Successf(format string, args ...interface{}) string {
	return successStyle.Render("  ✓  ") + fmt.Sprintf(format, args...)
}

// Errorf renders an error message in red.
func Errorf(format string, args ...interface{}) string {
	return errorStyle.Render(fmt.Sprintf(format, args...))
}

// Dimf renders a secondary/dim message.
func Dimf(format string, args ...interface{}) string {
	return dimStyle.Render(fmt.Sprintf(format, args...))
}

// Alias renders an alias name with accent color.
func Alias(name string) string {
	return aliasStyle.Render(name)
}

// Bullet renders a list bullet point.
func Bullet() string {
	return bulletStyle.Render("●")
}

// ListItem renders a single list entry: "  ● alias"
func ListItem(alias string) string {
	return fmt.Sprintf("  %s %s", Bullet(), listItemStyle.Render(alias))
}
