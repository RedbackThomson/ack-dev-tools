package views

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/samber/lo"

	ackconfig "github.com/aws-controllers-k8s/code-generator/pkg/config"
	ackmodel "github.com/aws-controllers-k8s/code-generator/pkg/model"
	"github.com/aws-controllers-k8s/dev-tools/pkg/generate/components/button"
	"github.com/aws-controllers-k8s/dev-tools/pkg/generate/utils"
)

type referencesFormInputs struct {
	serviceName textinput.Model
	resource    textinput.Model
	path        textinput.Model
}

type ReferencesForm struct {
	service   string
	fieldName string
	crd       *ackmodel.CRD
	config    *ackconfig.ReferencesConfig

	inputs *referencesFormInputs
}

func (m ReferencesForm) CRD() *ackmodel.CRD {
	return m.crd
}

func (m ReferencesForm) Keymap() help.KeyMap {
	return referencesFormKeys
}

func NewReferencesForm(service string, fieldName string, crd *ackmodel.CRD, config *ackconfig.ReferencesConfig) *ReferencesForm {
	form := &ReferencesForm{
		service:   service,
		fieldName: fieldName,
		crd:       crd,
		config:    config,

		inputs: &referencesFormInputs{
			serviceName: textinput.New(),
			resource:    textinput.New(),
			path:        textinput.New(),
		},
	}

	form.inputs.serviceName.Placeholder = service

	form.inputs.serviceName.SetValue(config.ServiceName)
	form.inputs.resource.SetValue(config.Resource)
	form.inputs.path.SetValue(config.Path)

	return form
}

func (m ReferencesForm) getInputFocusOrder() []utils.Focusable {
	return []utils.Focusable{
		&m.inputs.serviceName,
		&m.inputs.resource,
		&m.inputs.path,
	}
}

func (m *ReferencesForm) handleInputUpdates(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, referencesFormKeys.LineDown):
			utils.RotateFocus(m.getInputFocusOrder(), utils.FocusRotateDown)
		case key.Matches(msg, referencesFormKeys.LineUp):
			utils.RotateFocus(m.getInputFocusOrder(), utils.FocusRotateUp)
		case key.Matches(msg, referencesFormKeys.Quit):
			return *m, (func() tea.Msg {
				return ReturnMessage{}
			})
		}
	}

	switch {
	case m.inputs.serviceName.Focused():
		m.inputs.serviceName, cmd = m.inputs.serviceName.Update(msg)
		val := m.inputs.serviceName.Value()
		m.config.ServiceName = lo.Ternary(lo.IsEmpty(val), m.service, val)
	case m.inputs.resource.Focused():
		m.inputs.resource, cmd = m.inputs.resource.Update(msg)
		m.config.Resource = m.inputs.resource.Value()
	case m.inputs.path.Focused():
		m.inputs.path, cmd = m.inputs.path.Update(msg)
		m.config.Path = m.inputs.path.Value()
	}

	return *m, cmd
}

func (m ReferencesForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case button.ButtonSelectMessage:
		switch msg.GetID() {
		}
	}

	return m.handleInputUpdates(msg)
}

func (m ReferencesForm) View() string {
	fieldNames := []string{
		"Service Name",
		"Resource",
		"Path",
	}

	fieldViews := []string{
		m.inputs.serviceName.View(),
		m.inputs.resource.View(),
		m.inputs.path.View(),
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

func (m ReferencesForm) Init() tea.Cmd {
	return nil
}

type UpdateFieldResourceConfig struct {
	Kind   string
	Field  string
	Config *ackconfig.ReferencesConfig
}

type referencesFormKeyMap struct {
	LineUp   key.Binding
	LineDown key.Binding
	Quit     key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k referencesFormKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k referencesFormKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Quit}, // first column
	}
}

var referencesFormKeys = referencesFormKeyMap{
	LineUp: key.NewBinding(
		key.WithKeys("up", "k", "shift+tab"),
		key.WithHelp("↑/k", "up"),
	),
	LineDown: key.NewBinding(
		key.WithKeys("down", "j", "tab"),
		key.WithHelp("↓/j", "down"),
	),
	Quit: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "save/go back"),
	),
}
