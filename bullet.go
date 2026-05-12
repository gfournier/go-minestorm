package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	BulletSpeed    = 9.0
	BulletLifetime = 55
	BulletRadius   = 4.0
)

type Bullet struct {
	Pos      Vec2
	Vel      Vec2
	Lifetime int
}

func NewBullet(ship Ship) Bullet {
	return Bullet{
		Pos: ship.Pos,
		Vel: Vec2{
			math.Cos(ship.Angle)*BulletSpeed + ship.Vel.X,
			math.Sin(ship.Angle)*BulletSpeed + ship.Vel.Y,
		},
		Lifetime: BulletLifetime,
	}
}

func (b *Bullet) Update() {
	b.Pos = wrapPos(b.Pos.Add(b.Vel), ScreenWidth, ScreenHeight)
	b.Lifetime--
}

func (b *Bullet) Draw(screen *ebiten.Image) {
	x, y := float32(b.Pos.X), float32(b.Pos.Y)
	dx, dy := float32(b.Vel.X*0.6), float32(b.Vel.Y*0.6)
	strokeLine(screen, x-dx, y-dy, x+dx, y+dy, colorBullet)
}
