package generate

import (
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"

	ackconfig "github.com/aws-controllers-k8s/code-generator/pkg/config"
	ackmodel "github.com/aws-controllers-k8s/code-generator/pkg/model"
	"github.com/aws-controllers-k8s/dev-tools/pkg/generate/views"
)

var models []tea.Model

// sessionState is used to track which model is focused
type sessionState uint

const (
	resourcesSummary sessionState = iota
	resourceDetails
	fieldsView
)

type Wizard struct {
	config *ackconfig.Config
	model  *ackmodel.Model
	crds   []*ackmodel.CRD

	ready    bool
	quitting bool
	state    sessionState

	help help.Model

	resourceTable views.ResourceTable

	selectedResourceForm views.ResourceForm
}

func (w Wizard) Config() *ackconfig.Config {
	return w.config
}

func InitialState(model *ackmodel.Model, modelName string) (Wizard, error) {
	config := &ackconfig.Config{
		ModelName: modelName,
		Resources: map[string]ackconfig.ResourceConfig{},
		Ignore: ackconfig.IgnoreSpec{
			ResourceNames: []string{},
		},
		Operations:                     map[string]ackconfig.OperationConfig{},
		PrefixConfig:                   ackconfig.PrefixConfig{},
		IncludeACKMetadata:             true,
		SetManyOutputNotFoundErrReturn: "",
	}

	crds, err := model.GetCRDs()
	if err != nil {
		return Wizard{}, err
	}

	w := Wizard{
		config:        config,
		model:         model,
		crds:          crds,
		state:         resourcesSummary,
		resourceTable: *views.NewResourceTable(crds, config),
		ready:         false,
		help:          help.New(),
	}

	return w, nil
}
