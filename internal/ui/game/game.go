package game

import (
	// "github.com/YarikRevich/HideSeek-Client/internal/gameplay/pc"
	//

	"github.com/YarikRevich/HideSeek-Client/internal/direction"
	// "github.com/YarikRevich/HideSeek-Client/internal/gameplay/pc"
	"github.com/YarikRevich/HideSeek-Client/internal/gameplay/pc"
	"github.com/YarikRevich/HideSeek-Client/internal/history"
	// "github.com/YarikRevich/HideSeek-Client/internal/player_mechanics/animation"
	"github.com/YarikRevich/HideSeek-Client/internal/render"
	imageloader "github.com/YarikRevich/HideSeek-Client/internal/resource_manager/loader/image_loader"
	// metadataloader "github.com/YarikRevich/HideSeek-Client/internal/resource_manager/loader/metadata_loader"
	"github.com/hajimehoshi/ebiten/v2"
)

func Draw() {


	render.SetToRender(func(screen *ebiten.Image) {
		img := imageloader.GetImage("assets/images/maps/default/background/Game")

		opts := &ebiten.DrawImageOptions{}

		imageW, imageH := img.Size()
		screenW, screenH := screen.Size()
		opts.GeoM.Scale(float64(screenW)/float64(imageW), float64(screenH)/float64(imageH))

		screen.DrawImage(img, opts)
	})

	p := pc.GetPC()
	// c := animation.WithAnimation(
	// 	imageloader.GetImage("/images/heroes/pumpkinhero"),
	// 	metadataloader.MetadataCollection["/images/heroes/pumpkinhero"],
	// 	&p.Equipment.Skin.Animation)
	render.SetToRender(func(i *ebiten.Image) {
		opts := &ebiten.DrawImageOptions{}

		if history.GetDirection() == direction.LEFT {
			opts.GeoM.Scale(-1, 1)
		}
		if history.GetDirection() == direction.RIGHT {
			opts.GeoM.Scale(1, 1)
		}

		opts.GeoM.Translate(p.X, p.Y)
		

	})

	// for _, otherC := range pc.PCs{

	// }

	// for _, otherPCs := range {
	// 	img := 	imageloader.Images[players.Equipment.Skin.ImageHash]
	// 	render.SetImageToRender(img, func(i *ebiten.Image) *ebiten.DrawImageOptions {
	// 		return &ebiten.DrawImageOptions{}
	// 	})
	// }
	// screen.DrawImage(, &ebiten.DrawImageOptions{})
	// g.winConf.DrawGameBackground()

	// // g.winConf.DrawGoldChest()

	// g.mapComponents.GetCollisions().GetDoorsCollisions().DrawDoors(g.winConf.DrawHorDoor, g.winConf.DrawVerDoor)

	// Animation.NewDefaultSwordAnimator(g.winConf, g.userConfig).Move()
	// Animation.NewIconAnimator(g.winConf, g.userConfig).Move()

	// for _, value := range g.winConf.GameProcess.OtherUsers {
	// 	Animation.NewDefaultSwordAnimator(g.winConf, value).Move()
	// 	Animation.NewIconAnimator(g.winConf, value).Move()
	// }

	// g.winConf.DrawDarkness(pixel.V((float64(g.userConfig.Pos.X)*2.5)-31, (float64(g.userConfig.Pos.Y)*2.5)-30))

	// g.winConf.DrawElementsPanel()

	// g.mapComponents.GetCam().UpdateCam()

	// var bias float64
	// for i := 0; i <= g.userConfig.GameInfo.Health; i++ {
	// 	g.winConf.DrawHPHeart(
	// 		pixel.V(-40+bias, 1200),
	// 	)
	// 	bias += 100
	// }

	// g.winConf.DrawWeaponIcon(g.userConfig.GameInfo.WeaponName)

	// if g.userConfig.GameInfo.Health < 1 {
	// 	g.mapComponents.GetCam().SetDefaultCam()
	// 	g.currState.MainStates.SetStartMenu()
	// }
}
