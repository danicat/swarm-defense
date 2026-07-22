package main

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type GameState int

const (
	StateIntro GameState = iota
	StateTitle
	StateBuild
	StateAttack
	StateWin
	StateOver
)

type CatalogItem struct {
	Name       string
	Cost       int
	IsBuilding bool
	Type       int
}

var Catalog = []CatalogItem{
	{"ARROW", 100, true, int(BldArrow)},
	{"CANNON", 200, true, int(BldCannon)},
	{"ICE", 150, true, int(BldIce)},
	{"MINE", 120, true, int(BldGoldMine)},
	{"KNIGHT", 120, false, int(UnitKnight)},
	{"ARCHER", 100, false, int(UnitArcher)},
	{"MAGE", 150, false, int(UnitMage)},
	{"HEALER", 110, false, int(UnitHealer)},
}

type Game struct {
	// State Machine
	state       GameState
	introTimer  time.Time
	titleTicks  int
	winTicks    int
	overTicks   int

	// Resources & Metrics
	gold       int
	lives      int
	wave       int
	score      int
	highScore  int
	baseHealth int // Base health at coordinates (19, 5) or similar

	// Audio System
	audioContext *audio.Context
	sounds       *SoundSystem

	// Assets
	sprites *SpriteSystem

	// Entity Containers
	units       []*Unit
	buildings   []*Building
	enemies     []*Enemy
	projectiles []*Projectile
	particles   []*Particle

	// Level Map Grid (20x12)
	// 0: Empty, 1: Path, 2: Obstacle, 3: Building, 4: Enemy Base (Spawn), 5: Player Base
	grid [20][12]int

	// Selection & UI State
	selectedBuildType int // 0-3 for buildings
	selectedUnitType  int // 0-3 for units
	cursorX, cursorY  int // grid cursor coordinates
	isPlacingBuilding bool
	isPlacingUnit     bool

	// Wave Spawning Controller
	waveActive    bool
	spawnTimer    float64
	spawnQueue    []EnemyType
	waveTimer     float64
}

func NewGame() *Game {
	g := &Game{
		state:      StateIntro,
		introTimer: time.Now(),
		gold:       500,
		lives:      20,
		wave:       1,
		baseHealth: 100,
	}

	// Initialize Audio Context
	g.audioContext = audio.NewContext(44100)
	g.sounds = NewSoundSystem(g.audioContext)

	// Initialize Sprite System
	g.sprites = NewSpriteSystem()

	// High score
	g.highScore = LoadHighScore()

	// Default Grid Path (0,5 to 19,5)
	for x := 0; x < 20; x++ {
		for y := 0; y < 12; y++ {
			g.grid[x][y] = 0 // Empty grass
		}
	}
	
	// Create a winding path starting at 0,5 and ending at 19,5
	// Let's make it wind a bit
	// 0,5 -> 4,5 -> 4,2 -> 10,2 -> 10,9 -> 15,9 -> 15,5 -> 19,5
	for x := 0; x <= 4; x++ { g.grid[x][5] = 1 }
	for y := 2; y <= 5; y++ { g.grid[4][y] = 1 }
	for x := 4; x <= 10; x++ { g.grid[x][2] = 1 }
	for y := 2; y <= 9; y++ { g.grid[10][y] = 1 }
	for x := 10; x <= 15; x++ { g.grid[x][9] = 1 }
	for y := 5; y <= 9; y++ { g.grid[15][y] = 1 }
	for x := 15; x <= 19; x++ { g.grid[x][5] = 1 }

	g.grid[0][5] = 4 // Enemy Base
	g.grid[19][5] = 5 // Player Base

	return g
}

