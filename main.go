package main

import (
	"fmt"
	"image"
	_ "image/png"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

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
	g := newGame()

	cfg := pixelgl.WindowConfig{
		Title:  "Typeshot",
		Bounds: pixel.R(0, 0, g.winW, g.winH),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		log.Fatal(err)
	}

	var (
		frames = 0
		tick   = time.Tick(time.Second)
	)

	m := g.spawnMeteor()
	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		win.Clear(colornames.Black)

		g.drawBackground(win)
		g.drawPlayer(win)
		g.drawMeteor(win, m.pos)
		m.update(dt)

		win.Update()
		frames++
		select {
		case <-tick:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}
	}
}

func main() {
	rand.Seed(int64(time.Now().Second()))
	pixelgl.Run(run)
}
