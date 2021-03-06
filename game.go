package main

import (
	"fmt"
	"log"
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
)

type game struct {
	sprites      map[string]*pixel.Sprite
	winW, winH   float64
	current      string
	playerPos    pixel.Vec
	dictionary   []string
	dictIndicies []int
	dictIdx      int

	bgParallax float64

	currentText *text.Text
}

func newGame() *game {
	g := &game{winW: 800, winH: 600}
	g.sprites = map[string]*pixel.Sprite{}
	g.loadSprites()
	err := g.loadDictionary("data/words.txt")
	if err != nil {
		log.Fatal(err)
	}
	g.playerPos = pixel.Vec{g.winW / 2, g.sprites["player"].Frame().Max.Y}
	return g
}

func (g *game) initText(atlas *text.Atlas) {
	g.currentText = text.New(pixel.ZV, atlas)
}

func (g *game) loadSprites() {
	for _, name := range []string{"bg", "player", "meteor", "beam"} {
		pic, err := loadPicture("data/" + name + ".png")
		if err != nil {
			log.Fatal(err)
		}
		g.sprites[name] = pixel.NewSprite(pic, pic.Bounds())
	}
}

func (g *game) drawBackground(target pixel.Target) {

	s := g.sprites["bg"]
	w, h := getWH(s.Frame())
	m := pixel.IM.Moved(pixel.V(0, -g.bgParallax*h))

	for y := -h; y < g.winH+h; y += h {
		for x := 0.0; x < g.winW; x += w {
			s.Draw(target, m.Moved(pixel.Vec{x, y}))
		}
	}

	g.bgParallax += 0.005
	if g.bgParallax > 1 {
		g.bgParallax = 0
	}
}

func (g *game) drawPlayer(target pixel.Target, angle float64) {
	s := g.sprites["player"]
	mat := pixel.IM.Rotated(pixel.ZV, angle)
	mat = mat.Moved(g.playerPos)
	s.Draw(target, mat)
}

func (g *game) drawMeteor(target pixel.Target, m *meteor) {
	s := g.sprites["meteor"]
	_, h := getWH(s.Frame())
	mat := pixel.IM.Rotated(pixel.ZV, m.angle)
	mat = mat.Moved(m.pos)
	s.Draw(target, mat)
	mat = pixel.IM.Moved(m.pos)
	m.text.Draw(target, mat.Moved(pixel.V(0, h/2+5)))
}

func (g *game) drawCurrentInput(target pixel.Target, atlas *text.Atlas) {
	g.currentText.Dot.X -= g.currentText.BoundsOf(g.current).W()/2 + 2
	fmt.Fprint(g.currentText, g.current)
	mat := pixel.IM.Scaled(g.currentText.Orig, 2)
	mat = mat.Moved(pixel.V(g.winW/2, 5))
	g.currentText.Draw(target, mat)
	g.currentText.Clear()
}

func (g *game) drawBeam(target pixel.Target, b *beam) {
	s := g.sprites["beam"]
	x := math.Cos(b.angle+math.Pi/2) * b.speed * b.curTime
	y := math.Sin(b.angle+math.Pi/2) * b.speed * b.curTime
	mat := pixel.IM.Rotated(pixel.ZV, b.angle)
	mat = mat.Moved(g.playerPos)
	mat = mat.Moved(pixel.V(x, y))
	s.Draw(target, mat)
}
