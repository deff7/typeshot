package main

import (
	"image"
	_ "image/png"
	"log"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type game struct {
	sprites map[string]*pixel.Sprite
}

func newGame() *game {
	g := &game{}
	g.sprites = map[string]*pixel.Sprite{}
	g.loadSprites()
	return g
}

func (g *game) loadSprites() {
	for _, name := range []string{"bg", "player", "meteor"} {
		pic, err := loadPicture("data/" + name + ".png")
		if err != nil {
			log.Fatal(err)
		}
		g.sprites[name] = pixel.NewSprite(pic, pic.Bounds())
	}
}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func getWH(rect pixel.Rect) (float64, float64) {
	return rect.Max.X, rect.Max.Y
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Typeshot",
		Bounds: pixel.R(0, 0, 800, 600),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		log.Fatal(err)
	}

	g := newGame()

	win.Clear(colornames.Black)

	m := pixel.IM
	winW, winH := getWH(win.Bounds())

	spr := g.sprites["bg"]
	sprW, sprH := getWH(spr.Frame())

	for y := 0.0; y < winH; y += sprH {
		for x := 0.0; x < winW; x += sprW {
			spr.Draw(win, m.Moved(pixel.Vec{x, y}))
		}
	}

	for !win.Closed() {
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
