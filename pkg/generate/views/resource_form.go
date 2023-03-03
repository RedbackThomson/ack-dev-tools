package views

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/samber/lo"

	ackconfig "github.com/aws-controllers-k8s/code-generator/pkg/config"
	ackmodel "github.com/aws-controllers-k8s/code-generator/pkg/model"
	"github.com/aws-controllers-k8s/dev-tools/pkg/generate/components/button"
	"github.com/aws-controllers-k8s/dev-tools/pkg/generate/styles"
	"github.com/aws-controllers-k8s/dev-tools/pkg/generate/utils"
)

var (
	SpecFieldsButtonID   string = "spec"
	StatusFieldsButtonID string = "status"
)

type resourceFormKeyMap struct {
	LineUp   key.Binding
	LineDown key.Binding
	Quit     key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k resourceFormKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k resourceFormKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.LineUp, k.LineDown}, // first column
		{k.Quit},               // second column
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
	Quit: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "go back"),
	),
}

type resourceFormInputs struct {
	specFieldsButton   button.Model
	statusFieldsButton button.Model
}

type ResourceForm struct {
	crd    *ackmodel.CRD
	config *ackconfig.ResourceConfig

	inputs *resourceFormInputs
}

func (m ResourceForm) Keymap() help.KeyMap {
	return resourceFormKeys
}

func NewResourceForm(crd *ackmodel.CRD, config *ackconfig.ResourceConfig) *ResourceForm {
	form := &ResourceForm{
		crd:    crd,
		config: config,

		inputs: &resourceFormInputs{
			specFieldsButton:   button.New(SpecFieldsButtonID, fmt.Sprintf("%d fields", len(crd.SpecFields))),
			statusFieldsButton: button.New(StatusFieldsButtonID, fmt.Sprintf("%d fields", len(crd.StatusFields))),
		},
	}

	return form
}

func (m *ResourceForm) getInputFocusOrder() []utils.Focusable {
	return []utils.Focusable{
		&m.inputs.specFieldsButton,
		&m.inputs.statusFieldsButton,
	}
}

func (m *ResourceForm) rotateFocus(rotateDown bool) {
	focusOrder := m.getInputFocusOrder()
	current, currentIdx, exists := lo.FindIndexOf(focusOrder, func(item utils.Focusable) bool {
		return item.Focused()
	})

	if !exists {
		focusOrder[0].Focus()
		return
	}

	nextIndex := lo.Clamp(lo.Ternary(rotateDown, currentIdx+1, currentIdx-1), 0, len(focusOrder)-1)

	current.Blur()
	focusOrder[nextIndex].Focus()
}

func (m *ResourceForm) handleInputUpdates(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, resourceFormKeys.LineDown):
			m.rotateFocus(true)
		case key.Matches(msg, resourceFormKeys.LineUp):
			m.rotateFocus(false)
		case key.Matches(msg, resourceFormKeys.Quit):
			return *m, SendReturn
		}
	}

	switch {
	case m.inputs.specFieldsButton.Focused():
		m.inputs.specFieldsButton, cmd = m.inputs.specFieldsButton.Update(msg)
	case m.inputs.statusFieldsButton.Focused():
		m.inputs.statusFieldsButton, cmd = m.inputs.statusFieldsButton.Update(msg)
	}

	return *m, cmd
}

func (m *ResourceForm) getHeaderView() string {
	return styles.HeaderStyle.Render(m.crd.Kind)
}

func (m ResourceForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case button.ButtonSelectMessage:
		switch msg.GetID() {
		case SpecFieldsButtonID:
			log.Default().Println(msg.GetID())
		case StatusFieldsButtonID:
			log.Default().Println(msg.GetID())
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
		fmt.Sprintf("%t", m.crd.IsARNPrimaryKey()),
	}

	renderedFieldNames := lipgloss.NewStyle().
		AlignHorizontal(lipgloss.Right).
		BorderRight(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderRightForeground(lipgloss.Color("240")).
		Render(strings.Join(fieldNames, "\n"))

	content := lipgloss.JoinHorizontal(lipgloss.Right, renderedFieldNames, strings.Join(fieldViews, "\n"))
	return lipgloss.JoinVertical(lipgloss.Top, m.getHeaderView(), content)
}

func (m ResourceForm) Init() tea.Cmd {
	return nil
}

type ReturnMessage struct{}

func SendReturn() tea.Msg {
	return ReturnMessage{}
}
