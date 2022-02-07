package ui

import "github.com/YarikRevich/hide-seek-client/internal/core/screen"

type Component interface {
	Render(screen.ScreenManager)
}

type UIManager struct {
	components []Component
}

func (uim *UIManager) Clear() {
	uim.components = uim.components[:0]
}

func (uim *UIManager) AddComponent(c Component) {
	uim.components = append(uim.components, c)
}

func (uim *UIManager) Render(s screen.ScreenManager) {
	for _, c := range uim.components {
		c.Render(s)
	}
}
