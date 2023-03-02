package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	ackconfig "github.com/aws-controllers-k8s/code-generator/pkg/config"
	ackmodel "github.com/aws-controllers-k8s/code-generator/pkg/model"
	"github.com/aws-controllers-k8s/dev-tools/pkg/generate/components/button"
	"github.com/aws-controllers-k8s/dev-tools/pkg/generate/constants"
	"github.com/aws-controllers-k8s/dev-tools/pkg/generate/styles"
)

var (
	SpecFieldsButtonName   string = "spec"
	StatusFieldsButtonName string = "status"
)

type resourceFormKeyMap struct {
	LineUp   key.Binding
	LineDown key.Binding
}

var ResourceFormKeyMap = resourceFormKeyMap{
	LineUp: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	LineDown: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
}

type ResourceForm struct {
	crd    *ackmodel.CRD
	config *ackconfig.ResourceConfig

	specFieldsButton   button.Model
	statusFieldsButton button.Model
}

func NewResourceForm(crd *ackmodel.CRD, config *ackconfig.ResourceConfig) *ResourceForm {
	form := &ResourceForm{
		crd:    crd,
		config: config,

		specFieldsButton:   button.New(SpecFieldsButtonName, fmt.Sprintf("%d fields", len(crd.SpecFields))),
		statusFieldsButton: button.New(StatusFieldsButtonName, fmt.Sprintf("%d fields", len(crd.StatusFields))),
	}

	form.specFieldsButton.Focus()

	return form
}

func (m *ResourceForm) getHeaderView() string {
	return styles.HeaderStyle.Render(m.crd.Kind)
}

func (m ResourceForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case button.ButtonSelectMessage:
		switch msg.GetID() {
		case SpecFieldsButtonName:
			fmt.Printf("Pressed spec")
		case StatusFieldsButtonName:
			fmt.Printf("Pressed status")
		}

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, ResourceFormKeyMap.LineDown):
			fallthrough
		case key.Matches(msg, ResourceFormKeyMap.LineUp):
			if m.statusFieldsButton.Focused() {
				m.statusFieldsButton.Blur()
				m.specFieldsButton.Focus()
			} else {
				m.specFieldsButton.Blur()
				m.statusFieldsButton.Focus()
			}
		case key.Matches(msg, constants.Keymap.Back):
			return m, SendReturn
		}
	}
	switch {
	case m.specFieldsButton.Focused():
		m.specFieldsButton, cmd = m.specFieldsButton.Update(msg)
		return m, cmd
	case m.statusFieldsButton.Focused():
		m.statusFieldsButton, cmd = m.statusFieldsButton.Update(msg)
		return m, cmd
	}
	return m, cmd
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
		m.specFieldsButton.View(),
		m.statusFieldsButton.View(),
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