func (g *Game) Update() error {
	switch g.state {
	case StateIntro:
		if time.Since(g.introTimer).Seconds() > 3 {
			g.state = StateTitle
		}
	case StateTitle:
		g.titleTicks++
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			g.state = StateBuild
			g.waveTimer = 30 // 30 seconds for build phase
		}
	case StateBuild:
		g.updateBuildPhase()
	case StateAttack:
		g.updateAttackPhase()
	case StateWin:
		g.winTicks++
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			g.resetGame()
		}
	case StateOver:
		g.overTicks++
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			g.resetGame()
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	DrawUI(screen, g)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func (g *Game) updateBuildPhase() {
	g.waveTimer -= 1.0 / 60.0
	
	// Input logic
	g.handleInput()

	// Update gold mines or other buildings during build phase so mining works
	for _, b := range g.buildings {
		if b.Type == BldGoldMine {
			b.Update(g.enemies, g.units, func(gold int) {
				g.gold += gold
				g.score += gold / 2
				g.sounds.PlayGoldMineClick()
			}, func(p *Projectile) {})
		}
	}

	// Update particles during build phase
	var activeParticles []*Particle
	for _, part := range g.particles {
		dead := part.Update()
		if !dead {
			activeParticles = append(activeParticles, part)
		}
	}
	g.particles = activeParticles

	if g.waveTimer <= 0 || ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.startWave()
	}
}

func (g *Game) updateAttackPhase() {
	// Spawning
	if len(g.spawnQueue) > 0 {
		g.spawnTimer -= 1.0 / 60.0
		if g.spawnTimer <= 0 {
			g.spawnEnemy(g.spawnQueue[0])
			g.spawnQueue = g.spawnQueue[1:]
			g.spawnTimer = 1.0 // 1 second between spawns
		}
	} else if len(g.enemies) == 0 && len(g.spawnQueue) == 0 {
		// Wave complete
		g.wave++
		if g.wave > 10 {
			g.state = StateWin
			g.sounds.PlayWin()
			SaveHighScore(g.score)
		} else {
			g.state = StateBuild
			g.waveTimer = 30
		}
	}

	// 1. Update Units
	for _, u := range g.units {
		u.Update(g.enemies, g.buildings, g.units, func(p *Projectile) {
			g.projectiles = append(g.projectiles, p)
		})
	}

	// 2. Update Buildings
	for _, b := range g.buildings {
		b.Update(g.enemies, g.units, func(gold int) {
			g.gold += gold
			g.score += gold / 2
			g.sounds.PlayGoldMineClick()
		}, func(p *Projectile) {
			g.projectiles = append(g.projectiles, p)
			g.sounds.PlayArrowShoot()
		})
	}

	// 3. Update Enemies
	var activeEnemies []*Enemy
	for _, e := range g.enemies {
		if e.Path == nil {
			e.Path = FindPath(&g.grid, 0, 5, 19, 5)
		}

		reachedBase := e.Update(19*32+16, 5*32+16, func(newEnemy *Enemy) {
			activeEnemies = append(activeEnemies, newEnemy)
		}, func(p *Projectile) {
			g.projectiles = append(g.projectiles, p)
		}, g.units, func(part *Particle) {
			g.particles = append(g.particles, part)
		})

		if reachedBase {
			g.baseHealth -= int(e.Damage)
			g.sounds.PlayBaseHurt()
			if g.baseHealth <= 0 {
				g.lives--
				g.baseHealth = 100
				if g.lives <= 0 {
					g.state = StateOver
					SaveHighScore(g.score)
					g.sounds.PlayGameOver()
				}
			}
		} else if e.Health <= 0 {
			g.gold += e.GoldReward
			g.score += e.GoldReward * 5
			g.sounds.PlayEnemyDeath()
			for i := 0; i < 15; i++ {
				g.particles = append(g.particles, NewParticle(e.X, e.Y, 0xFFFFFFFF))
			}
		} else {
			activeEnemies = append(activeEnemies, e)
		}
	}
	g.enemies = activeEnemies

	// 4. Update Projectiles
	var activeProjectiles []*Projectile
	for _, p := range g.projectiles {
		hit := p.Update(g.enemies, g.units, func(part *Particle) {
			g.particles = append(g.particles, part)
		})
		if hit {
			if p.Type == ProjCannonBall {
				g.sounds.PlayExplosion()
			} else if p.Type == ProjIceShard {
				g.sounds.PlayIceFreeze()
			} else {
				g.sounds.PlayUnitHurt()
			}
		} else {
			activeProjectiles = append(activeProjectiles, p)
		}
	}
	g.projectiles = activeProjectiles

	// 5. Update Particles
	var activeParticles []*Particle
	for _, part := range g.particles {
		dead := part.Update()
		if !dead {
			activeParticles = append(activeParticles, part)
		}
	}
	g.particles = activeParticles
}

