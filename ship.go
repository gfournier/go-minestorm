package main

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ShipMaxSpeed  = 7.0
	ShipThrust    = 0.20
	ShipRotSpeed  = 0.07
	ShipFriction  = 0.985
	ShipRadius    = 10.0
	ShipScale     = 14.0
	ShootCooldown = 10
)

type Ship struct {
	Pos          Vec2
	Vel          Vec2
	Angle        float64
	Thrusting    bool
	Dead         bool
	Invincible   bool
	InvTimer     int
	RespawnTimer int
}

func NewShip() Ship {
	return Ship{
		Pos:   Vec2{ScreenWidth / 2.0, ScreenHeight / 2.0},
		Angle: -math.Pi / 2,
	}
}

func (s *Ship) Update() {
	if s.Invincible {
		s.InvTimer--
		if s.InvTimer <= 0 {
			s.Invincible = false
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		s.Angle -= ShipRotSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		s.Angle += ShipRotSpeed
	}

	s.Thrusting = ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW)
	if s.Thrusting {
		s.Vel.X += math.Cos(s.Angle) * ShipThrust
		s.Vel.Y += math.Sin(s.Angle) * ShipThrust
	}

	if spd := s.Vel.Mag(); spd > ShipMaxSpeed {
		s.Vel = s.Vel.Scale(ShipMaxSpeed / spd)
	}
	s.Vel = s.Vel.Scale(ShipFriction)
	s.Pos = wrapPos(s.Pos.Add(s.Vel), ScreenWidth, ScreenHeight)
	s.Angle = math.Mod(s.Angle, 2*math.Pi)
}

// Hyperspace teleports the ship to a random position with brief invincibility.
func (s *Ship) Hyperspace() {
	s.Pos = Vec2{rand.Float64() * ScreenWidth, rand.Float64() * ScreenHeight}
	s.Vel = Vec2{}
	s.Invincible = true
	s.InvTimer = 90
}

func (s *Ship) Draw(screen *ebiten.Image) {
	if s.Dead {
		return
	}
	// Flash during invincibility frames
	if s.Invincible && (s.InvTimer/5)%2 == 0 {
		return
	}

	// pts[0]=front, pts[1]=back-left, pts[2]=back-notch, pts[3]=back-right
	pts := transformPoints(shipShape, s.Pos, s.Angle, ShipScale)
	strokeLine(screen, pts[0][0], pts[0][1], pts[1][0], pts[1][1], colorShip)
	strokeLine(screen, pts[0][0], pts[0][1], pts[3][0], pts[3][1], colorShip)
	strokeLine(screen, pts[1][0], pts[1][1], pts[2][0], pts[2][1], colorShip)
	strokeLine(screen, pts[3][0], pts[3][1], pts[2][0], pts[2][1], colorShip)

	if s.Thrusting {
		backX := float32(s.Pos.X - math.Cos(s.Angle)*14)
		backY := float32(s.Pos.Y - math.Sin(s.Angle)*14)
		flameX := float32(s.Pos.X - math.Cos(s.Angle)*26)
		flameY := float32(s.Pos.Y - math.Sin(s.Angle)*26)
		strokeLine(screen, backX, backY, flameX, flameY, colorThrust)
	}
}
