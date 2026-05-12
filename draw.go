package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

var (
	colorBg     = color.RGBA{0, 0, 0, 255}
	colorShip   = color.RGBA{160, 220, 255, 255}
	colorThrust = color.RGBA{255, 100, 0, 255}
	colorBullet = color.RGBA{255, 255, 80, 255}
	colorMineL  = color.RGBA{255, 70, 70, 255}
	colorMineM  = color.RGBA{255, 165, 0, 255}
	colorMineS  = color.RGBA{200, 255, 60, 255}
	colorHUD    = color.RGBA{60, 255, 100, 255}
	colorWhite  = color.RGBA{255, 255, 255, 255}
)

func mineColor(size MineSize) color.RGBA {
	switch size {
	case SizeLarge:
		return colorMineL
	case SizeMedium:
		return colorMineM
	default:
		return colorMineS
	}
}

func strokeLine(screen *ebiten.Image, x0, y0, x1, y1 float32, clr color.RGBA) {
	vector.StrokeLine(screen, x0, y0, x1, y1, 1.5, clr, true)
}

func drawPoly(screen *ebiten.Image, pts [][2]float32, clr color.RGBA) {
	n := len(pts)
	for i := 0; i < n; i++ {
		p0, p1 := pts[i], pts[(i+1)%n]
		strokeLine(screen, p0[0], p0[1], p1[0], p1[1], clr)
	}
}

func transformPoints(shape [][2]float64, pos Vec2, angle, scale float64) [][2]float32 {
	cos, sin := math.Cos(angle), math.Sin(angle)
	out := make([][2]float32, len(shape))
	for i, p := range shape {
		x := p[0]*cos - p[1]*sin
		y := p[0]*sin + p[1]*cos
		out[i] = [2]float32{float32(x*scale+pos.X), float32(y*scale+pos.Y)}
	}
	return out
}

// starPolygon builds a star/spike outline with n spikes.
// outerR = spike-tip radius, innerR = valley radius (both normalized 0..1).
func starPolygon(n int, outerR, innerR float64) [][2]float64 {
	pts := make([][2]float64, n*2)
	step := math.Pi / float64(n)
	for i := 0; i < n; i++ {
		a := float64(i)*2*math.Pi/float64(n) - math.Pi/2
		pts[i*2] = [2]float64{math.Cos(a) * outerR, math.Sin(a) * outerR}
		pts[i*2+1] = [2]float64{math.Cos(a+step) * innerR, math.Sin(a+step) * innerR}
	}
	return pts
}

// Ship: front-tip, back-left, back-notch, back-right (pointing right when angle=0)
var shipShape = [][2]float64{
	{1.0, 0},
	{-0.75, -0.6},
	{-0.3, 0},
	{-0.75, 0.6},
}

var (
	mineLargeShape  = starPolygon(8, 1.0, 0.5)
	mineMediumShape = starPolygon(6, 1.0, 0.5)
	mineSmallShape  = starPolygon(4, 1.0, 0.5)
)

func getMineShape(size MineSize) [][2]float64 {
	switch size {
	case SizeLarge:
		return mineLargeShape
	case SizeMedium:
		return mineMediumShape
	default:
		return mineSmallShape
	}
}

func drawCenteredText(screen *ebiten.Image, s string, cx, cy int, clr color.RGBA) {
	w := len(s) * 7
	text.Draw(screen, s, basicfont.Face7x13, cx-w/2, cy, clr)
}