func (g *Game) handleInput() {
	mx, my := ebiten.CursorPosition()
	g.cursorX = mx / 32
	g.cursorY = my / 32

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if my >= 384 {
			catalogX := 160
			catalogY := 400
			itemSize := 32
			
			for i := 0; i < 8; i++ {
				x1 := catalogX + i*(itemSize+12)
				x2 := x1 + itemSize
				y1 := catalogY
				y2 := y1 + itemSize
				
				if mx >= x1 && mx <= x2 && my >= y1 && my <= y2 {
					g.sounds.PlayGoldMineClick()
					if i <= 3 {
						g.selectedBuildType = i
						g.isPlacingBuilding = true
						g.isPlacingUnit = false
					} else {
						g.selectedUnitType = i - 4
						g.isPlacingBuilding = false
						g.isPlacingUnit = true
					}
					return
				}
			}
		} else {
			if g.cursorX >= 0 && g.cursorX < 20 && g.cursorY >= 0 && g.cursorY < 12 {
				currentTile := g.grid[g.cursorX][g.cursorY]
				
				if g.isPlacingBuilding {
					cost := Catalog[g.selectedBuildType].Cost
					if g.gold >= cost && currentTile == 0 {
						g.gold -= cost
						g.buildings = append(g.buildings, NewBuilding(BuildingType(g.selectedBuildType), g.cursorX, g.cursorY))
						g.grid[g.cursorX][g.cursorY] = 3
						g.sounds.PlayExplosion()
					}
				} else if g.isPlacingUnit {
					cost := Catalog[g.selectedUnitType+4].Cost
					if g.gold >= cost && currentTile == 0 {
						g.gold -= cost
						g.units = append(g.units, NewUnit(UnitType(g.selectedUnitType), g.cursorX, g.cursorY))
						g.sounds.PlayUnitHurt()
					}
				}
			}
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		g.selectedBuildType = int(BldArrow)
		g.isPlacingBuilding = true
		g.isPlacingUnit = false
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		g.selectedBuildType = int(BldCannon)
		g.isPlacingBuilding = true
		g.isPlacingUnit = false
	}
	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		g.selectedBuildType = int(BldIce)
		g.isPlacingBuilding = true
		g.isPlacingUnit = false
	}
	if inpututil.IsKeyJustPressed(ebiten.Key4) {
		g.selectedBuildType = int(BldGoldMine)
		g.isPlacingBuilding = true
		g.isPlacingUnit = false
	}
	if inpututil.IsKeyJustPressed(ebiten.Key5) {
		g.selectedUnitType = int(UnitKnight)
		g.isPlacingBuilding = false
		g.isPlacingUnit = true
	}
	if inpututil.IsKeyJustPressed(ebiten.Key6) {
		g.selectedUnitType = int(UnitArcher)
		g.isPlacingBuilding = false
		g.isPlacingUnit = true
	}
	if inpututil.IsKeyJustPressed(ebiten.Key7) {
		g.selectedUnitType = int(UnitMage)
		g.isPlacingBuilding = false
		g.isPlacingUnit = true
	}
	if inpututil.IsKeyJustPressed(ebiten.Key8) {
		g.selectedUnitType = int(UnitHealer)
		g.isPlacingBuilding = false
		g.isPlacingUnit = true
	}
}

func (g *Game) startWave() {
	g.state = StateAttack
	g.waveActive = true
	g.spawnQueue = g.generateWave(g.wave)
	g.spawnTimer = 1.0
}

func (g *Game) generateWave(waveNum int) []EnemyType {
	var queue []EnemyType
	if waveNum == 10 {
		queue = append(queue, EnemyDragonBoss)
		return queue
	}
	count := waveNum * 5
	for i := 0; i < count; i++ {
		queue = append(queue, EnemyGoblin)
	}
	return queue
}

func (g *Game) spawnEnemy(t EnemyType) {
	e := NewEnemy(t, 0*32+16, 5*32+16)
	g.enemies = append(g.enemies, e)
}

func (g *Game) resetGame() {
	g.state = StateIntro
	g.introTimer = time.Now()
	g.gold = 500
	g.lives = 20
	g.wave = 1
	g.score = 0
	g.enemies = nil
	g.buildings = nil
	g.units = nil
	g.projectiles = nil
	g.particles = nil
	// Reset grid...
}
