package sources

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"path/filepath"
	"sort"
	"strings"

	"github.com/YarikRevich/hide-seek-client/assets"
	"github.com/YarikRevich/hide-seek-client/internal/core/screen"
	"github.com/YarikRevich/hide-seek-client/internal/core/types"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lafriks/go-tiled"
	"github.com/sirupsen/logrus"
)

type Tile struct {
	Image      *ebiten.Image
	Position   image.Point
	Layer      string
	LayerNum   int
	TileNum    int
	Properties struct {
		Collision, Spawn bool
	}
}

type AnimationFrame struct {
	Duration uint32
	TileID   image.Point
}

type Animation struct {
	Frames                     []*AnimationFrame
	CurrentFrame, DelayTrigger int
}

//Starts animation if hasn't been started
//or continues if not
func (a *Animation) Proceed() {
	a.DelayTrigger++
	a.DelayTrigger %= int(a.Frames[a.CurrentFrame].Duration)
	if a.DelayTrigger == 0 {
		a.CurrentFrame++
		a.CurrentFrame %= len(a.Frames)
	}
}

//Stops animation and returns frame and frame
//change trigger to start position
func (a *Animation) Reset() {
	a.DelayTrigger = 0
	a.CurrentFrame = 0
}

type Graph map[*Tile][]*Tile

func (g *Graph) AddNode(key, value *Tile) {
	_, ok := (*g)[key]
	if !ok {
		(*g)[key] = append((*g)[key], value)
	}

	_, ok = (*g)[value]
	if !ok {
		(*g)[value] = append((*g)[value], key)
	}
}

type OrthographicTile struct {
	Tile            *Tile
	Rotation, Pitch float64
}

type OrthographicTileFace int

const (
	Floor OrthographicTileFace = iota
	South
	North
	East
	West
	Top
)

type OrthographicTilebatch struct {
	IsWall bool
	Tiles  map[OrthographicTileFace]*OrthographicTile
	// Floor, Top, North, South, West, East *OrthographicTile
}

//Extension for tilemap
type OrthographicTilemap map[image.Point]*OrthographicTilebatch

type Quad struct {
	Tile
	Points [4]types.Vec3
}

type CubeOpts struct {
	sm *screen.ScreenManager

	Scale, Position types.Vec2
	Angle, Pitch    float64
	CameraPosition  types.Vec3
}

func CreateCube(opts CubeOpts) [8]types.Vec3 {

	var unitCube, rotCube, worldCube, projCube [8]types.Vec3

	fmt.Println(opts.Scale)

	unitCube[1] = types.Vec3{X: opts.Scale.X}
	unitCube[2] = types.Vec3{X: opts.Scale.X, Y: -opts.Scale.Y}
	unitCube[3] = types.Vec3{Y: -opts.Scale.Y}
	unitCube[4] = types.Vec3{Z: opts.Scale.Y}
	unitCube[5] = types.Vec3{X: opts.Scale.X, Z: opts.Scale.Y}
	unitCube[6] = types.Vec3{X: opts.Scale.X, Y: -opts.Scale.Y, Z: opts.Scale.Y}
	unitCube[7] = types.Vec3{Y: -opts.Scale.Y, Z: opts.Scale.Y}

	// unitCube[1] = types.Vec3{X: 0}
	// unitCube[2] = types.Vec3{X: 10, Y: 2}
	// unitCube[3] = types.Vec3{Y: 43}
	// unitCube[4] = types.Vec3{Z: 20}
	// unitCube[5] = types.Vec3{X: 20, Z: 0}
	// unitCube[6] = types.Vec3{X: 40, Y: -5, Z: 23}
	// unitCube[7] = types.Vec3{Y: -10, Z: 3}

	for i := 0; i < 8; i++ {
		unitCube[i].X += (opts.Position.X*opts.Scale.X - opts.CameraPosition.X)
		unitCube[i].Y += -opts.CameraPosition.Y
		unitCube[i].Z += (opts.Position.Y*opts.Scale.Y - opts.CameraPosition.Z)
	}

	s := math.Sin(opts.Angle)
	c := math.Cos(opts.Angle)
	for i := 0; i < 8; i++ {
		rotCube[i].X = unitCube[i].X*c + unitCube[i].Z*s
		rotCube[i].Y = unitCube[i].Y
		rotCube[i].Z = unitCube[i].X*-s + unitCube[i].Z*c
	}

	s = math.Sin(opts.Pitch)
	c = math.Cos(opts.Pitch)
	for i := 0; i < 8; i++ {
		worldCube[i].X = rotCube[i].X
		worldCube[i].Y = rotCube[i].Y*c - rotCube[i].Z*s
		worldCube[i].Z = rotCube[i].Y*s + rotCube[i].Z*c
	}

	screenSize := opts.sm.GetSize()

	for i := 0; i < 8; i++ {
		projCube[i].X = worldCube[i].X + screenSize.X*0.5
		projCube[i].Y = worldCube[i].Y + screenSize.Y*0.5
		projCube[i].Z = worldCube[i].Z
	}

	return projCube
}

