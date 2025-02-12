package game

import (
	"github.com/YarikRevich/hide-seek-client/internal/core/middlewares"
	"github.com/YarikRevich/hide-seek-client/internal/core/statemachine"
	"github.com/YarikRevich/hide-seek-client/internal/core/world"
)

func Show() {
	world.UseWorld().DebugInit()

	middlewares.UseMiddlewares().UI().UseAfter(func() {
		statemachine.UseStateMachine().UI().SetState(statemachine.UI_GAME)
	})

	statemachine.UseStateMachine().Input().SetState(statemachine.INPUT_GAME)
}
