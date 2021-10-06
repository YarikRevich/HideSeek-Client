package settingsmenu

import (
	"image/color"

	buffercollection "github.com/YarikRevich/HideSeek-Client/internal/hid/keyboard/buffers/collection"
	"github.com/YarikRevich/HideSeek-Client/internal/interface/fonts"

	"github.com/YarikRevich/HideSeek-Client/internal/interface/positioning/button"
	"github.com/YarikRevich/HideSeek-Client/internal/render"
	imagecollection "github.com/YarikRevich/HideSeek-Client/internal/resource_manager/image_loader/collection"
	metadatacollection "github.com/YarikRevich/HideSeek-Client/internal/resource_manager/metadata_loader/collection"
	"github.com/hajimehoshi/ebiten/v2"
	// "github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

func Draw() {
	render.SetToRender(func(screen *ebiten.Image) {
		img := imagecollection.GetImage("assets/images/menues/background/background")

		opts := &ebiten.DrawImageOptions{}
		imageW, imageH := img.Size()
		screenW, screenH := screen.Size()
		opts.GeoM.Scale(float64(screenW)/float64(imageW), float64(screenH)/float64(imageH))

		screen.DrawImage(img, opts)
	})

	render.SetToRender(func(screen *ebiten.Image) {
		img := ebiten.NewImageFromImage(imagecollection.GetImage("assets/images/menues/inputs/input"))
		m := metadatacollection.GetMetadata("assets/images/menues/inputs/input")

		opts := &ebiten.DrawImageOptions{}

		opts.GeoM.Translate(m.Margins.LeftMargin, m.Margins.TopMargin)
		opts.GeoM.Scale(m.Scale.CoefficiantX, m.Scale.CoefficiantY)

		f := fonts.GetFont(*m)
		t := buffercollection.SettingsMenuNameBuffer.Read()
		tx, ty := button.ChooseButtonTextPosition(f, t, *m)

		text.Draw(img, t, f, tx, ty, color.White)

		screen.DrawImage(img, opts)
	})
}
