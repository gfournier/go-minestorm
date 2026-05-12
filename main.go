package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 600
)

func main() {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("MineStorm — A Vectrex Classic")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	// Prevent the window from shrinking to zero, which causes CAMetalLayer
	// to fail allocating a drawable on macOS.
	ebiten.SetWindowSizeLimits(320, 240, -1, -1)
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
