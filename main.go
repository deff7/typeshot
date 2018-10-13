package main

import (
	"fmt"
	"image"
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
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

	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	g.initText(atlas)
	scoreText := text.New(pixel.V(60, 50), atlas)
	fmt.Fprintf(scoreText, "Score: %d", 0)

	m := g.spawnMeteor()
	m.initText(atlas)

	angle := 0.0

	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		win.Clear(colornames.Black)

		g.drawBackground(win)
		angle = math.Atan(-(m.pos.X - g.playerPos.X) / (m.pos.Y - g.playerPos.Y))
		g.drawPlayer(win, angle)
		g.drawMeteor(win, m)

		g.drawCurrentInput(win, atlas)
		scoreText.Draw(win, pixel.IM.Scaled(scoreText.Orig, 2))

		g.current += win.Typed()

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
