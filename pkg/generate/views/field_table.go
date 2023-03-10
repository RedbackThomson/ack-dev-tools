package views

import (
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

	referencesDialog Dialog
}

func (m *FieldTable) initialiseFieldsTable() error {
	width := constants.UsableViewSize.Width
	height := constants.UsableViewSize.Height

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

func (m *FieldTable) listDialogs() []Dialog {
	return []Dialog{
		m.referencesDialog,
	}
}

func (m *FieldTable) showReferencesDialog() {
	selectedItem := m.table.SelectedRow()
	selectedFieldName := selectedItem[0]

	m.referencesDialog = *NewDialog("Add references to "+selectedFieldName, DefaultDialogWidth)
	m.referencesDialog.SetVisible(true)
}

func (m FieldTable) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.loaded {
		m.initialiseFieldsTable()
		m.loaded = true
	}

	var cmd tea.Cmd

	switch {
	case m.referencesDialog.IsVisible():
		m.referencesDialog, cmd = m.referencesDialog.Update(msg)
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.initialiseFieldsTable()
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, fieldTableKeys.References):
			m.showReferencesDialog()
		case key.Matches(msg, fieldTableKeys.Quit):
			openDialog, _, isOpen := lo.FindIndexOf(m.listDialogs(), func(item Dialog) bool {
				return item.IsVisible()
			})
			if !isOpen {
				return m, (func() tea.Msg {
					return ReturnMessage{}
				})
			}
			openDialog.SetVisible(false)

		case key.Matches(msg, fieldTableKeys.Select):
			return m, nil
		default:
			m.table, cmd = m.table.Update(msg)
		}
	}

	return m, cmd
}

func (m FieldTable) View() string {
	switch {
	case m.referencesDialog.IsVisible():
		return m.referencesDialog.View()
	default:
		return m.table.View()
	}
}

func (m FieldTable) Init() tea.Cmd {
	return nil
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
