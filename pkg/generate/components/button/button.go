package button

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// KeyMap is the key bindings for different actions within the button.
type KeyMap struct {
	Select key.Binding
}

// DefaultKeyMap is the default set of key bindings for navigating and acting
// upon the button.
var DefaultKeyMap = KeyMap{
	Select: key.NewBinding(key.WithKeys("enter")),
}

// Model is the Bubble Tea model for this button element.
type Model struct {
	// General settings.
	SelectionPrefix string

	// Styles. These will be applied as inline styles.
	//
	// For an introduction to styling with Lip Gloss see:
	// https://github.com/charmbracelet/lipgloss
	FocusedStyle lipgloss.Style
	BlurredStyle lipgloss.Style
	TextStyle    lipgloss.Style

	id    string
	label string

	// focus indicates whether user input focus should be on this input
	// component. When false, ignore keyboard input and hide the cursor.
	focus bool

	// KeyMap encodes the keybindings recognized by the widget.
	KeyMap KeyMap
}

// New creates a new model with default settings.
func New(id string, text string) Model {
	return Model{
		SelectionPrefix: "",
		id:              id,
		label:           text,
		focus:           false,
		KeyMap:          DefaultKeyMap,
		FocusedStyle:    lipgloss.NewStyle().Background(lipgloss.Color("#777777")),
	}
}

// SetText sets the label of the button.
func (m *Model) SetText(s string) {
	m.label = s
}

// Text returns the text string of the button.
func (m Model) Text() string {
	return string(m.label)
}

// Focused returns the focus state on the model.
func (m Model) Focused() bool {
	return m.focus
}

// Focus sets the focus state on the model. When the model is in focus it can
// receive keyboard input.
func (m *Model) Focus() tea.Cmd {
	m.focus = true
	return nil
}

// Blur removes the focus state on the model.  When the model is blurred it can
// not receive keyboard input.
func (m *Model) Blur() {
	m.focus = false
}

// Update is the Bubble Tea update loop.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.focus {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, m.KeyMap.Select) {
			return m, m.SelectButton
		}
	}

	return m, nil
}

// View renders the textinput in its current state.
func (m Model) View() string {
	styleText := m.TextStyle.Inline(true).Render

	v := styleText(m.label)

	if m.focus {
		return m.FocusedStyle.Render(m.SelectionPrefix + v)
	}

	return v
}

type ButtonSelectMessage struct {
	id string
}

func (m ButtonSelectMessage) GetID() string {
	return m.id
}

func (m Model) SelectButton() tea.Msg {
	return ButtonSelectMessage{id: m.id}
}
