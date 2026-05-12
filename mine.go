package main

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type MineSize int

const (
	SizeLarge  MineSize = 0
	SizeMedium MineSize = 1
	SizeSmall  MineSize = 2
)

var (
	mineRadii     = [3]float64{22, 13, 7}
	minePoints    = [3]int{10, 20, 40}
	mineMaxSpeeds = [3]float64{2.0, 3.2, 4.8}
)

type Mine struct {
	Pos        Vec2
	Vel        Vec2
	Size       MineSize
	Angle      float64
	AngularVel float64
}

func NewMine(size MineSize, pos, shipPos Vec2, speedMult float64) Mine {
	dir := wrappedDelta(pos, shipPos, ScreenWidth, ScreenHeight).Norm()
	// Spread direction randomly within ±50° so not all mines beeline straight at you
	dir = rotateVec(dir, (rand.Float64()-0.5)*math.Pi*0.55)
	speed := mineMaxSpeeds[size] * speedMult
	return Mine{
		Pos:        pos,
		Vel:        dir.Scale(speed),
		Size:       size,
		Angle:      rand.Float64() * math.Pi * 2,
		AngularVel: (rand.Float64() - 0.5) * 0.07,
	}
}

func RandomEdgePos() Vec2 {
	switch rand.Intn(4) {
	case 0:
		return Vec2{rand.Float64() * ScreenWidth, -5}
	case 1:
		return Vec2{rand.Float64() * ScreenWidth, ScreenHeight + 5}
	case 2:
		return Vec2{-5, rand.Float64() * ScreenHeight}
	default:
		return Vec2{ScreenWidth + 5, rand.Float64() * ScreenHeight}
	}
}

func (m *Mine) Update(shipPos Vec2, trackStr float64) {
	// Softly steer toward the ship each frame
	dir := wrappedDelta(m.Pos, shipPos, ScreenWidth, ScreenHeight).Norm()
	m.Vel.X += dir.X * trackStr
	m.Vel.Y += dir.Y * trackStr
	maxSpd := mineMaxSpeeds[m.Size]
	if spd := m.Vel.Mag(); spd > maxSpd {
		m.Vel = m.Vel.Scale(maxSpd / spd)
	}
	m.Pos = wrapPos(m.Pos.Add(m.Vel), ScreenWidth, ScreenHeight)
	m.Angle += m.AngularVel
}

func (m *Mine) Radius() float64 { return mineRadii[m.Size] }

func (m *Mine) Draw(screen *ebiten.Image) {
	pts := transformPoints(getMineShape(m.Size), m.Pos, m.Angle, mineRadii[m.Size])
	drawPoly(screen, pts, mineColor(m.Size))
}
