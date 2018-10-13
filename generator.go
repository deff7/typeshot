package main

import (
	"math/rand"

	"github.com/faiface/pixel"
)

var dictionary = []string{
	"foo",
	"bar",
	"test",
}

func generateWord() string {
	i := rand.Int() % len(dictionary)
	return dictionary[i]
}

func (g *game) generatePosition(w, h float64) pixel.Vec {
	x := w/2 + rand.Float64()*(g.winW-w)
	return pixel.Vec{x, g.winH - h/2}
}

func (g *game) spawnMeteor() *meteor {
	w, h := getWH(g.sprites["meteor"].Frame())
	m := &meteor{
		word:  generateWord(),
		pos:   g.generatePosition(w, h),
		speed: 30,
	}
	return m
}
