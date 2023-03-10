package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	ackconfig "github.com/aws-controllers-k8s/code-generator/pkg/config"
	ackmodel "github.com/aws-controllers-k8s/code-generator/pkg/model"
	"github.com/aws-controllers-k8s/dev-tools/pkg/generate/components/button"
	"github.com/aws-controllers-k8s/dev-tools/pkg/generate/utils"
)

var (
	ResourceFormSpecFieldsButtonID    string = "spec"
	ResourceFormStatusFieldsButtonID  string = "status"
	ResourceFormARNPrimaryKeyButtonID string = "arn-primary-key"
)

type resourceFormInputs struct {
	specFieldsButton    button.Model
	statusFieldsButton  button.Model
	arnPrimaryKeyButton button.Model
}

type ResourceForm struct {
	crd    *ackmodel.CRD
	config *ackconfig.ResourceConfig

	inputs *resourceFormInputs
}

func (m ResourceForm) CRD() *ackmodel.CRD {
	return m.crd
}

func (m ResourceForm) Keymap() help.KeyMap {
	return resourceFormKeys
}

func NewResourceForm(crd *ackmodel.CRD, config *ackconfig.ResourceConfig) *ResourceForm {
	form := &ResourceForm{
		crd:    crd,
		config: config,

		inputs: &resourceFormInputs{
			specFieldsButton:    button.New(ResourceFormSpecFieldsButtonID, fmt.Sprintf("%d fields", len(crd.SpecFields))),
			statusFieldsButton:  button.New(ResourceFormStatusFieldsButtonID, fmt.Sprintf("%d fields", len(crd.StatusFields))),
			arnPrimaryKeyButton: button.New(ResourceFormARNPrimaryKeyButtonID, fmt.Sprintf("%t", crd.IsARNPrimaryKey())),
		},
	}

	return form
}

func (m ResourceForm) getInputFocusOrder() []utils.Focusable {
	return []utils.Focusable{
		&m.inputs.specFieldsButton,
		&m.inputs.statusFieldsButton,
		&m.inputs.arnPrimaryKeyButton,
	}
}

func (m *ResourceForm) handleInputUpdates(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, resourceFormKeys.LineDown):
			utils.RotateFocus(m.getInputFocusOrder(), utils.FocusRotateDown)
		case key.Matches(msg, resourceFormKeys.LineUp):
			utils.RotateFocus(m.getInputFocusOrder(), utils.FocusRotateUp)
		case key.Matches(msg, resourceFormKeys.Quit):
			return *m, (func() tea.Msg {
				return ReturnMessage{}
			})
		}
	}

	switch {
	case m.inputs.specFieldsButton.Focused():
		m.inputs.specFieldsButton, cmd = m.inputs.specFieldsButton.Update(msg)
	case m.inputs.statusFieldsButton.Focused():
		m.inputs.statusFieldsButton, cmd = m.inputs.statusFieldsButton.Update(msg)
	case m.inputs.arnPrimaryKeyButton.Focused():
		m.inputs.arnPrimaryKeyButton, cmd = m.inputs.arnPrimaryKeyButton.Update(msg)
	}

	return *m, cmd
}

func (m ResourceForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case button.ButtonSelectMessage:
		switch msg.GetID() {
		case ResourceFormSpecFieldsButtonID:
			return m, (func() tea.Msg {
				return OpenSpecFieldsMessage{}
			})
		case ResourceFormStatusFieldsButtonID:
			return m, (func() tea.Msg {
				return OpenStatusFieldsMessage{}
			})
		case ResourceFormARNPrimaryKeyButtonID:
			m.config.IsARNPrimaryKey = !m.config.IsARNPrimaryKey
			m.inputs.arnPrimaryKeyButton.SetLabel(fmt.Sprintf("%t", m.config.IsARNPrimaryKey))
			return m, (func() tea.Msg {
				return UpdateResourceConfig{
					Kind:   m.crd.Kind,
					Config: m.config,
				}
			})
		}
	}

	return m.handleInputUpdates(msg)
}

func (m ResourceForm) View() string {
	fieldNames := []string{
		"Kind",
		"Plural",
		"Spec Fields",
		"Status Fields",
		"Is ARN Primary Key?",
	}

	fieldViews := []string{
		m.crd.Kind,
		m.crd.Plural,
		m.inputs.specFieldsButton.View(),
		m.inputs.statusFieldsButton.View(),
		m.inputs.arnPrimaryKeyButton.View(),
	}

	renderedFieldNames := lipgloss.NewStyle().
		AlignHorizontal(lipgloss.Left).
		PaddingLeft(1).
		BorderRight(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderRightForeground(lipgloss.Color("240")).
		Render(strings.Join(fieldNames, "\n"))

	return lipgloss.JoinHorizontal(lipgloss.Right, renderedFieldNames, strings.Join(fieldViews, "\n"))
}

func (m ResourceForm) Init() tea.Cmd {
	return nil
}

type UpdateResourceConfig struct {
	Kind   string
	Config *ackconfig.ResourceConfig
}

type ReturnMessage struct{}

type OpenSpecFieldsMessage struct{}

type OpenStatusFieldsMessage struct{}

type resourceFormKeyMap struct {
	LineUp   key.Binding
	LineDown key.Binding
	Select   key.Binding
	Quit     key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k resourceFormKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Select, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k resourceFormKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Select}, // first column
		{k.Quit},   // second column
	}
}

var resourceFormKeys = resourceFormKeyMap{
	LineUp: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	LineDown: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select/toggle"),
	),
	Quit: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "go back"),
	),
}
