package main

import (
	"github.com/faiface/pixel"
)

type beam struct {
	speed    float64
	angle    float64
	lifetime float64
	curTime  float64
	dead     bool

	target *meteor
}

func newBeam(playerPos pixel.Vec, target *meteor) *beam {
	b := &beam{
		speed: 500,
		angle: angleBetweenVecs(playerPos, target.pos),
	}
	b.lifetime = distantionBetweenVecs(playerPos, target.pos) / b.speed
	b.target = target
	return b
}

func (b *beam) update(dt float64) {
	b.curTime += dt
	if b.curTime > b.lifetime {
		b.dead = true
		b.target.destroy()
	}
}
