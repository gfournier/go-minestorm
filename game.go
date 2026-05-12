package main

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

const (
	InitLives        = 3
	RespawnDelay     = 120
	InvincibleFrames = 180
	BonusLifeScore   = 10000
	MaxLives         = 9
)

type GameState int

const (
	StateTitle GameState = iota
	StatePlaying
	StateLevelUp
	StateGameOver
)

type Game struct {
	state     GameState
	ship      Ship
	mines     []Mine
	bullets   []Bullet
	particles []Particle
	score     int
	level     int
	lives     int
	shootCD   int
	lvlTimer  int
	nextLife  int
}

func NewGame() *Game {
	return &Game{state: StateTitle}
}

func (g *Game) Layout(_, _ int) (int, int) { return ScreenWidth, ScreenHeight }

func (g *Game) startGame() {
	g.score = 0
	g.lives = InitLives
	g.level = 1
	g.nextLife = BonusLifeScore
	g.startLevel()
	g.state = StatePlaying
}

func (g *Game) startLevel() {
	g.ship = NewShip()
	g.ship.Invincible = true
	g.ship.InvTimer = InvincibleFrames
	g.bullets = nil
	g.particles = nil
	g.shootCD = 0

	count := 2 + g.level*2
	speedMult := 0.65 + float64(g.level)*0.12
	g.mines = make([]Mine, 0, count)
	for i := 0; i < count; i++ {
		pos := RandomEdgePos()
		g.mines = append(g.mines, NewMine(SizeLarge, pos, g.ship.Pos, speedMult))
	}
}

func (g *Game) Update() error {
	switch g.state {
	case StateTitle:
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.startGame()
		}

	case StatePlaying:
		g.updatePlaying()

	case StateLevelUp:
		// Keep particles ticking so explosions finish animating between levels.
		aliveP := g.particles[:0]
		for i := range g.particles {
			g.particles[i].Update()
			if g.particles[i].Life > 0 {
				aliveP = append(aliveP, g.particles[i])
			}
		}
		g.particles = aliveP
		g.lvlTimer--
		if g.lvlTimer <= 0 {
			g.level++
			g.startLevel()
			g.state = StatePlaying
		}

	case StateGameOver:
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.state = StateTitle
		}
	}
	return nil
}

func (g *Game) updatePlaying() {
	// Ship update / respawn
	if g.ship.Dead {
		g.ship.RespawnTimer--
		if g.ship.RespawnTimer <= 0 {
			g.ship = NewShip()
			g.ship.Invincible = true
			g.ship.InvTimer = InvincibleFrames
		}
	} else {
		g.ship.Update()
		// Hyperspace
		if inpututil.IsKeyJustPressed(ebiten.KeyZ) || inpututil.IsKeyJustPressed(ebiten.KeyH) {
			g.ship.Hyperspace()
		}
	}

	// Shooting
	if g.shootCD > 0 {
		g.shootCD--
	}
	if !g.ship.Dead && ebiten.IsKeyPressed(ebiten.KeySpace) && g.shootCD == 0 {
		g.bullets = append(g.bullets, NewBullet(g.ship))
		g.shootCD = ShootCooldown
	}

	// Bullets
	aliveBullets := g.bullets[:0]
	for i := range g.bullets {
		g.bullets[i].Update()
		if g.bullets[i].Lifetime > 0 {
			aliveBullets = append(aliveBullets, g.bullets[i])
		}
	}
	g.bullets = aliveBullets

	// Mines
	trackStr := 0.003 + 0.001*float64(g.level)
	for i := range g.mines {
		g.mines[i].Update(g.ship.Pos, trackStr)
	}

	g.resolveBulletMine()

	// Ship-mine collision
	if !g.ship.Dead && !g.ship.Invincible {
		for i := range g.mines {
			d := wrappedDelta(g.ship.Pos, g.mines[i].Pos, ScreenWidth, ScreenHeight).Mag()
			if d < ShipRadius+g.mines[i].Radius() {
				g.killShip()
				break
			}
		}
	}

	// Particles
	aliveParticles := g.particles[:0]
	for i := range g.particles {
		g.particles[i].Update()
		if g.particles[i].Life > 0 {
			aliveParticles = append(aliveParticles, g.particles[i])
		}
	}
	g.particles = aliveParticles

	// Level complete
	if len(g.mines) == 0 && !g.ship.Dead {
		g.state = StateLevelUp
		g.lvlTimer = 150
	}
}

