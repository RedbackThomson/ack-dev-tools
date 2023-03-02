package views

import (
	"fmt"

	"github.com/aws-controllers-k8s/dev-tools/pkg/generate/styles"
)

func HeaderView(controllerName string) string {
	title := styles.TitleStyle.Render(fmt.Sprintf("Generating %q", controllerName))
	return title
}
