package utils

import "github.com/samber/lo"

type Breadcrumbs struct {
	parts []string
}

func (b *Breadcrumbs) Push(val string) {
	b.parts = append(b.parts, val)
}

func (b *Breadcrumbs) Pop() {
	if len(b.parts) > 0 {
		b.parts = b.parts[:len(b.parts)-1]
	}
}

func (b *Breadcrumbs) ReplaceAt(val string, index int) {
	newParts := lo.Slice(b.parts, 0, index)
	newParts = append(newParts, val)
	newParts = append(newParts, lo.Slice(b.parts, index+1, len(b.parts))...)

	b.parts = newParts
}

func (b *Breadcrumbs) Size() int {
	return len(b.parts)
}

func (b *Breadcrumbs) Parts() []string {
	return b.parts
}

func NewBreadcrumbs() *Breadcrumbs {
	return &Breadcrumbs{}
}
