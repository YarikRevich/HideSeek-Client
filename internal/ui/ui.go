package ui

import (
	"github.com/YarikRevich/HideSeek-Client/internal/player_mechanics/state_machine"
	"github.com/YarikRevich/HideSeek-Client/internal/player_mechanics/state_machine/constants/ui"
	"github.com/YarikRevich/HideSeek-Client/internal/profiling"
	"github.com/YarikRevich/HideSeek-Client/internal/ui/debug"
	"github.com/YarikRevich/HideSeek-Client/internal/ui/game"
	"github.com/YarikRevich/HideSeek-Client/internal/ui/pop_up_messages"
	settingsmenu "github.com/YarikRevich/HideSeek-Client/internal/ui/settings_menu"
	startmenu "github.com/YarikRevich/HideSeek-Client/internal/ui/start_menu"
	waitroom "github.com/YarikRevich/HideSeek-Client/internal/ui/wait_room"
	"github.com/YarikRevich/HideSeek-Client/tools/cli"
)

func Process() {
	profiling.UseProfiler().StartMonitoring(profiling.UI)
	
	switch statemachine.UseStateMachine().UI().GetState() {
	case ui.GAME:
		game.Draw()
	case ui.START_MENU:
		startmenu.Draw()
	case ui.SETTINGS_MENU:
		settingsmenu.Draw()
	case ui.CREATE_LOBBY_MENU:
	case ui.JOIN_LOBBY_MENU:

	case ui.CHOOSE_EQUIPMENT:
	case ui.WAIT_ROOM:
		waitroom.Draw()
	}

	if cli.GetDebug() {
		debug.Draw()
	}
	popupmessages.Draw()

	profiling.UseProfiler().EndMonitoring()
}
