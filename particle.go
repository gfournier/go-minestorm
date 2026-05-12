package main

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type Particle struct {
	Pos     Vec2
	Vel     Vec2
	Life    int
	MaxLife int
	Clr     color.RGBA
}

func newParticle(pos Vec2, clr color.RGBA) Particle {
	angle := rand.Float64() * math.Pi * 2
	speed := rand.Float64()*3 + 0.5
	life := rand.Intn(30) + 20
	return Particle{
		Pos:     pos,
		Vel:     Vec2{math.Cos(angle) * speed, math.Sin(angle) * speed},
		Life:    life,
		MaxLife: life,
		Clr:     clr,
	}
}

func SpawnExplosion(pos Vec2, size MineSize, particles *[]Particle) {
	clr := mineColor(size)
	count := 14 - int(size)*2 // 14 large, 12 medium, 10 small
	for i := 0; i < count; i++ {
		*particles = append(*particles, newParticle(pos, clr))
	}
}

func SpawnShipExplosion(pos Vec2, particles *[]Particle) {
	for i := 0; i < 20; i++ {
		*particles = append(*particles, newParticle(pos, colorShip))
	}
}

func (p *Particle) Update() {
	p.Pos = wrapPos(p.Pos.Add(p.Vel), ScreenWidth, ScreenHeight)
	p.Vel = p.Vel.Scale(0.95)
	p.Life--
}

func (p *Particle) Draw(screen *ebiten.Image) {
	alpha := float64(p.Life) / float64(p.MaxLife)
	clr := color.RGBA{
		p.Clr.R, p.Clr.G, p.Clr.B,
		uint8(float64(p.Clr.A) * alpha),
	}
	x, y := float32(p.Pos.X), float32(p.Pos.Y)
	dx, dy := float32(p.Vel.X), float32(p.Vel.Y)
	strokeLine(screen, x-dx, y-dy, x+dx, y+dy, clr)
}
