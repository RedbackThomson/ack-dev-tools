package generate

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/samber/lo"

	ackconfig "github.com/aws-controllers-k8s/code-generator/pkg/config"
	ackmodel "github.com/aws-controllers-k8s/code-generator/pkg/model"
	"github.com/aws-controllers-k8s/dev-tools/pkg/generate/constants"
	"github.com/aws-controllers-k8s/dev-tools/pkg/generate/styles"
	"github.com/aws-controllers-k8s/dev-tools/pkg/generate/views"
)

func (w *Wizard) upsertResourceConfig(kind string) *ackconfig.ResourceConfig {
	config, ok := w.config.Resources[kind]
	if !ok {
		w.config.Resources[kind] = ackconfig.ResourceConfig{}
		config, _ = w.config.Resources[kind]
	}

	return &config
}

func (w *Wizard) getCRDByKind(kind string) *ackmodel.CRD {
	crd, _ := lo.Find(w.crds, func(crd *ackmodel.CRD) bool {
		return crd.Kind == kind
	})
	return crd
}

func (w Wizard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case views.ReturnMessage:
		w.state = resourcesSummary
	case views.SelectResource:
		w.selectedResourceKind = msg.ResourceKind
		w.state = resourceDetails
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, constants.Keymap.Quit):
			w.quitting = true
			return w, tea.Quit
		}
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(views.HeaderView(w.config.ModelName))
		constants.WindowSize = msg
		// Subtract 2 for the borders
		constants.ContainerViewSize = tea.WindowSizeMsg{Height: msg.Height - 2 - headerHeight, Width: msg.Width - 2}
		constants.UsableViewSize = tea.WindowSizeMsg{Height: msg.Height - 4 - headerHeight, Width: msg.Width - 4}
	}

	switch w.state {
	case resourcesSummary:
		newResourcesList, newCmd := w.resourceTable.Update(msg)
		resourcesList, ok := newResourcesList.(views.ResourceTable)
		if !ok {
			panic(fmt.Errorf(ErrAssertUpdate, "ResourceList"))
		}
		w.resourceTable = resourcesList
		cmd = newCmd
	case resourceDetails:
		w.selectedResourceForm = *views.NewResourceForm(w.getCRDByKind(w.selectedResourceKind), w.upsertResourceConfig(w.selectedResourceKind))
		newSelectedResourceForm, newCmd := w.selectedResourceForm.Update(msg)
		selectedResourceForm, ok := newSelectedResourceForm.(views.ResourceForm)
		if !ok {
			panic(fmt.Errorf(ErrAssertUpdate, "ResourceForm"))
		}
		w.selectedResourceForm = selectedResourceForm
		cmd = newCmd
	}

	cmds = append(cmds, cmd)

	return w, tea.Batch(cmds...)
}

func (w Wizard) View() string {
	if w.quitting {
		return ""
	}

	header := views.HeaderView(fmt.Sprintf("%s-controller", w.config.ModelName))

	var view string

	switch w.state {
	case resourceDetails:
		view = w.selectedResourceForm.View()
	case resourcesSummary:
		fallthrough
	default:
		view = w.resourceTable.View()
	}

	view = styles.FocusedStyle.
		Width(constants.ContainerViewSize.Width).
		Height(constants.ContainerViewSize.Height).
		Render(view)

	return lipgloss.JoinVertical(lipgloss.Top, header, view)
}

func (w Wizard) Init() tea.Cmd {
	return nil
}
