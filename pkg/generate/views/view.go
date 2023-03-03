package views

import (
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
)

type View interface {
	tea.Model
	Keymap() help.KeyMap
}
