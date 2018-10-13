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

func findMeteor(meteors map[*meteor]bool, text string) *meteor {
	var nearest *meteor = nil
	for m, _ := range meteors {
		if m.word == text {
			if nearest == nil || m.pos.Y < nearest.pos.Y {
				nearest = m
			}
		}
	}
	return nearest
}

func processInput(input string) (string, bool) {
	ss := strings.Split(input, " ")
	return ss[0], len(ss) > 1
}

func angleBetweenVecs(v1, v2 pixel.Vec) float64 {
	return math.Atan(-(v1.X - v2.X) / (v1.Y - v2.Y))
}

func distantionBetweenVecs(v1, v2 pixel.Vec) float64 {
	return math.Sqrt(math.Pow(v2.X-v1.X, 2) + math.Pow(v2.Y-v1.Y, 2))
}

func animTowards(src, dest, speed float64) float64 {
	if math.Abs(src-dest) <= speed {
		return dest
	}

	if src < dest {
		src += speed
	} else {
		src -= speed
	}
	return src
}

func run() {
	g := newGame()
	gameSpeed := 1.0

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
		spawner = time.Tick(2 * time.Second)
	)

	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	g.initText(atlas)
	score := 0
	scoreText := text.New(pixel.V(60, 50), atlas)

	meteors := map[*meteor]bool{}
	beams := map[*beam]bool{}

	angle := 0.0
	angleDest := 0.0

	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		win.Clear(colornames.Black)

		g.drawBackground(win)

		met := findMeteor(meteors, g.current)
		if met != nil {
			angleDest = angleBetweenVecs(met.pos, g.playerPos)
		} else {
			angleDest = 0
		}
		angle = animTowards(angle, angleDest, 0.09)

		for m, _ := range meteors {
			if m.dead {
				score++
				delete(meteors, m)
				continue
			}
			g.drawMeteor(win, m)
			m.update(dt)
		}

		for b, _ := range beams {
			if b.dead {
				delete(beams, b)
				continue
			}
			if b.curTime < b.lifetime {
				g.drawBeam(win, b)
			}
			b.update(dt)
		}
		g.drawPlayer(win, angle)

		g.drawCurrentInput(win, atlas)

		scoreText.Clear()
		fmt.Fprintf(scoreText, "Score: %d", score)
		scoreText.Draw(win, pixel.IM.Scaled(scoreText.Orig, 2))

		var shot bool
		g.current, shot = processInput(g.current + win.Typed())
		if shot || win.JustPressed(pixelgl.KeyEnter) {
			if met != nil {
				b := newBeam(g.playerPos, met)
				beams[b] = true
			}
			g.current = ""
		}
		if win.JustPressed(pixelgl.KeyBackspace) {
			if len(g.current) > 0 {
				g.current = g.current[:len(g.current)-1]
			}
		}

		win.Update()

		frames++
		select {
		case <-spawner:
			m := g.spawnMeteor(gameSpeed)
			m.initText(atlas)
			meteors[m] = true
		case <-tick:
			gameSpeed += 0.05
			win.SetTitle(fmt.Sprintf("%s | FPS: %d | spd: %f", cfg.Title, frames, gameSpeed))
			frames = 0
		default:
		}
	}
}

func main() {
	rand.Seed(int64(time.Now().Second()))
	pixelgl.Run(run)
}
