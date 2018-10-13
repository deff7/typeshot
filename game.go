package main

import (
	"fmt"
	"log"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
)

type game struct {
	sprites    map[string]*pixel.Sprite
	winW, winH float64
	current    string

	currentText *text.Text
}

func newGame() *game {
	g := &game{winW: 800, winH: 600}
	g.sprites = map[string]*pixel.Sprite{}
	g.loadSprites()
	return g
}

func (g *game) initText(atlas *text.Atlas) {
	g.currentText = text.New(pixel.ZV, atlas)
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

func (g *game) drawBackground(target pixel.Target) {
	m := pixel.IM

	s := g.sprites["bg"]
	w, h := getWH(s.Frame())

	for y := 0.0; y < g.winH; y += h {
		for x := 0.0; x < g.winW; x += w {
			s.Draw(target, m.Moved(pixel.Vec{x, y}))
		}
	}
}

func (g *game) drawPlayer(target pixel.Target) {
	s := g.sprites["player"]
	m := pixel.IM.Moved(pixel.Vec{g.winW / 2, s.Frame().Max.Y})
	s.Draw(target, m)
}

func (g *game) drawMeteor(target pixel.Target, m *meteor) {
	s := g.sprites["meteor"]
	_, h := getWH(s.Frame())
	mat := pixel.IM.Moved(m.pos)
	s.Draw(target, mat)
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
