package utils

import tea "github.com/charmbracelet/bubbletea"

type Focusable interface {
	Focused() bool
	Focus() tea.Cmd
	Blur()
}
