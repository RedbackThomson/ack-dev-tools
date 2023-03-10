package views

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/aws-controllers-k8s/dev-tools/pkg/generate/components/button"
	"github.com/aws-controllers-k8s/dev-tools/pkg/generate/constants"
	"github.com/aws-controllers-k8s/dev-tools/pkg/generate/utils"
)

var (
	subtle = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}

	dialogBoxStyle = lipgloss.NewStyle().
			Padding(1, 0)
)

var (
	DialogOKButtonID     string = "ok"
	DialogCancelButtonID string = "cancel"
)

var (
	DefaultDialogWidth int = 64
)

type dialogInputs struct {
	okButton     button.Model
	cancelButton button.Model
}

type Dialog struct {
	content string
	width   int
	visible bool

	inputs *dialogInputs
}

func NewDialog(content string, width int) *Dialog {
	inputs := &dialogInputs{
		okButton:     button.New(DialogOKButtonID, "Ok"),
		cancelButton: button.New(DialogCancelButtonID, "Cancel"),
	}
	inputs.okButton.Focus()

	m := &Dialog{
		content: content,
		width:   width,
		visible: false,
		inputs:  inputs,
	}

	return m
}

func (m *Dialog) SetVisible(visible bool) {
	m.visible = visible
}

func (m Dialog) IsVisible() bool {
	return m.visible
}

func (m Dialog) getInputFocusOrder() []utils.Focusable {
	return []utils.Focusable{
		&m.inputs.okButton,
		&m.inputs.cancelButton,
	}
}

func (m *Dialog) handleInputUpdates(msg tea.Msg) (Dialog, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, dialogKeys.Right):
			utils.RotateFocus(m.getInputFocusOrder(), utils.FocusRotateDown)
		case key.Matches(msg, dialogKeys.Left):
			utils.RotateFocus(m.getInputFocusOrder(), utils.FocusRotateUp)
		}
	}

	switch {
	case m.inputs.cancelButton.Focused():
		m.inputs.cancelButton, cmd = m.inputs.cancelButton.Update(msg)
	case m.inputs.okButton.Focused():
		m.inputs.okButton, cmd = m.inputs.okButton.Update(msg)
	}

	return *m, cmd
}

func (m Dialog) Update(msg tea.Msg) (Dialog, tea.Cmd) {
	switch msg := msg.(type) {
	case button.ButtonSelectMessage:
		switch msg.GetID() {
		case DialogOKButtonID:
			return m, (func() tea.Msg {
				return DialogOKMessage{}
			})
		case DialogCancelButtonID:
			return m, (func() tea.Msg {
				return DialogCancelMessage{}
			})
		}
	}

	return m.handleInputUpdates(msg)
}

func (m Dialog) View() string {
	if !m.visible {
		return ""
	}

	okButton := m.inputs.okButton.View()
	cancelButton := m.inputs.cancelButton.View()

	buttonSpacing := lipgloss.NewStyle().
		MarginRight(1)

	question := lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).MarginBottom(1).Render(m.content)
	buttons := lipgloss.JoinHorizontal(lipgloss.Top, buttonSpacing.Render(okButton), cancelButton)
	ui := lipgloss.JoinVertical(lipgloss.Center, question, buttons)

	return lipgloss.Place(constants.UsableViewSize.Width, constants.UsableViewSize.Height,
		lipgloss.Center, lipgloss.Center,
		dialogBoxStyle.Render(ui),
		lipgloss.WithWhitespaceChars("*"),
		lipgloss.WithWhitespaceForeground(subtle),
	)
}

func (m Dialog) Init() tea.Cmd {
	return nil
}

type DialogOKMessage struct{}
type DialogCancelMessage struct{}

type dialogKeyMap struct {
	Right  key.Binding
	Left   key.Binding
	Select key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k dialogKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Left, k.Right}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k dialogKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Left, k.Right}, // first column
	}
}

var dialogKeys = dialogKeyMap{
	Left: key.NewBinding(
		key.WithKeys("left"),
		key.WithHelp("←", "left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right"),
		key.WithHelp("→", "right"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select/toggle"),
	),
}
