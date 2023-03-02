package views

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/samber/lo"

	ackconfig "github.com/aws-controllers-k8s/code-generator/pkg/config"
	ackmodel "github.com/aws-controllers-k8s/code-generator/pkg/model"
	"github.com/aws-controllers-k8s/dev-tools/pkg/generate/constants"
	"github.com/aws-controllers-k8s/dev-tools/pkg/generate/styles"
)

type ResourceTable struct {
	crds   []*ackmodel.CRD
	config *ackconfig.Config

	loaded bool

	table table.Model
}

func (m *ResourceTable) initialiseResourcesTable() error {
	headerHeight := lipgloss.Height(m.getHeaderView())

	// Subtract all offsets for unknown reasons
	// TODO: Figure out why offests are necessary
	width := constants.UsableViewSize.Width - 3
	height := constants.UsableViewSize.Height - headerHeight

	columns := []table.Column{
		{Title: "Kind", Width: width / 3},
		{Title: "# Spec Fields", Width: width / 3},
		{Title: "# Status Fields", Width: width / 3},
	}

	rows := lo.Map(m.crds, func(crd *ackmodel.CRD, index int) table.Row {
		return table.Row{
			crd.Names.Camel,
			fmt.Sprintf("%d", len(crd.SpecFields)),
			fmt.Sprintf("%d", len(crd.StatusFields)),
		}
	})

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(height),
		table.WithWidth(width),
	)
	t.SetStyles(styles.DefaultTableStyle)

	m.table = t

	return nil
}

func (m *ResourceTable) getHeaderView() string {
	return styles.HeaderStyle.Render("Resources")
}

func NewResourceTable(crds []*ackmodel.CRD, config *ackconfig.Config) *ResourceTable {
	form := &ResourceTable{
		crds:   crds,
		config: config,
		loaded: false,
	}

	return form
}

func (m *ResourceTable) SelectCurrentResource() tea.Msg {
	selectedItem := m.table.SelectedRow()
	selectedKind := selectedItem[0]

	return SelectResource{ResourceKind: selectedKind}
}

func (m ResourceTable) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.loaded {
		m.initialiseResourcesTable()
		m.loaded = true
	}

	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.initialiseResourcesTable()
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, constants.Keymap.Quit):
			return m, tea.Quit
		case key.Matches(msg, constants.Keymap.Enter):
			return m, m.SelectCurrentResource
		default:
			m.table, cmd = m.table.Update(msg)
		}
	}
	return m, cmd
}

func (m ResourceTable) View() string {
	return lipgloss.JoinVertical(lipgloss.Top, m.getHeaderView(), m.table.View())
}

func (m ResourceTable) Init() tea.Cmd {
	return nil
}

type ResourceListItem struct {
	list.DefaultItem

	kind            string
	numSpecFields   int
	numStatusFields int
}

func NewResourceItem(kind string, numSpecFields, numStatusFields int) ResourceListItem {
	return ResourceListItem{kind: kind, numSpecFields: numSpecFields, numStatusFields: numStatusFields}
}

// implement the list.Item interface
func (i ResourceListItem) FilterValue() string {
	return i.kind
}

func (i ResourceListItem) Title() string {
	return i.kind
}

func (i ResourceListItem) Description() string {
	return fmt.Sprintf("Fields: %d|%d", i.numSpecFields, i.numStatusFields)
}

type SelectResource struct {
	ResourceKind string
}
