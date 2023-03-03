package views

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/samber/lo"

	ackconfig "github.com/aws-controllers-k8s/code-generator/pkg/config"
	ackmodel "github.com/aws-controllers-k8s/code-generator/pkg/model"
	"github.com/aws-controllers-k8s/dev-tools/pkg/generate/constants"
	"github.com/aws-controllers-k8s/dev-tools/pkg/generate/styles"
)

type FieldTableType string

var (
	FieldTableTypeSpec   FieldTableType = "Spec"
	FieldTableTypeStatus FieldTableType = "Status"
)

type FieldTable struct {
	tableType FieldTableType
	fields    map[string]*ackmodel.Field
	config    *ackconfig.ResourceConfig

	loaded bool

	table table.Model
}

func (m *FieldTable) initialiseFieldsTable() error {
	headerHeight := lipgloss.Height(m.getHeaderView())

	width := constants.UsableViewSize.Width
	height := constants.UsableViewSize.Height - headerHeight

	columns := []table.Column{
		{Title: "Name", Width: width},
	}

	rows := lo.MapToSlice(m.fields, func(key string, field *ackmodel.Field) table.Row {
		return table.Row{
			field.Names.Camel,
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

func (m *FieldTable) getHeaderView() string {
	return styles.HeaderStyle.Render(fmt.Sprintf("%s Fields", m.tableType))
}

func (m FieldTable) Keymap() help.KeyMap {
	return fieldTableKeys
}

func NewFieldTable(tableType FieldTableType, fields map[string]*ackmodel.Field, config *ackconfig.ResourceConfig) *FieldTable {
	form := &FieldTable{
		tableType: tableType,
		fields:    fields,
		config:    config,
		loaded:    false,
	}

	return form
}

func (m FieldTable) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.loaded {
		m.initialiseFieldsTable()
		m.loaded = true
	}

	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.initialiseFieldsTable()
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, fieldTableKeys.Quit):
			return m, (func() tea.Msg {
				return ReturnMessage{}
			})
		case key.Matches(msg, fieldTableKeys.Select):
			return m, nil
		default:
			m.table, cmd = m.table.Update(msg)
		}
	}
	return m, cmd
}

func (m FieldTable) View() string {
	return lipgloss.JoinVertical(lipgloss.Top, m.getHeaderView(), m.table.View())
}

func (m FieldTable) Init() tea.Cmd {
	return nil
}

type fieldTableKeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Ignore key.Binding
	Help   key.Binding
	Quit   key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k fieldTableKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Ignore, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k fieldTableKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Select},
		{k.Ignore, k.Help, k.Quit},
	}
}

var fieldTableKeys = fieldTableKeyMap{
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
