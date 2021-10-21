package game

import (
	"fmt"

	"github.com/YarikRevich/HideSeek-Client/internal/direction"
	"github.com/YarikRevich/HideSeek-Client/internal/gameplay/pc"
	"github.com/YarikRevich/HideSeek-Client/internal/history"
	"github.com/YarikRevich/HideSeek-Client/internal/physics/jump"
	"github.com/hajimehoshi/ebiten/v2"
)

func Exec() {

	for _, v := range ebiten.GamepadIDs() {
		fmt.Println(v)
		fmt.Println(ebiten.IsGamepadButtonPressed(v, ebiten.GamepadButton0), "0")
		fmt.Println(ebiten.IsGamepadButtonPressed(v, ebiten.GamepadButton1), "1")
		fmt.Println(ebiten.IsGamepadButtonPressed(v, ebiten.GamepadButton2), "2")
		fmt.Println(ebiten.IsGamepadButtonPressed(v, ebiten.GamepadButton3), "3")
		fmt.Println(ebiten.IsGamepadButtonPressed(v, ebiten.GamepadButton4), "4")
		fmt.Println(ebiten.IsGamepadButtonPressed(v, ebiten.GamepadButton5), "5")
		fmt.Println(ebiten.IsGamepadButtonPressed(v, ebiten.GamepadButton6), "6")
		fmt.Println(ebiten.IsGamepadButtonPressed(v, ebiten.GamepadButton7), "7")
		fmt.Println(ebiten.IsGamepadButtonPressed(v, ebiten.GamepadButton8), "8")
		fmt.Println(ebiten.IsGamepadButtonPressed(v, ebiten.GamepadButton9), "9")
		fmt.Println(ebiten.IsGamepadButtonPressed(v, ebiten.GamepadButton10), "10")
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) || isGamepadButtonPressed(gamepadUPButton) {
		history.SetDirection(direction.UP)
		pc.UsePC().SetY(pc.UsePC().Y - pc.UsePC().Buffs.SpeedY)
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) || isGamepadButtonPressed(gamepadDOWNButton)  {
		history.SetDirection(direction.DOWN)
		pc.UsePC().SetY(pc.UsePC().Y + pc.UsePC().Buffs.SpeedY)
	}

	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) || isGamepadButtonPressed(gamepadRIGHTButton)  {
		history.SetDirection(direction.RIGHT)
		pc.UsePC().SetX(pc.UsePC().X + pc.UsePC().Buffs.SpeedX)
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || isGamepadButtonPressed(gamepadLEFTButton)  {
		history.SetDirection(direction.LEFT)
		pc.UsePC().SetX(pc.UsePC().X - pc.UsePC().Buffs.SpeedX)
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) || isGamepadButtonPressed(gamepadRIGHTUPPERCLICKERButton) {
		jump.CalculateJump(pc.UsePC())
	}

	// 	currPosition := pixel.V(float64(g.userConfig.Pos.X), float64(g.userConfig.Pos.Y))
	// 	g.mapComponents.GetCollisions().GetDoorsCollisions().DoorTraker(currPosition)

	// 	switch {
	// 	case ((g.winConf.Win.Pressed(pixelgl.KeyW) ||
	// 		g.winConf.Win.Pressed(pixelgl.KeyA) ||
	// 		g.winConf.Win.Pressed(pixelgl.KeyS) ||
	// 		g.winConf.Win.Pressed(pixelgl.KeyD)) && g.winConf.Win.JustPressed(pixelgl.KeySpace)) ||
	// 		g.winConf.Win.JustPressed(pixelgl.KeySpace):
	// 		ok, user := g.mapComponents.GetCollisions().GetHeroCollisions().IsHero(currPosition, g.winConf.GameProcess.OtherUsers).Near(30, 37)
	// 		if ok {
	// 			g.userConfig.Context.Additional = append(g.userConfig.Context.Additional, user, "1")
	// 			parser := Server.GameParser(new(Server.GameRequest))
	// 			server := Server.Network(new(Server.N))
	// 			server.Init(nil, g.userConfig, 1, nil, parser.Parse, "UpdateUsersHealth")

	// 			server.Write()
	// 			server.ReadGame(parser.Unparse)

	// 			g.userConfig.Context.Additional = []string{}
	// 		}
	// 		fallthrough
	// 	default:
	// 		switch {
	// 		case g.winConf.Win.Pressed(pixelgl.KeyW):
	// 			coll := g.mapComponents.GetCollisions().IsCritical(pixel.V(float64(g.userConfig.Pos.X), float64(g.userConfig.Pos.Y+2)), g.winConf.GameProcess.OtherUsers, "top")
	// 			if coll {
	// 				return
	// 			}

	// 			if g.userConfig.Pos.Y <= g.mapComponents.GetHeroBorder().Top() {
	// 				g.userConfig.Pos.Y += 3
	// 			}
	// 			if g.winConf.Cam.CamPos.Y < g.mapComponents.GetCamBorder().Top()+50 {
	// 				if g.userConfig.Pos.Y >= int(g.winConf.Win.Bounds().Center().Y) {
	// 					g.winConf.Cam.CamPos.Y += 5
	// 				}
	// 			}
	// 		case g.winConf.Win.Pressed(pixelgl.KeyA):
	// 			coll := g.mapComponents.GetCollisions().IsCritical(pixel.V(float64(g.userConfig.Pos.X-2), float64(g.userConfig.Pos.Y)), g.winConf.GameProcess.OtherUsers, "left")
	// 			if coll {
	// 				return
	// 			}

	// 			if g.userConfig.Pos.X >= g.mapComponents.GetHeroBorder().Left() {
	// 				g.userConfig.Pos.X -= 3
	// 			}
	// 			if g.winConf.Cam.CamPos.X >= g.mapComponents.GetCamBorder().Left() {
	// 				if g.userConfig.Pos.X <= int(g.winConf.Win.Bounds().Center().X) {
	// 					g.winConf.Cam.CamPos.X -= 5
	// 				}
	// 			}
	// 		case g.winConf.Win.Pressed(pixelgl.KeyS):
	// 			coll := g.mapComponents.GetCollisions().IsCritical(pixel.V(float64(g.userConfig.Pos.X), float64(g.userConfig.Pos.Y-2)), g.winConf.GameProcess.OtherUsers, "bottom")
	// 			if coll {
	// 				return
	// 			}

	// 			if g.userConfig.Pos.Y >= g.mapComponents.GetHeroBorder().Bottom() {
	// 				g.userConfig.Pos.Y -= 3
	// 			}
	// 			if g.winConf.Cam.CamPos.Y >= g.mapComponents.GetCamBorder().Bottom() {
	// 				if g.userConfig.Pos.Y <= int(g.winConf.Win.Bounds().Center().Y) {
	// 					g.winConf.Cam.CamPos.Y -= 5
	// 				}
	// 			}
	// 		case g.winConf.Win.Pressed(pixelgl.KeyD):
	// 			coll := g.mapComponents.GetCollisions().IsCritical(pixel.V(float64(g.userConfig.Pos.X+2), float64(g.userConfig.Pos.Y)), g.winConf.GameProcess.OtherUsers, "right")
	// 			if coll {
	// 				return
	// 			}

	// 			if g.userConfig.Pos.X <= g.mapComponents.GetHeroBorder().Right() {
	// 				g.userConfig.Pos.X += 3
	// 			}
	// 			if g.winConf.Cam.CamPos.X <= g.mapComponents.GetCamBorder().Right() {
	// 				if g.userConfig.Pos.X >= int(g.winConf.Win.Bounds().Center().X) {
	// 					g.winConf.Cam.CamPos.X += 5
	// 				}
	// 			}
	// 		}
	// 	}
	// }
}
