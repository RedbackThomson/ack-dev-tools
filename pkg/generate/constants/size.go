package constants

import (
	tea "github.com/charmbracelet/bubbletea"
)

var (
	// WindowSize determines the full size of the screen window
	WindowSize tea.WindowSizeMsg

	// ContainerViewSize determines the full size for the container of a view.
	// This is used by the main view to all child views inside a bordered box.
	ContainerViewSize tea.WindowSizeMsg

	// UsableViewSize determines the full size of any one particular view
	UsableViewSize tea.WindowSizeMsg
)
