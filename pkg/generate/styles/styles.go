package styles

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

var (
	BaseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240"))

	ColumnStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.HiddenBorder())

	FocusedStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62"))

	HelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	// TitleStyle is the style used for the program title
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Background(lipgloss.Color("#320ba8")).
			Padding(0, 1).
			Margin(1, 0, 1, 1)

	// HeaderStyle is the style used for the current view's header
	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			Margin(1, 0, 1, 1)
)

var (
	DefaultTableStyle table.Styles = table.DefaultStyles()
)

func init() {
	DefaultTableStyle.Header = DefaultTableStyle.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	DefaultTableStyle.Selected = DefaultTableStyle.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
}
