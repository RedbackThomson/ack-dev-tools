package views

import (
	"fmt"
	"math"

	"github.com/charmbracelet/bubbles/help"
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

type resourceTableKeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Ignore key.Binding
	Help   key.Binding
	Quit   key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k resourceTableKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Ignore, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k resourceTableKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Select},   // first column
		{k.Ignore, k.Help, k.Quit}, // second column
	}
}

var resourceTableKeys = resourceTableKeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type ResourceTable struct {
	crds   []*ackmodel.CRD
	config *ackconfig.Config

	loaded bool

	table table.Model
}

func (m *ResourceTable) initialiseResourcesTable() error {
	headerHeight := lipgloss.Height(m.getHeaderView())

	width := constants.UsableViewSize.Width
	height := constants.UsableViewSize.Height - headerHeight

	columns := []table.Column{
		{Title: "Ignored", Width: 10},
		{Title: "Kind", Width: (int)(math.Round((float64)(width-10) / 3))},
		{Title: "# Spec Fields", Width: (int)(math.Round((float64)(width-10) / 3))},
		{Title: "# Status Fields", Width: (int)(math.Round((float64)(width-10) / 3))},
	}

	rows := lo.Map(m.crds, func(crd *ackmodel.CRD, index int) table.Row {
		ignored := lo.Contains(m.config.Ignore.ResourceNames, crd.Kind)

		return table.Row{
			lo.Ternary(ignored, "✓", ""),
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
		table.WithWidth(width), // Subtract because of cell padding
	)
	t.SetStyles(styles.DefaultTableStyle)

	m.table = t

	return nil
}

func (m *ResourceTable) getHeaderView() string {
	return styles.HeaderStyle.Render("Resources")
}

func (m ResourceTable) Keymap() help.KeyMap {
	return resourceTableKeys
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
	selectedKind := selectedItem[1]

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
		case key.Matches(msg, resourceTableKeys.Quit):
			return m, tea.Quit
		case key.Matches(msg, resourceTableKeys.Select):
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
