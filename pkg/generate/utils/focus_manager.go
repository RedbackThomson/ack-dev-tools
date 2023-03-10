package utils

import "github.com/samber/lo"

type FocusRotateDirection int

const (
	FocusRotateUp FocusRotateDirection = iota
	FocusRotateDown
)

func RotateFocus(focusOrder []Focusable, direction FocusRotateDirection) {
	current, currentIdx, exists := lo.FindIndexOf(focusOrder, func(item Focusable) bool {
		return item.Focused()
	})

	if !exists {
		focusOrder[0].Focus()
		return
	}

	nextIndex := lo.Clamp(lo.Ternary(direction == FocusRotateDown, currentIdx+1, currentIdx-1), 0, len(focusOrder)-1)

	current.Blur()
	focusOrder[nextIndex].Focus()
}
