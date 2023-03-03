package views

import (
	"fmt"

	"github.com/aws-controllers-k8s/dev-tools/pkg/generate/styles"
	"github.com/charmbracelet/bubbles/help"
)

func HeaderView(controllerName string) string {
	title := styles.TitleStyle.Render(fmt.Sprintf("Generating %q", controllerName))
	return title
}

func FooterView(help help.Model, keymap help.KeyMap) string {
	return help.View(keymap)
}
