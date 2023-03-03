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

func (w Wizard) currentView() views.View {
	switch w.state {
	case resourceDetails:
		return w.selectedResourceForm
	case resourcesSummary:
		return w.resourceTable
	case fieldsView:
		return w.fieldTable
	default:
		panic(fmt.Errorf(ErrNoCurrentViewDefined, w.state))
	}
}

func (w *Wizard) replaceCurrentView(r tea.Model) {
	switch w.state {
	case resourceDetails:
		selectedResourceForm, ok := r.(views.ResourceForm)
		if !ok {
			panic(fmt.Errorf(ErrAssertUpdate, "ResourceForm"))
		}
		w.selectedResourceForm = selectedResourceForm
	case resourcesSummary:
		resourceTable, ok := r.(views.ResourceTable)
		if !ok {
			panic(fmt.Errorf(ErrAssertUpdate, "ResourceTable"))
		}
		w.resourceTable = resourceTable
	case fieldsView:
		fieldTable, ok := r.(views.FieldTable)
		if !ok {
			panic(fmt.Errorf(ErrAssertUpdate, "FieldTable"))
		}
		w.fieldTable = fieldTable
	default:
		panic(fmt.Errorf(ErrNoReplaceViewDefined, w.state))
	}
}

func (w Wizard) recalculateScreenSizes(view views.View, windowSize tea.WindowSizeMsg) {
	headerHeight := lipgloss.Height(views.HeaderView(w.config.ModelName))
	footerHeight := lipgloss.Height(views.FooterView(w.help, view.Keymap()))

	constants.WindowSize = windowSize
	// Subtract 2 for the borders
	constants.ContainerViewSize = tea.WindowSizeMsg{Height: windowSize.Height - 2 - headerHeight - footerHeight, Width: windowSize.Width - 2}
	constants.UsableViewSize = tea.WindowSizeMsg{Height: constants.ContainerViewSize.Height - 2, Width: constants.ContainerViewSize.Width - 2}
}

func (w Wizard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	view := w.currentView()

	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case views.ReturnMessage:
		switch w.state {
		case resourceDetails:
			w.state = resourcesSummary
		case fieldsView:
			w.state = resourceDetails
		}
		return w, nil
	case views.SelectResource:
		crd := w.getCRDByKind(msg.ResourceKind)
		w.selectedResourceForm = *views.NewResourceForm(crd, w.upsertResourceConfig(msg.ResourceKind))
		w.state = resourceDetails
	case views.OpenSpecFieldsMessage:
		crd := w.selectedResourceForm.CRD()
		w.fieldTable = *views.NewFieldTable(views.FieldTableTypeSpec, crd.SpecFields, w.upsertResourceConfig(crd.Kind))
		w.state = fieldsView
	case views.OpenStatusFieldsMessage:
		crd := w.selectedResourceForm.CRD()
		w.fieldTable = *views.NewFieldTable(views.FieldTableTypeStatus, crd.StatusFields, w.upsertResourceConfig(crd.Kind))
		w.state = fieldsView
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("q", "ctrl+c"))):
			w.quitting = true
			return w, tea.Quit
			// TODO: Resizing down when ShowAll = false breakes layout
			// case key.Matches(msg, key.NewBinding(key.WithKeys("?", "h"))):
			// 	w.help.ShowAll = !w.help.ShowAll
			// 	w.recalculateScreenSizes(view, constants.WindowSize)
			// 	return w, nil
		}
	case tea.WindowSizeMsg:
		w.recalculateScreenSizes(view, msg)
	}

	updatedView, newCmd := w.currentView().Update(msg)
	w.replaceCurrentView(updatedView)

	cmds = append(cmds, newCmd)

	return w, tea.Batch(cmds...)
}

func (w Wizard) View() string {
	if w.quitting {
		return ""
	}

	header := views.HeaderView(fmt.Sprintf("%s-controller", w.config.ModelName))

	view := w.currentView()

	renderedView := styles.FocusedStyle.
		Width(constants.ContainerViewSize.Width).
		Height(constants.ContainerViewSize.Height).
		Render(view.View())

	footer := views.FooterView(w.help, view.Keymap())

	return lipgloss.JoinVertical(lipgloss.Top, header, renderedView, footer)
}

func (w Wizard) Init() tea.Cmd {
	return nil
}
