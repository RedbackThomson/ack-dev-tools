package views

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/samber/lo"

	ackconfig "github.com/aws-controllers-k8s/code-generator/pkg/config"
	ackmodel "github.com/aws-controllers-k8s/code-generator/pkg/model"
	"github.com/aws-controllers-k8s/dev-tools/pkg/generate/constants"
	"github.com/aws-controllers-k8s/dev-tools/pkg/generate/styles"
)

type ResourceTable struct {
	service string
	crds    []*ackmodel.CRD
	config  *ackconfig.Config

	loaded bool

	table table.Model
}

func (m *ResourceTable) initialiseResourcesTable() {
	style := styles.DefaultTableStyle

	// subtract padding for every header cell
	width := constants.UsableViewSize.Width - style.Header.GetHorizontalPadding()*4
	height := constants.UsableViewSize.Height

	columns := []table.Column{
		{Title: "Ignored", Width: 7},
		{Title: "Kind", Width: (int)(math.Round((float64)(width-7) / 2))},
		{Title: "API Group", Width: (int)(math.Round((float64)(width-7) / 2))},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(height),
		table.WithWidth(width), // Subtract because of cell padding
	)
	t.SetStyles(style)

	m.table = t

	m.loadTableRows()
}

func (m *ResourceTable) loadTableRows() {
	rows := lo.Map(m.crds, func(crd *ackmodel.CRD, index int) table.Row {
		ignored := lo.Contains(m.config.Ignore.ResourceNames, crd.Kind)

		return table.Row{
			lo.Ternary(ignored, " ✓", ""),
			crd.Names.Camel,
			strings.ToLower(fmt.Sprintf("%s.%s.services.k8s.aws", crd.Kind, m.service)),
		}
	})
	m.table.SetRows(rows)
}

func (m ResourceTable) Keymap() help.KeyMap {
	return resourceTableKeys
}

func NewResourceTable(service string, crds []*ackmodel.CRD, config *ackconfig.Config) *ResourceTable {
	form := &ResourceTable{
		service: service,
		crds:    crds,
		config:  config,
		loaded:  false,
	}

	return form
}

func (m *ResourceTable) selectCurrentResource() tea.Msg {
	selectedItem := m.table.SelectedRow()
	selectedKind := selectedItem[1]

	return SelectResource{ResourceKind: selectedKind}
}

func (m *ResourceTable) toggleIgnoreCurrentResource() {
	selectedItem := m.table.SelectedRow()
	selectedKind := selectedItem[1]

	isIgnored := lo.Contains(m.config.Ignore.ResourceNames, selectedKind)
	if isIgnored {
		m.config.Ignore.ResourceNames = lo.Filter(m.config.Ignore.ResourceNames, func(name string, index int) bool {
			return name != selectedKind
		})
	} else {
		m.config.Ignore.ResourceNames = append(m.config.Ignore.ResourceNames, selectedKind)
	}
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
		case key.Matches(msg, resourceTableKeys.Select):
			return m, m.selectCurrentResource
		case key.Matches(msg, resourceTableKeys.Ignore):
			m.toggleIgnoreCurrentResource()
			m.loadTableRows()
			return m, nil
		case key.Matches(msg, resourceTableKeys.Quit):
			return m, tea.Quit
		default:
			m.table, cmd = m.table.Update(msg)
		}
	}
	return m, cmd
}

func (m ResourceTable) View() string {
	return m.table.View()
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
	return []key.Binding{k.Ignore, k.Select, k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k resourceTableKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Ignore, k.Select},           // first column
		{k.Up, k.Down, k.Help, k.Quit}, // second column
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
	Ignore: key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "toggle ignore"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q/esc/ctrl+c", "quit"),
	),
}