// type FaceQuadOpts struct {
// 	CubeOpts
// }

// func GetFaceQuad(opts FaceQuadOpts) Quad {
// 	CreateCube(CubeOpts{
// 		sm:             opts.sm,
// 		CameraPosition: opts.CameraPosition,
// 		Scale:          opts.Scale,
// 	})
// 	return Quad{}
// }

type Tilemap struct {
	Name string

	//Collection of animations which
	//can be applied to the tilemap
	//SAVES STATE
	Animations map[int]*Animation
	Graph
	OrthographicTilemap

	Tiles      map[image.Point]*Tile
	Properties struct {
		//PHYSICS

		//Gravity acceleration
		G int

		//Acceleration
		A int

		//WORLD MAP properties

		//Contains IDs of Spawn Tiles
		Spawns []int64
	}

	MapSize, TileSize, TileCount types.Vec2
}

func (tm *Tilemap) ToAPIMessage() {

}

func (tm *Tilemap) load(path string) error {
	gameMap, err := tiled.LoadFile(fmt.Sprintf("%s.%s", path, "tmx"), tiled.WithFileSystem(assets.Assets))
	if err != nil {
		return err
	}

	baseDir := filepath.Dir(path)

	tm.MapSize = types.Vec2{
		X: float64((gameMap.Width) * gameMap.TileWidth),
		Y: float64((gameMap.Height) * gameMap.TileHeight)}
	tm.TileSize = types.Vec2{
		X: float64(gameMap.TileWidth),
		Y: float64(gameMap.TileHeight)}
	tm.TileCount = types.Vec2{
		X: float64(gameMap.Width),
		Y: float64(gameMap.Height),
	}

	tempTileCollection := make(map[image.Point]*Tile)
	tempTileImageCache := make(map[string]*ebiten.Image)

	for n, l := range gameMap.Layers {
		y := 0
		for i, t := range l.Tiles {

			if !t.IsNil() {
				tile := new(Tile)

				if _, ok := tempTileImageCache[t.Tileset.Image.Source]; !ok {
					file, err := assets.Assets.Open(filepath.Join(baseDir, t.Tileset.Image.Source))
					if err != nil {
						logrus.Fatalln(err)
					}
					pngFile, _, err := image.Decode(file)
					if err != nil {
						logrus.Fatalln(err)
					}
					if err := file.Close(); err != nil {
						logrus.Fatalln(err)
					}
					ebitenImage := ebiten.NewImageFromImage(pngFile)
					tempTileImageCache[t.Tileset.Image.Source] = ebitenImage
				}

				ebitenImage := tempTileImageCache[t.Tileset.Image.Source]
				subImage := ebitenImage.SubImage(t.GetTileRect())

				tile.Image = ebiten.NewImageFromImage(subImage)

				tile.Layer = l.Name
				tile.LayerNum = n

				x := (i % gameMap.Width) * gameMap.TileWidth

				for _, w := range t.Tileset.Tiles {
					if w.ID == t.ID {
						animation := new(Animation)
						for _, a := range w.Animation {
							animation.Frames = append(animation.Frames, &AnimationFrame{
								Duration: a.Duration,
								TileID:   image.Point{X: x, Y: y},
							})
						}
						tm.Animations[len(tm.Animations)+1] = animation
						tile.Properties.Collision = w.Properties.GetBool("collision")
					}
				}

				if top, ok := tempTileCollection[image.Point{X: x, Y: y - gameMap.TileHeight}]; ok {
					tm.Graph.AddNode(
						tile, top)
				}

				if bottom, ok := tempTileCollection[image.Point{X: x, Y: y + gameMap.TileHeight}]; ok {
					tm.Graph.AddNode(
						tile, bottom)
				}

				if right, ok := tempTileCollection[image.Point{X: x + gameMap.TileWidth, Y: y}]; ok {
					tm.Graph.AddNode(
						tile, right)
				}

				if left, ok := tempTileCollection[image.Point{X: x + gameMap.TileWidth, Y: y}]; ok {
					tm.Graph.AddNode(
						tile, left)
				}

				tempTileCollection[image.Point{X: x, Y: y}] = tile

				tm.OrthographicTilemap[image.Point{X: x, Y: y}] = &OrthographicTilebatch{
					Tiles: map[OrthographicTileFace]*OrthographicTile{
						Floor: {Tile: tile},
						Top:   {Tile: tile},
						South: {Tile: tile},
						North: {Tile: tile},
						West:  {Tile: tile},
						East:  {Tile: tile},
					},
					IsWall: false,
				}
				tile.TileNum = len(tm.Tiles) + 1
				tile.Position = image.Point{X: x, Y: y}
				tm.Tiles[tile.Position] = tile

				if (i%gameMap.Width)*gameMap.TileWidth == ((gameMap.Width * gameMap.TileWidth) - gameMap.TileWidth) {
					y += gameMap.TileHeight
				}
			}
		}
	}

	tm.Name = strings.Split(path, ".")[0]
	tileMapCollection[tm.Name] = *tm

	return nil
}

