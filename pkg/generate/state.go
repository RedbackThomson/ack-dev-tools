package generate

import (
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"

	ackconfig "github.com/aws-controllers-k8s/code-generator/pkg/config"
	ackmodel "github.com/aws-controllers-k8s/code-generator/pkg/model"
	"github.com/aws-controllers-k8s/dev-tools/pkg/generate/utils"
	"github.com/aws-controllers-k8s/dev-tools/pkg/generate/views"
)

var models []tea.Model

// sessionState is used to track which model is focused
type sessionState uint

const (
	resourcesSummary sessionState = iota
	resourceDetails
	fieldsView
	fieldReferencesView
)

type Wizard struct {
	config  *ackconfig.Config
	service string
	model   *ackmodel.Model
	crds    []*ackmodel.CRD

	ready    bool
	quitting bool
	state    sessionState

	breadcrumbs *utils.Breadcrumbs

	help help.Model

	resourceTable views.ResourceTable

	selectedResourceForm views.ResourceForm
	fieldTable           views.FieldTable

	fieldReferencesForm views.ReferencesForm
}

func (w Wizard) Config() *ackconfig.Config {
	return w.config
}

func InitialState(config *ackconfig.Config, model *ackmodel.Model, service, modelName, apiVersion string) (Wizard, error) {
	crds, err := model.GetCRDs()
	if err != nil {
		return Wizard{}, err
	}

	w := Wizard{
		service:       service,
		config:        config,
		model:         model,
		crds:          crds,
		state:         resourcesSummary,
		resourceTable: *views.NewResourceTable(service, crds, config),
		ready:         false,
		help:          help.New(),
		breadcrumbs:   utils.NewBreadcrumbs(),
	}

	return w, nil
}
