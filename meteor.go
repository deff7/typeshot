package main

import "github.com/faiface/pixel"

type meteor struct {
	pos   pixel.Vec
	word  string
	speed float64
}

func (m *meteor) update(dt float64) {
	m.pos = m.pos.Add(pixel.Vec{0, -m.speed * dt})
}