func (g *Game) resolveBulletMine() {
	origMines := make([]Mine, len(g.mines))
	copy(origMines, g.mines)

	mineHit := make([]bool, len(origMines))
	bulletHit := make([]bool, len(g.bullets))
	var newMines []Mine

	for mi := range origMines {
		for bi := range g.bullets {
			if bulletHit[bi] {
				continue
			}
			d := wrappedDelta(origMines[mi].Pos, g.bullets[bi].Pos, ScreenWidth, ScreenHeight).Mag()
			if d < origMines[mi].Radius()+BulletRadius {
				mineHit[mi] = true
				bulletHit[bi] = true
				g.score += minePoints[origMines[mi].Size]
				SpawnExplosion(origMines[mi].Pos, origMines[mi].Size, &g.particles)
				if origMines[mi].Size < SizeSmall {
					newMines = append(newMines, g.spawnChildren(origMines[mi])...)
				}
				break
			}
		}
	}

	var aliveMines []Mine
	for i, m := range origMines {
		if !mineHit[i] {
			aliveMines = append(aliveMines, m)
		}
	}
	g.mines = append(aliveMines, newMines...)

	var aliveBullets []Bullet
	for i, b := range g.bullets {
		if !bulletHit[i] {
			aliveBullets = append(aliveBullets, b)
		}
	}
	g.bullets = aliveBullets

	// Award bonus lives
	for g.score >= g.nextLife {
		if g.lives < MaxLives {
			g.lives++
		}
		g.nextLife += BonusLifeScore
	}
}

func (g *Game) spawnChildren(parent Mine) []Mine {
	childSize := parent.Size + 1
	speedMult := 0.65 + float64(g.level)*0.12
	children := make([]Mine, 2)
	for i := range children {
		child := NewMine(childSize, parent.Pos, g.ship.Pos, speedMult)
		// Inherit a fraction of parent velocity for natural spread
		child.Vel.X += parent.Vel.X * 0.3
		child.Vel.Y += parent.Vel.Y * 0.3
		children[i] = child
	}
	return children
}

func (g *Game) killShip() {
	SpawnShipExplosion(g.ship.Pos, &g.particles)
	g.lives--
	g.ship.Dead = true
	if g.lives <= 0 {
		g.state = StateGameOver
	} else {
		g.ship.RespawnTimer = RespawnDelay
	}
}

// ---- Drawing ----

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(colorBg)

	switch g.state {
	case StateTitle:
		g.drawTitle(screen)

	case StatePlaying:
		g.drawPlay(screen)

	case StateLevelUp:
		g.drawPlay(screen)
		drawCenteredText(screen, fmt.Sprintf("LEVEL %d", g.level+1), ScreenWidth/2, ScreenHeight/2-10, colorHUD)
		drawCenteredText(screen, "GET READY!", ScreenWidth/2, ScreenHeight/2+18, colorWhite)

	case StateGameOver:
		g.drawPlay(screen)
		drawCenteredText(screen, "GAME OVER", ScreenWidth/2, ScreenHeight/2-20, colorMineL)
		drawCenteredText(screen, fmt.Sprintf("SCORE  %d", g.score), ScreenWidth/2, ScreenHeight/2+10, colorHUD)
		drawCenteredText(screen, "PRESS SPACE OR ENTER", ScreenWidth/2, ScreenHeight/2+38, colorWhite)
	}
}

func (g *Game) drawPlay(screen *ebiten.Image) {
	for i := range g.particles {
		g.particles[i].Draw(screen)
	}
	for i := range g.bullets {
		g.bullets[i].Draw(screen)
	}
	for i := range g.mines {
		g.mines[i].Draw(screen)
	}
	g.ship.Draw(screen)
	g.drawHUD(screen)
}

func (g *Game) drawHUD(screen *ebiten.Image) {
	f := basicfont.Face7x13
	text.Draw(screen, fmt.Sprintf("SCORE %06d", g.score), f, 10, 18, colorHUD)
	text.Draw(screen, fmt.Sprintf("LVL %d", g.level), f, ScreenWidth-70, 18, colorHUD)

	// Draw small ship icons for remaining lives
	for i := 0; i < g.lives; i++ {
		pos := Vec2{float64(12 + i*22), 38}
		pts := transformPoints(shipShape, pos, -math.Pi/2, 9)
		strokeLine(screen, pts[0][0], pts[0][1], pts[1][0], pts[1][1], colorShip)
		strokeLine(screen, pts[0][0], pts[0][1], pts[3][0], pts[3][1], colorShip)
		strokeLine(screen, pts[1][0], pts[1][1], pts[2][0], pts[2][1], colorShip)
		strokeLine(screen, pts[3][0], pts[3][1], pts[2][0], pts[2][1], colorShip)
	}
}

func (g *Game) drawTitle(screen *ebiten.Image) {
	drawCenteredText(screen, "M I N E S T O R M", ScreenWidth/2, ScreenHeight/2-80, colorMineL)
	drawCenteredText(screen, "A VECTREX CLASSIC CLONE", ScreenWidth/2, ScreenHeight/2-50, colorMineM)

	drawCenteredText(screen, "ARROWS / WASD    ROTATE & THRUST", ScreenWidth/2, ScreenHeight/2+10, colorWhite)
	drawCenteredText(screen, "SPACE            FIRE", ScreenWidth/2, ScreenHeight/2+32, colorWhite)
	drawCenteredText(screen, "Z / H            HYPERSPACE", ScreenWidth/2, ScreenHeight/2+54, colorWhite)

	drawCenteredText(screen, "PRESS SPACE OR ENTER TO START", ScreenWidth/2, ScreenHeight/2+100, colorHUD)
}
