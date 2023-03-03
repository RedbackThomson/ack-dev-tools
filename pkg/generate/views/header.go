package views

import (
	"fmt"

	"github.com/aws-controllers-k8s/dev-tools/pkg/generate/styles"
	"github.com/aws-controllers-k8s/dev-tools/pkg/generate/utils"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
	"github.com/samber/lo"
)

func HeaderView(service string, breadcrumbs *utils.Breadcrumbs) string {
	title := styles.TitleStyle.Render(fmt.Sprintf("Generating %q", fmt.Sprintf("%s-controller", service)))

	breadcrumbTitles := lo.Map(breadcrumbs.Parts(), func(part string, index int) string {
		return styles.BreadcrumbStyle.Render(part)
	})

	headerParts := []string{title}
	headerParts = append(headerParts, breadcrumbTitles...)

	return lipgloss.JoinHorizontal(lipgloss.Left, headerParts...)
}

func FooterView(help help.Model, keymap help.KeyMap) string {
	return help.View(keymap)
}
