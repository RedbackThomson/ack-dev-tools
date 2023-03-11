package views

import (
	"sort"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
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

func (m *FieldTable) getFieldConfig(fieldName string) *ackconfig.FieldConfig {
	if m.config.Fields == nil {
		return nil
	}
	config, exists := m.config.Fields[fieldName]
	if !exists {
		return nil
	}

	return config
}

func (m *FieldTable) InitialiseFieldsTable() error {
	style := styles.DefaultTableStyle

	// subtract padding for every header cell
	width := constants.UsableViewSize.Width - style.Header.GetHorizontalPadding()*7
	height := constants.UsableViewSize.Height

	nameColumnWidth := 30
	boolColumnWidth := (width - nameColumnWidth) / 6

	columns := []table.Column{
		{Title: "Name", Width: nameColumnWidth},
		{Title: "Is Required", Width: boolColumnWidth},
		{Title: "Is Primary Key", Width: boolColumnWidth},
		{Title: "Is Secret", Width: boolColumnWidth},
		{Title: "Is Immutable", Width: boolColumnWidth},
		{Title: "Is ARN", Width: boolColumnWidth},
		{Title: "References", Width: boolColumnWidth},
	}

	rows := lo.MapToSlice(m.fields, func(key string, field *ackmodel.Field) table.Row {
		fieldConfig := m.getFieldConfig(field.Names.Camel)

		return table.Row{
			field.Names.Camel,
			lo.Ternary(fieldConfig != nil && fieldConfig.IsRequired != nil && *fieldConfig.IsRequired, " ✓", ""),
			lo.Ternary(fieldConfig != nil && fieldConfig.IsPrimaryKey, " ✓", ""),
			lo.Ternary(fieldConfig != nil && fieldConfig.IsSecret, " ✓", ""),
			lo.Ternary(fieldConfig != nil && fieldConfig.IsImmutable, " ✓", ""),
			lo.Ternary(fieldConfig != nil && fieldConfig.IsARN, " ✓", ""),
			lo.Ternary(fieldConfig != nil && fieldConfig.References != nil, " ✓", ""),
		}
	})
	sort.Slice(rows, func(i, j int) bool {
		return rows[i][0] < rows[j][0]
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

func (m *FieldTable) showReferencesForm() tea.Msg {
	selectedItem := m.table.SelectedRow()
	selectedFieldName := selectedItem[0]

	return OpenFieldReferences{
		FieldName: selectedFieldName,
	}
}

func (m FieldTable) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.loaded {
		m.InitialiseFieldsTable()
		m.loaded = true
	}

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.InitialiseFieldsTable()
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, fieldTableKeys.References):
			return m, m.showReferencesForm
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
	return m.table.View()
}

func (m FieldTable) Init() tea.Cmd {
	return nil
}

type OpenFieldReferences struct {
	FieldName string
}

type fieldTableKeyMap struct {
	Up         key.Binding
	Down       key.Binding
	Select     key.Binding
	Ignore     key.Binding
	Help       key.Binding
	References key.Binding
	Quit       key.Binding
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
		key.WithKeys("esc"),
		key.WithHelp("esc", "go back"),
	),
	References: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "set reference"),
	),
}
