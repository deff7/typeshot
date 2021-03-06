package main

import (
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
)

type meteor struct {
	pos   pixel.Vec
	word  string
	speed float64
	angle float64
	dead  bool

	text *text.Text
}

func (m *meteor) update(dt float64) {
	m.pos = m.pos.Add(pixel.Vec{0, -m.speed * dt})
	m.angle += 0.01
	if m.angle > 2*math.Pi {
		m.angle = 0
	}
}

func (m *meteor) initText(atlas *text.Atlas) {
	m.text = text.New(pixel.V(0, 0), atlas)
	m.text.Dot.X -= m.text.BoundsOf(m.word).W() / 2
	m.text.Write([]byte(m.word))
}

func (m *meteor) destroy() {
	m.dead = true
}