type RenderTilemapOpts struct {
	StickedTo                            *Tilemap
	StickedToPosition                    types.Vec2
	SurfacePosition, Scale               types.Vec2
	CameraAngle, CameraPitch             float64
	CameraPosition                       types.Vec3
	AutoScaleForbidden, CenterizedOffset bool

	OrthigraphicProjection bool
}

func (t *Tilemap) Render(sm *screen.ScreenManager, opts RenderTilemapOpts) {
	screenSize := sm.GetSize()
	screenScale := sm.GetScale()

	var orthographicPostRender []*Tile
	for k, v := range t.Tiles {
		if (float64(k.X)+opts.SurfacePosition.X-t.TileSize.X < screenSize.X && float64(k.Y)+opts.SurfacePosition.Y-t.TileSize.Y < screenSize.Y) &&
			(float64(k.X)+opts.SurfacePosition.X+t.TileSize.X > 0 && float64(k.Y)+opts.SurfacePosition.Y+t.TileSize.Y > 0) {
			drawOpts := &ebiten.DrawImageOptions{}

			if opts.OrthigraphicProjection {
				orthographicPostRender = append(orthographicPostRender, v)
				// projectionCube := CreateCube(CubeOpts{
				// 	sm:             sm,
				// 	Scale:          opts.Scale,
				// 	Position:       opts.SurfacePosition,
				// 	Angle:          opts.CameraAngle,
				// 	Pitch:          opts.CameraPitch,
				// 	CameraPosition: opts.CameraPosition,
				// })
				// for _, c := range [][5]int{
				// 	{4, 0, 1, 5, int(Floor)},
				// 	{3, 0, 1, 2, int(South)},
				// 	{6, 5, 4, 7, int(North)},
				// 	{7, 4, 0, 3, int(East)},
				// 	{2, 1, 5, 6, int(West)},
				// 	{7, 3, 2, 6, int(Top)},
				// } {
				// v1, v2, v3, v4 := projectionCube[c[0]], projectionCube[c[1]], projectionCube[c[2]], projectionCube[c[3]]
				// var path vector.Path
				// path.MoveTo(float32(k.X*int(opts.Scale.X)), float32(k.Y*int(opts.Scale.Y)))
				// path.LineTo(float32(v1.X), float32(v2.Y))
				// path.LineTo(float32(v2.X), float32(v3.Y))
				// path.LineTo(float32(v3.X), float32(v4.Y))
				// path.LineTo(float32(v4.X), float32(v1.Y))
				// path.Fill(sm.Image, &vector.FillOptions{Color: color.Opaque})

				// w, h := v.Image.Size()
				// imageWidth, imageHeight := float64(w), float64(h)
				// imageWidth += float64(v.TileNum / 2)
				// // imageHeight -= float64(v.TileNum)
				// // imageWidth -= (opts.CameraAngle * 1.2)
				// // imageHeight *= math.Cos(opts.CameraAngle)
				// fmt.Println(v.TileNum % int(t.TileCount.X))
				// opts := &ebiten.DrawImageOptions{}
				// for i := 0; i < int(imageHeight); i++ {
				// 	opts.GeoM.Reset()

				// 	// opts.GeoM.Translate(-screenSize.X/2, -screenSize.Y/2)
				// 	opts.GeoM.Translate(-imageWidth/2-float64(v.TileNum), -imageHeight/2)
				// 	lineW := (int(imageWidth) + i*3/4)
				// 	x := -float64(lineW) / imageWidth / 2
				// 	opts.GeoM.Scale(float64(lineW)/imageWidth, 1)
				// 	opts.GeoM.Translate(x, float64(i))
				// 	opts.GeoM.Translate(screenSize.X/2, screenSize.Y/2)
				// 	opts.GeoM.Translate(float64(k.X), float64(k.Y))
				// 	// sm.Image.DrawImage(t.OrthographicTilemap[k].Tiles[OrthographicTileFace(c[4])].Tile.Image, opts)
				// 	sm.Image.DrawImage(v.Image.SubImage(image.Rect(0, i, int(imageWidth), i+1)).(*ebiten.Image), opts)
				// }

				// opts.GeoM.Rotate(float64(10%360) * 2 * math.Pi / 360)

				// opts.GeoM.Translate(v3.X, v3.X)
				// opts.GeoM.Translate(-float64(screenSize.X)/2, -float64(screenSize.Y)/2)

				// opts.GeoM.Translate(float64(screenSize.X)/2, float64(screenSize.Y)/2)

				// sm.Image.DrawImage(t.OrthographicTilemap[k].Tiles[OrthographicTileFace(c[4])].Tile.Image, opts)
				// fmt.Println(c)
				// ebitenutil.DrawRect(sm.Image, float64(v1.X), float64(v1.Y), 200, 200, color.Opaque)
				// ebitenutil.DrawRect(sm.Image, float64(v2.X)+10, float64(v2.Y), 200, 200, color.RGBA{200, 220, 110, 255})
				// ebitenutil.DrawRect(sm.Image, float64(v3.X), float64(v3.Y), 200, 200, color.RGBA{220, 210, 140, 255})
				// ebitenutil.DrawRect(sm.Image, float64(v4.X)+10, float64(v4.Y), 200, 200, color.RGBA{110, 210, 140, 255})
				// }
			} else {
				if !opts.AutoScaleForbidden {
					drawOpts.GeoM.Scale(1/screenScale.X, 1/screenScale.Y)
				}

				if opts.StickedTo != nil {
					drawOpts.GeoM.Translate(opts.StickedToPosition.X, opts.StickedToPosition.Y)
				}

				drawOpts.GeoM.Translate(float64(k.X), float64(k.Y))

				if opts.CenterizedOffset {
					drawOpts.GeoM.Translate(-t.MapSize.X/2, -t.MapSize.Y/2)
				}

				if opts.Scale.X != 0 && opts.Scale.Y != 0 {
					drawOpts.GeoM.Scale(opts.Scale.X, opts.Scale.Y)
				}

				drawOpts.GeoM.Translate(opts.SurfacePosition.X, opts.SurfacePosition.Y)
				sm.Image.DrawImage(v.Image, drawOpts)
			}
		}
	}

	if opts.OrthigraphicProjection {
		sort.Slice(orthographicPostRender, func(i, j int) bool {
			return orthographicPostRender[i].TileNum < orthographicPostRender[j].TileNum
		})

		orthographicSurface := ebiten.NewImage(int(t.MapSize.X), int(t.MapSize.Y))
		w, h := orthographicSurface.Size()
		fmt.Println(w, t.MapSize.X)
		for _, v := range orthographicPostRender {
			for y := 0; y < h; y++ {
				for x := 0; x < w; x++ {
					// fmt.Println(v.Position, x, y, v.Image.At(x, y))
					orthographicSurface.Set(v.Position.X+int(t.TileSize.X)+x, v.Position.Y+int(t.TileSize.Y)+y, v.Image.At(x, y))
				}
			}
		}

		// op := &ebiten.DrawImageOptions{}
		// op.GeoM.Translate(100, 100)
		// // w, h := orthographicSurface.Size()
		// // for ?

		// sm.Image.DrawImage(orthographicSurface, op)

		// fmt.Println(w)
		op := &ebiten.DrawImageOptions{}
		for i := 0; i < h; i++ {
			op.GeoM.Reset()

			op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
			lineW := (w + i*3/4)
			x := -float64(lineW) / float64(w) / 2

			op.GeoM.Scale(float64(lineW)/float64(w)*math.Sin(opts.CameraPitch), 1)
			// fmt.Println(float64(lineW) / float64(w))
			op.GeoM.Translate(x, float64(i))
			op.GeoM.Translate(screenSize.X/2, screenSize.Y/2)
			// sm.Image.DrawImage(t.OrthographicTilemap[k].Tiles[OrthographicTileFace(c[4])].Tile.Image, opts)
			sm.Image.DrawImage(orthographicSurface.SubImage(image.Rect(0, i, w, i+1)).(*ebiten.Image), op)
		}
	}
}

func NewTilemap() Tilemap {
	return Tilemap{
		Tiles:               make(map[image.Point]*Tile),
		Animations:          make(map[int]*Animation),
		OrthographicTilemap: make(OrthographicTilemap),
		Graph:               make(Graph),
	}
}
