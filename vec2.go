package main

import "math"

type Vec2 struct{ X, Y float64 }

func (v Vec2) Add(u Vec2) Vec2      { return Vec2{v.X + u.X, v.Y + u.Y} }
func (v Vec2) Sub(u Vec2) Vec2      { return Vec2{v.X - u.X, v.Y - u.Y} }
func (v Vec2) Scale(s float64) Vec2 { return Vec2{v.X * s, v.Y * s} }
func (v Vec2) Mag() float64         { return math.Sqrt(v.X*v.X + v.Y*v.Y) }

func (v Vec2) Norm() Vec2 {
	m := v.Mag()
	if m < 1e-9 {
		return Vec2{}
	}
	return Vec2{v.X / m, v.Y / m}
}

func wrapPos(pos Vec2, w, h float64) Vec2 {
	x, y := pos.X, pos.Y
	if x < 0 {
		x += w
	} else if x >= w {
		x -= w
	}
	if y < 0 {
		y += h
	} else if y >= h {
		y -= h
	}
	return Vec2{x, y}
}

// wrappedDelta returns the shortest vector from `from` to `to` in a wrapped world.
func wrappedDelta(from, to Vec2, w, h float64) Vec2 {
	dx := to.X - from.X
	dy := to.Y - from.Y
	if dx > w/2 {
		dx -= w
	} else if dx < -w/2 {
		dx += w
	}
	if dy > h/2 {
		dy -= h
	} else if dy < -h/2 {
		dy += h
	}
	return Vec2{dx, dy}
}

func rotateVec(v Vec2, angle float64) Vec2 {
	cos, sin := math.Cos(angle), math.Sin(angle)
	return Vec2{v.X*cos - v.Y*sin, v.X*sin + v.Y*cos}
}
