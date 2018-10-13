package main

import (
	"bufio"
	"math"
	"math/rand"
	"os"
	"strings"

	"github.com/faiface/pixel"
)

var dictionary = []string{
	"foo",
	"bar",
	"test",
}

func (g *game) loadDictionary(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	s := bufio.NewScanner(f)
	for s.Scan() {
		g.dictionary = append(g.dictionary, strings.TrimSpace(s.Text()))
	}
	g.dictIndicies = g.shuffleDict()

	return s.Err()
}

func (g *game) shuffleDict() []int {
	return rand.Perm(len(g.dictionary))
}

func (g *game) nextWord() int {
	wordIdx := g.dictIndicies[g.dictIdx]

	g.dictIdx++
	if g.dictIdx >= len(g.dictIndicies) {
		g.dictIdx = 0
		g.dictIndicies = g.shuffleDict()
	}
	return wordIdx
}

func (g *game) generateWord() string {
	return g.dictionary[g.nextWord()]
}

func (g *game) generatePosition(w, h float64) pixel.Vec {
	x := w/2 + rand.Float64()*(g.winW-w)
	return pixel.Vec{x, g.winH - h/2}
}

func (g *game) spawnMeteor(gameSpeed float64) *meteor {
	w, h := getWH(g.sprites["meteor"].Frame())
	m := &meteor{
		word:  g.generateWord(),
		pos:   g.generatePosition(w, h),
		angle: rand.Float64() * 2 * math.Pi,
		speed: 30 * gameSpeed,
	}
	return m
}
