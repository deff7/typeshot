package main

import (
	"fmt"
	"image"
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"os"
	"strings"
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

func findMeteorIndex(meteors []*meteor, text string) int {
	for i, m := range meteors {
		if m.word == text {
			return i
		}
	}
	return -1
}

func processInput(input string) (string, bool) {
	ss := strings.Split(input, " ")
	return ss[0], len(ss) > 1
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
		frames  = 0
		tick    = time.Tick(time.Second)
		spawner = time.Tick(3 * time.Second)
	)

	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	g.initText(atlas)
	scoreText := text.New(pixel.V(60, 50), atlas)
	fmt.Fprintf(scoreText, "Score: %d", 0)

	meteors := []*meteor{}

	angle := 0.0

	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		win.Clear(colornames.Black)

		g.drawBackground(win)

		mIdx := findMeteorIndex(meteors, g.current)
		if mIdx != -1 {
			m := meteors[mIdx]
			angle = math.Atan(-(m.pos.X - g.playerPos.X) / (m.pos.Y - g.playerPos.Y))
		}
		g.drawPlayer(win, angle)

		for _, m := range meteors {
			g.drawMeteor(win, m)
			m.update(dt)
		}

		g.drawCurrentInput(win, atlas)
		scoreText.Draw(win, pixel.IM.Scaled(scoreText.Orig, 2))

		var shot bool
		g.current, shot = processInput(g.current + win.Typed())
		if shot {
			g.current = ""
		}

		win.Update()

		frames++
		select {
		case <-spawner:
			m := g.spawnMeteor()
			m.initText(atlas)
			meteors = append(meteors, m)
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
