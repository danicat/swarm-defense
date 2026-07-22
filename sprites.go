package main

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// SpriteSystem manages all cached in-memory procedural Ebitengine images.
type SpriteSystem struct {
	Tiles       [10]*ebiten.Image
	Units       map[UnitType][4][8]*ebiten.Image   // [Type][Direction][Frame]
	Buildings   map[BuildingType]*ebiten.Image
	Enemies     map[EnemyType][4][8]*ebiten.Image  // [Type][Direction][Frame]
	Projectiles map[ProjectileType]*ebiten.Image
}

// NewSpriteSystem compiles and returns the full SpriteSystem struct populated with cacheable frames.
func NewSpriteSystem() *SpriteSystem {
	sys := &SpriteSystem{
		Units:       make(map[UnitType][4][8]*ebiten.Image),
		Buildings:   make(map[BuildingType]*ebiten.Image),
		Enemies:     make(map[EnemyType][4][8]*ebiten.Image),
		Projectiles: make(map[ProjectileType]*ebiten.Image),
	}

	// 1. Generate Tiles (Grass, Path, Obstacle, Building base, Enemy base, Player base, etc.)
	sys.Tiles[0] = makeGrassTile()
	sys.Tiles[1] = makePathTile()
	sys.Tiles[2] = makeObstacleTile()
	sys.Tiles[3] = makeBuildingBaseTile()
	sys.Tiles[4] = makeEnemyBaseTile()
	sys.Tiles[5] = makePlayerBaseTile()
	sys.Tiles[6] = makeWaterTile()
	sys.Tiles[7] = makeLavaTile()
	sys.Tiles[8] = makeSandTile()
	sys.Tiles[9] = makeSnowTile()

	// 2. Generate Buildings
	buildingTypes := []BuildingType{BldArrow, BldCannon, BldIce, BldGoldMine}
	for _, bt := range buildingTypes {
		sys.Buildings[bt] = makeBuildingSprite(bt)
	}

	// 3. Generate Projectiles
	projTypes := []ProjectileType{ProjArrow, ProjCannonBall, ProjIceShard, ProjFireBall, ProjMagicBall, ProjHealBeam}
	for _, pt := range projTypes {
		sys.Projectiles[pt] = makeProjectileSprite(pt)
	}

	// 4. Generate Units (Knight, Archer, Mage, Healer)
	unitTypes := []UnitType{UnitKnight, UnitArcher, UnitMage, UnitHealer}
	for _, ut := range unitTypes {
		var dirArr [4][8]*ebiten.Image
		for d := 0; d < 4; d++ {
			for f := 0; f < 8; f++ {
				dirArr[d][f] = makeUnitSprite(ut, d, f)
			}
		}
		sys.Units[ut] = dirArr
	}

	// 5. Generate Enemies (Goblin, Orc, Troll, Bat, FireElem, Assassin, Necro, Skeleton, DragonBoss)
	enemyTypes := []EnemyType{
		EnemyGoblin, EnemyOrc, EnemyTroll, EnemyBat, EnemyFireElem,
		EnemyAssassin, EnemyNecro, EnemySkeleton, EnemyDragonBoss,
	}
	for _, et := range enemyTypes {
		var dirArr [4][8]*ebiten.Image
		for d := 0; d < 4; d++ {
			for f := 0; f < 8; f++ {
				dirArr[d][f] = makeEnemySprite(et, d, f)
			}
		}
		sys.Enemies[et] = dirArr
	}

	return sys
}

// --- TILE GENERATORS ---

func makeGrassTile() *ebiten.Image {
	img := image.NewRGBA(image.Rect(0, 0, 32, 32))
	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			noise := float64(x*x + y*y*7)
			n := int(math.Mod(noise, 15)) - 7
			r := 40 + n/2
			g := 120 + n
			b := 40 + n/3
			img.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), 255})
		}
	}
	// Draw cute little grass blade tufts
	drawTuft(img, 5, 8)
	drawTuft(img, 18, 5)
	drawTuft(img, 12, 18)
	drawTuft(img, 24, 22)
	drawTuft(img, 8, 25)
	return ebiten.NewImageFromImage(img)
}

func drawTuft(img *image.RGBA, cx, cy int) {
	c := color.RGBA{60, 160, 80, 255}
	c2 := color.RGBA{80, 200, 100, 255}
	img.Set(cx, cy, c)
	img.Set(cx-1, cy, c)
	img.Set(cx+1, cy, c)
	img.Set(cx, cy-1, c)
	img.Set(cx, cy-2, c2)
	img.Set(cx-1, cy-1, c2)
	img.Set(cx+1, cy-1, c2)
}

func makePathTile() *ebiten.Image {
	img := image.NewRGBA(image.Rect(0, 0, 32, 32))
	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			noise := float64(x*x*3 + y*y*11)
			n := int(math.Mod(noise, 20)) - 10
			r := 180 + n
			g := 150 + n
			b := 110 + n
			img.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), 255})
		}
	}
	// Draw path gravel textures
	drawPebble(img, 4, 10, 3)
	drawPebble(img, 20, 6, 2)
	drawPebble(img, 12, 22, 3)
	drawPebble(img, 26, 18, 2)
	drawPebble(img, 8, 28, 4)
	return ebiten.NewImageFromImage(img)
}

func drawPebble(img *image.RGBA, cx, cy, size int) {
	col := color.RGBA{110, 95, 80, 255}
	colHigh := color.RGBA{140, 125, 110, 255}
	for dy := 0; dy < size; dy++ {
		for dx := 0; dx < size; dx++ {
			if dx == 0 && dy == 0 {
				img.Set(cx+dx, cy+dy, colHigh)
			} else {
				img.Set(cx+dx, cy+dy, col)
			}
		}
	}
}

func makeObstacleTile() *ebiten.Image {
	img := image.NewRGBA(image.Rect(0, 0, 32, 32))
	// Draw grass background first
	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			noise := float64(x*x + y*y*7)
			n := int(math.Mod(noise, 15)) - 7
			img.Set(x, y, color.RGBA{uint8(40 + n/2), uint8(120 + n), uint8(40 + n/3), 255})
		}
	}
	// Ground shadow for rock
	fillEllipse(img, 16, 23, 10, 5, color.RGBA{10, 40, 15, 180})
	// Boulder grey
	fillCircle(img, 16, 16, 9, color.RGBA{120, 120, 125, 255})
	// Shading (bottom/right)
	fillEllipse(img, 18, 18, 7, 6, color.RGBA{90, 90, 95, 255})
	// Highlight (top/left)
	fillCircle(img, 13, 13, 5, color.RGBA{160, 160, 165, 255})
	fillCircle(img, 11, 11, 2, color.RGBA{210, 210, 215, 255})

	applyOutline(img, color.RGBA{20, 20, 25, 255})
	return ebiten.NewImageFromImage(img)
}

func makeBuildingBaseTile() *ebiten.Image {
	img := image.NewRGBA(image.Rect(0, 0, 32, 32))
	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			bx, by := x/8, y/8
			col := color.RGBA{80, 85, 95, 255}
			if (bx+by)%2 == 0 {
				col = color.RGBA{65, 70, 80, 255}
			}
			if x%8 == 0 || y%8 == 0 {
				col = color.RGBA{35, 35, 40, 255}
			}
			img.Set(x, y, col)
		}
	}
	return ebiten.NewImageFromImage(img)
}

func makeEnemyBaseTile() *ebiten.Image {
	img := image.NewRGBA(image.Rect(0, 0, 32, 32))
	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			noise := float64(x*x*5 + y*y*13)
			n := int(math.Mod(noise, 20)) - 10
			img.Set(x, y, color.RGBA{uint8(25 + n/3), uint8(20 + n/3), uint8(35 + n/2), 255})
		}
	}
	// Dark purple portals
	fillCircle(img, 16, 16, 11, color.RGBA{60, 20, 90, 255})
	fillCircle(img, 16, 16, 8, color.RGBA{100, 30, 150, 255})
	fillCircle(img, 16, 16, 5, color.RGBA{160, 50, 220, 255})
	fillCircle(img, 16, 16, 2, color.RGBA{230, 150, 255, 255})

	// Corner pillars
	drawPillar(img, 2, 2)
	drawPillar(img, 26, 2)
	drawPillar(img, 2, 26)
	drawPillar(img, 26, 26)
	return ebiten.NewImageFromImage(img)
}

func drawPillar(img *image.RGBA, x, y int) {
	col := color.RGBA{100, 100, 110, 255}
	colHigh := color.RGBA{130, 130, 140, 255}
	colDark := color.RGBA{60, 60, 65, 255}
	for dy := 0; dy < 4; dy++ {
		for dx := 0; dx < 4; dx++ {
			if dx == 0 || dy == 0 {
				img.Set(x+dx, y+dy, colHigh)
			} else if dx == 3 || dy == 3 {
				img.Set(x+dx, y+dy, colDark)
			} else {
				img.Set(x+dx, y+dy, col)
			}
		}
	}
}

func makePlayerBaseTile() *ebiten.Image {
	img := image.NewRGBA(image.Rect(0, 0, 32, 32))
	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			noise := float64(x*x*7 + y*y*17)
			n := int(math.Mod(noise, 15)) - 7
			img.Set(x, y, color.RGBA{uint8(150 + n), uint8(130 + n), uint8(50 + n/2), 255})
		}
	}
	// Brick lines
	for y := 0; y < 32; y += 8 {
		for x := 0; x < 32; x++ {
			if y > 0 {
				img.Set(x, y, color.RGBA{90, 80, 30, 255})
			}
			offset := 0
			if (y/8)%2 == 1 {
				offset = 8
			}
			if (x-offset)%16 == 0 {
				for dy := 0; dy < 8; dy++ {
					if y+dy < 32 {
						img.Set(x, y+dy, color.RGBA{90, 80, 30, 255})
					}
				}
			}
		}
	}
	// Royal Crest
	fillCircle(img, 16, 16, 9, color.RGBA{20, 60, 150, 255})
	fillCircle(img, 16, 16, 6, color.RGBA{40, 100, 220, 255})
	fillRect(img, 15, 12, 2, 8, color.RGBA{240, 200, 40, 255})
	fillRect(img, 12, 15, 8, 2, color.RGBA{240, 200, 40, 255})
	return ebiten.NewImageFromImage(img)
}

func makeWaterTile() *ebiten.Image {
	img := image.NewRGBA(image.Rect(0, 0, 32, 32))
	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			noise := float64(x*x*3 + y*y*5)
			n := int(math.Mod(noise, 10))
			img.Set(x, y, color.RGBA{uint8(30 + n), uint8(110 + n), uint8(220 + n/2), 255})
		}
	}
	for i := 4; i < 32; i += 8 {
		drawLine(img, i, i, i+3, i, color.RGBA{130, 210, 255, 255})
	}
	return ebiten.NewImageFromImage(img)
}

func makeLavaTile() *ebiten.Image {
	img := image.NewRGBA(image.Rect(0, 0, 32, 32))
	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			noise := float64(x*x*7 + y*y*3)
			n := int(math.Mod(noise, 15))
			img.Set(x, y, color.RGBA{uint8(200 + n), uint8(50 + n/2), 10, 255})
		}
	}
	for i := 2; i < 32; i += 10 {
		drawLine(img, i, 31-i, i+4, 31-(i+4), color.RGBA{255, 210, 40, 255})
	}
	return ebiten.NewImageFromImage(img)
}

func makeSandTile() *ebiten.Image {
	img := image.NewRGBA(image.Rect(0, 0, 32, 32))
	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			noise := float64(x*x*11 + y*y*13)
			n := int(math.Mod(noise, 10)) - 5
			img.Set(x, y, color.RGBA{uint8(225 + n), uint8(200 + n), uint8(140 + n/2), 255})
		}
	}
	return ebiten.NewImageFromImage(img)
}

func makeSnowTile() *ebiten.Image {
	img := image.NewRGBA(image.Rect(0, 0, 32, 32))
	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			noise := float64(x*x*3 + y*y*7)
			n := int(math.Mod(noise, 10)) - 5
			img.Set(x, y, color.RGBA{uint8(240 + n), uint8(245 + n), 255, 255})
		}
	}
	return ebiten.NewImageFromImage(img)
}

// --- BUILDING GENERATORS ---

func makeBuildingSprite(t BuildingType) *ebiten.Image {
	img := image.NewRGBA(image.Rect(0, 0, 32, 32))
	// Base is building slate (Tile 3)
	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			bx, by := x/8, y/8
			col := color.RGBA{80, 85, 95, 255}
			if (bx+by)%2 == 0 {
				col = color.RGBA{65, 70, 80, 255}
			}
			if x%8 == 0 || y%8 == 0 {
				col = color.RGBA{35, 35, 40, 255}
			}
			img.Set(x, y, col)
		}
	}

	// Tower shadow
	fillEllipse(img, 16, 26, 9, 4, color.RGBA{10, 10, 15, 160})

	switch t {
	case BldArrow:
		// Wooden Tower with Ballista
		fillRect(img, 11, 8, 3, 18, color.RGBA{110, 70, 35, 255})
		fillRect(img, 18, 8, 3, 18, color.RGBA{110, 70, 35, 255})
		fillRect(img, 14, 12, 4, 3, color.RGBA{130, 85, 45, 255})
		fillRect(img, 14, 18, 4, 3, color.RGBA{130, 85, 45, 255})
		fillRect(img, 8, 6, 16, 4, color.RGBA{80, 50, 25, 255})
		drawLine(img, 10, 4, 22, 4, color.RGBA{160, 160, 160, 255})
		fillRect(img, 15, 2, 2, 7, color.RGBA{100, 65, 30, 255})
		drawLine(img, 16, 1, 16, 5, color.RGBA{240, 210, 40, 255})

	case BldCannon:
		// Heavy stone fort with dark cannon pointing up-right
		fillCircle(img, 16, 18, 9, color.RGBA{120, 120, 125, 255})
		fillCircle(img, 16, 15, 9, color.RGBA{140, 140, 145, 255})
		fillRect(img, 7, 6, 18, 5, color.RGBA{100, 100, 105, 255})
		fillRect(img, 7, 3, 3, 3, color.RGBA{100, 100, 105, 255})
		fillRect(img, 15, 3, 3, 3, color.RGBA{100, 100, 105, 255})
		fillRect(img, 22, 3, 3, 3, color.RGBA{100, 100, 105, 255})
		drawLine(img, 16, 14, 24, 6, color.RGBA{30, 30, 35, 255})
		drawLine(img, 15, 14, 23, 5, color.RGBA{50, 50, 55, 255})
		drawLine(img, 16, 15, 25, 7, color.RGBA{50, 50, 55, 255})
		img.Set(24, 6, color.RGBA{230, 170, 30, 255})
		img.Set(23, 5, color.RGBA{230, 170, 30, 255})
		img.Set(25, 7, color.RGBA{230, 170, 30, 255})

	case BldIce:
		// Pale blue crystalline spire
		fillCircle(img, 16, 20, 9, color.RGBA{40, 100, 150, 255})
		fillCircle(img, 16, 20, 7, color.RGBA{80, 180, 230, 255})
		for y := 4; y < 20; y++ {
			width := (y - 4) / 2
			if width < 1 {
				width = 1
			}
			if width > 6 {
				width = 6
			}
			fillRect(img, 16-width, y, width*2, 1, color.RGBA{120, 220, 255, 255})
			fillRect(img, 16-width, y, width, 1, color.RGBA{200, 245, 255, 255})
		}
		img.Set(8, 8, color.RGBA{200, 255, 255, 255})
		img.Set(24, 10, color.RGBA{200, 255, 255, 255})
		img.Set(10, 16, color.RGBA{150, 230, 255, 255})
		img.Set(22, 15, color.RGBA{150, 230, 255, 255})

	case BldGoldMine:
		// Wooden shaft with sparkling nuggets
		fillRect(img, 6, 8, 20, 4, color.RGBA{100, 60, 30, 255})
		fillRect(img, 6, 12, 4, 15, color.RGBA{100, 60, 30, 255})
		fillRect(img, 22, 12, 4, 15, color.RGBA{100, 60, 30, 255})
		fillRect(img, 10, 12, 12, 15, color.RGBA{25, 15, 10, 255})
		drawLine(img, 12, 25, 20, 25, color.RGBA{120, 120, 120, 255})
		drawLine(img, 12, 23, 12, 27, color.RGBA{80, 50, 25, 255})
		drawLine(img, 16, 23, 16, 27, color.RGBA{80, 50, 25, 255})
		drawLine(img, 20, 23, 20, 27, color.RGBA{80, 50, 25, 255})
		fillCircle(img, 26, 24, 4, color.RGBA{220, 170, 20, 255})
		img.Set(25, 23, color.RGBA{255, 230, 100, 255})
		img.Set(27, 24, color.RGBA{255, 230, 100, 255})
	}

	applyOutline(img, color.RGBA{15, 15, 20, 255})
	return ebiten.NewImageFromImage(img)
}

// --- PROJECTILE GENERATORS ---

func makeProjectileSprite(t ProjectileType) *ebiten.Image {
	img := image.NewRGBA(image.Rect(0, 0, 32, 32))

	switch t {
	case ProjArrow:
		drawLine(img, 8, 16, 24, 16, color.RGBA{120, 80, 40, 255})
		img.Set(24, 16, color.RGBA{190, 190, 200, 255})
		img.Set(23, 15, color.RGBA{130, 130, 140, 255})
		img.Set(23, 17, color.RGBA{130, 130, 140, 255})
		img.Set(8, 15, color.RGBA{240, 240, 240, 255})
		img.Set(8, 17, color.RGBA{240, 240, 240, 255})
		img.Set(9, 14, color.RGBA{220, 220, 220, 255})
		img.Set(9, 18, color.RGBA{220, 220, 220, 255})

	case ProjCannonBall:
		fillCircle(img, 16, 16, 5, color.RGBA{45, 45, 50, 255})
		fillCircle(img, 15, 15, 3, color.RGBA{70, 70, 75, 255})
		img.Set(14, 14, color.RGBA{120, 120, 125, 255})
		drawLine(img, 18, 13, 21, 10, color.RGBA{110, 80, 50, 255})
		img.Set(21, 10, color.RGBA{255, 230, 100, 255})
		img.Set(22, 9, color.RGBA{255, 100, 30, 255})

	case ProjIceShard:
		for x := 8; x <= 24; x++ {
			width := (x - 8) / 3
			if width < 1 {
				width = 1
			}
			if width > 4 {
				width = 4
			}
			fillRect(img, 32-x, 16-width/2, 1, width, color.RGBA{100, 220, 255, 255})
			img.Set(32-x, 16, color.RGBA{200, 250, 255, 255})
		}

	case ProjFireBall:
		fillCircle(img, 18, 16, 5, color.RGBA{255, 80, 20, 255})
		fillCircle(img, 19, 16, 3, color.RGBA{255, 180, 50, 255})
		img.Set(20, 16, color.RGBA{255, 255, 180, 255})
		drawLine(img, 14, 16, 8, 16, color.RGBA{255, 80, 20, 255})
		drawLine(img, 15, 14, 10, 13, color.RGBA{230, 50, 10, 255})
		drawLine(img, 15, 18, 10, 19, color.RGBA{230, 50, 10, 255})
		img.Set(6, 16, color.RGBA{200, 20, 0, 180})

	case ProjMagicBall:
		fillCircle(img, 16, 16, 6, color.RGBA{120, 40, 180, 255})
		fillCircle(img, 16, 16, 4, color.RGBA{180, 80, 240, 255})
		fillCircle(img, 16, 16, 2, color.RGBA{230, 180, 255, 255})
		fillCircle(img, 16, 16, 8, color.RGBA{180, 80, 240, 100})

	case ProjHealBeam:
		fillRect(img, 15, 10, 2, 12, color.RGBA{50, 230, 100, 255})
		fillRect(img, 10, 15, 12, 2, color.RGBA{50, 230, 100, 255})
		fillRect(img, 15, 13, 2, 6, color.RGBA{220, 255, 230, 255})
		fillRect(img, 13, 15, 6, 2, color.RGBA{220, 255, 230, 255})
		img.Set(11, 11, color.RGBA{200, 255, 200, 200})
		img.Set(21, 11, color.RGBA{200, 255, 200, 200})
		img.Set(11, 21, color.RGBA{200, 255, 200, 200})
		img.Set(21, 21, color.RGBA{200, 255, 200, 200})
	}

	applyOutline(img, color.RGBA{10, 10, 15, 255})
	return ebiten.NewImageFromImage(img)
}

// --- UNIT AND ENEMY WALKING ANIMATION GENERATORS ---

func makeUnitSprite(ut UnitType, dir int, frame int) *ebiten.Image {
	img := image.NewRGBA(image.Rect(0, 0, 32, 32))

	bobY := int(math.Sin(float64(frame)*math.Pi/4.0) * 1.2)
	swing := math.Sin(float64(frame) * math.Pi / 4.0)
	legOffset := int(swing * 2.0)

	// Ground shadow
	fillEllipse(img, 16, 26, 7, 3, color.RGBA{0, 0, 0, 80})

	switch ut {
	case UnitKnight:
		bootCol := color.RGBA{100, 110, 120, 255}
		if dir == 2 || dir == 3 {
			fillRect(img, 14, 22+bobY, 2, 4+legOffset, bootCol)
			fillRect(img, 17, 22+bobY, 2, 4-legOffset, bootCol)
		} else {
			fillRect(img, 13, 22+bobY, 2, 4+legOffset, bootCol)
			fillRect(img, 18, 22+bobY, 2, 4-legOffset, bootCol)
		}

		chestCol := color.RGBA{150, 160, 170, 255}
		if dir == 1 {
			fillRect(img, 11, 14+bobY, 10, 9, color.RGBA{180, 30, 30, 255})
			fillRect(img, 12, 14+bobY, 8, 8, chestCol)
		} else if dir == 2 {
			fillRect(img, 13, 14+bobY, 6, 8, chestCol)
			fillRect(img, 18, 15+bobY, 2, 7, color.RGBA{180, 30, 30, 255})
		} else if dir == 3 {
			fillRect(img, 13, 14+bobY, 6, 8, chestCol)
			fillRect(img, 12, 15+bobY, 2, 7, color.RGBA{180, 30, 30, 255})
		} else {
			fillRect(img, 12, 14+bobY, 8, 8, chestCol)
			fillRect(img, 16, 14+bobY, 4, 8, color.RGBA{110, 120, 130, 255})
			fillRect(img, 15, 16+bobY, 2, 4, color.RGBA{220, 180, 40, 255})
			fillRect(img, 14, 17+bobY, 4, 2, color.RGBA{220, 180, 40, 255})
		}

		fillCircle(img, 16, 9+bobY, 4, chestCol)
		if dir == 0 {
			fillRect(img, 14, 9+bobY, 4, 1, color.RGBA{20, 20, 30, 255})
		} else if dir == 2 {
			fillRect(img, 12, 9+bobY, 3, 1, color.RGBA{20, 20, 30, 255})
		} else if dir == 3 {
			fillRect(img, 17, 9+bobY, 3, 1, color.RGBA{20, 20, 30, 255})
		}
		drawLine(img, 16, 5+bobY, 13, 3+bobY, color.RGBA{200, 40, 40, 255})

		shieldCol := color.RGBA{200, 40, 40, 255}
		goldCol := color.RGBA{220, 180, 40, 255}
		steelCol := color.RGBA{190, 195, 200, 255}

		if dir == 0 {
			fillRect(img, 8, 14+bobY+legOffset, 3, 6, shieldCol)
			drawLine(img, 8, 14+bobY+legOffset, 10, 14+bobY+legOffset, goldCol)
			drawLine(img, 22, 17+bobY-legOffset, 25, 10+bobY-legOffset, steelCol)
			drawLine(img, 21, 18+bobY-legOffset, 23, 16+bobY-legOffset, goldCol)
		} else if dir == 1 {
			fillRect(img, 21, 14+bobY, 3, 6, shieldCol)
			drawLine(img, 10, 17+bobY, 10, 9+bobY, steelCol)
		} else if dir == 2 {
			fillRect(img, 9, 14+bobY+legOffset, 3, 6, shieldCol)
			drawLine(img, 9, 14+bobY+legOffset, 11, 14+bobY+legOffset, goldCol)
			drawLine(img, 20, 17+bobY, 23, 11+bobY, steelCol)
		} else {
			drawLine(img, 21, 17+bobY-legOffset, 26, 11+bobY-legOffset, steelCol)
			drawLine(img, 20, 18+bobY-legOffset, 22, 16+bobY-legOffset, goldCol)
			fillRect(img, 9, 14+bobY+legOffset, 3, 6, shieldCol)
		}

	case UnitArcher:
		tunicCol := color.RGBA{40, 140, 60, 255}
		hairCol := color.RGBA{240, 210, 80, 255}
		skinCol := color.RGBA{245, 200, 170, 255}
		bowCol := color.RGBA{130, 90, 50, 255}
		bootCol := color.RGBA{110, 70, 30, 255}

		fillRect(img, 13, 22+bobY, 2, 4+legOffset, bootCol)
		fillRect(img, 17, 22+bobY, 2, 4-legOffset, bootCol)

		fillRect(img, 12, 14+bobY, 8, 8, tunicCol)
		if dir == 0 {
			fillRect(img, 12, 18+bobY, 8, 1, bootCol)
			img.Set(16, 18+bobY, color.RGBA{220, 180, 40, 255})
		} else if dir == 1 {
			fillRect(img, 12, 18+bobY, 8, 1, bootCol)
		} else if dir == 2 {
			fillRect(img, 13, 14+bobY, 6, 8, tunicCol)
			fillRect(img, 15, 18+bobY, 4, 1, bootCol)
		} else {
			fillRect(img, 13, 14+bobY, 6, 8, tunicCol)
			fillRect(img, 13, 18+bobY, 4, 1, bootCol)
		}

		fillCircle(img, 16, 9+bobY, 4, skinCol)
		fillCircle(img, 16, 8+bobY, 4, hairCol)
		fillRect(img, 12, 7+bobY, 8, 2, hairCol)
		if dir == 1 {
			fillRect(img, 12, 7+bobY, 8, 6, hairCol)
		} else if dir == 2 {
			fillRect(img, 16, 7+bobY, 4, 5, hairCol)
		} else if dir == 3 {
			fillRect(img, 12, 7+bobY, 4, 5, hairCol)
		} else {
			img.Set(14, 10+bobY, color.RGBA{30, 30, 50, 255})
			img.Set(18, 10+bobY, color.RGBA{30, 30, 50, 255})
		}

		if dir == 0 {
			drawLine(img, 9, 11+bobY, 9, 21+bobY, bowCol)
			drawLine(img, 10, 12+bobY, 10, 20+bobY, color.RGBA{240, 240, 240, 180})
		} else if dir == 1 {
			drawLine(img, 13, 13+bobY, 19, 19+bobY, bowCol)
		} else if dir == 2 {
			drawLine(img, 8, 11+bobY, 8, 21+bobY, bowCol)
			drawLine(img, 8, 11+bobY, 14, 16+bobY, color.RGBA{240, 240, 240, 180})
			drawLine(img, 8, 21+bobY, 14, 16+bobY, color.RGBA{240, 240, 240, 180})
			drawLine(img, 7, 16+bobY, 15, 16+bobY, color.RGBA{200, 170, 40, 255})
		} else {
			drawLine(img, 23, 11+bobY, 23, 21+bobY, bowCol)
			drawLine(img, 23, 11+bobY, 17, 16+bobY, color.RGBA{240, 240, 240, 180})
			drawLine(img, 23, 21+bobY, 17, 16+bobY, color.RGBA{240, 240, 240, 180})
			drawLine(img, 16, 16+bobY, 24, 16+bobY, color.RGBA{200, 170, 40, 255})
		}

	case UnitMage:
		robeCol := color.RGBA{110, 50, 150, 255}
		goldCol := color.RGBA{240, 190, 50, 255}
		beardCol := color.RGBA{215, 215, 220, 255}
		skinCol := color.RGBA{245, 200, 170, 255}
		staffCol := color.RGBA{100, 65, 30, 255}

		fillRect(img, 14, 24+bobY, 1, 2, color.RGBA{40, 30, 50, 255})
		fillRect(img, 17, 24+bobY, 1, 2, color.RGBA{40, 30, 50, 255})

		if dir == 2 || dir == 3 {
			fillRect(img, 13, 14+bobY, 6, 11, robeCol)
			fillRect(img, 13, 24+bobY, 6, 1, goldCol)
		} else {
			fillRect(img, 12, 14+bobY, 8, 11, robeCol)
			fillRect(img, 12, 24+bobY, 8, 1, goldCol)
			if dir == 0 {
				fillRect(img, 15, 14+bobY, 2, 11, goldCol)
			}
		}

		fillCircle(img, 16, 10+bobY, 3, skinCol)
		if dir == 1 {
			fillCircle(img, 16, 10+bobY, 3, robeCol)
		} else {
			fillRect(img, 14, 12+bobY, 4, 3, beardCol)
			img.Set(15, 15+bobY, beardCol)
			img.Set(16, 15+bobY, beardCol)
			if dir == 0 {
				img.Set(14, 11+bobY, color.RGBA{30, 30, 50, 255})
				img.Set(18, 11+bobY, color.RGBA{30, 30, 50, 255})
			}
		}
		for h := 0; h < 6; h++ {
			width := 4 - h/2
			if width < 1 {
				width = 1
			}
			fillRect(img, 16-width, 7+bobY-h, width*2, 1, robeCol)
			if h == 0 {
				fillRect(img, 11, 7+bobY, 10, 1, goldCol)
			}
		}

		pulse := math.Sin(float64(frame)*math.Pi/4.0)*30 + 225
		gCrystal := color.RGBA{80, uint8(pulse), 255, 255}

		if dir == 2 {
			drawLine(img, 10, 10+bobY, 10, 24+bobY, staffCol)
			fillCircle(img, 10, 8+bobY, 2, gCrystal)
		} else {
			drawLine(img, 22, 10+bobY, 22, 24+bobY, staffCol)
			fillCircle(img, 22, 8+bobY, 2, gCrystal)
		}

	case UnitHealer:
		robeCol := color.RGBA{240, 240, 245, 255}
		goldCol := color.RGBA{240, 180, 30, 255}
		hairCol := color.RGBA{240, 120, 160, 255}
		skinCol := color.RGBA{245, 200, 170, 255}
		staffCol := color.RGBA{210, 180, 100, 255}

		fillRect(img, 14, 23+bobY, 1, 3, color.RGBA{200, 100, 120, 255})
		fillRect(img, 17, 23+bobY, 1, 3, color.RGBA{200, 100, 120, 255})

		fillRect(img, 12, 14+bobY, 8, 10, robeCol)
		if dir == 0 {
			fillRect(img, 15, 16+bobY, 2, 5, color.RGBA{220, 40, 40, 255})
			fillRect(img, 13, 17+bobY, 6, 2, color.RGBA{220, 40, 40, 255})
		}
		fillRect(img, 12, 23+bobY, 8, 1, goldCol)

		fillCircle(img, 16, 9+bobY, 4, skinCol)
		fillCircle(img, 16, 8+bobY, 4, hairCol)
		fillRect(img, 12, 6+bobY, 8, 3, hairCol)
		if dir == 1 {
			fillRect(img, 12, 6+bobY, 8, 6, hairCol)
		} else if dir == 2 {
			fillRect(img, 15, 6+bobY, 4, 6, hairCol)
		} else if dir == 3 {
			fillRect(img, 13, 6+bobY, 4, 6, hairCol)
		} else {
			img.Set(14, 10+bobY, color.RGBA{50, 50, 120, 255})
			img.Set(18, 10+bobY, color.RGBA{50, 50, 120, 255})
		}

		sX := 22
		if dir == 2 {
			sX = 10
		}
		drawLine(img, sX, 10+bobY, sX, 24+bobY, staffCol)
		fillRect(img, sX-1, 6+bobY, 3, 3, goldCol)
		fillRect(img, sX-2, 7+bobY, 5, 1, goldCol)
		img.Set(sX, 7+bobY, color.RGBA{255, 255, 180, 255})
	}

	applyOutline(img, color.RGBA{15, 15, 25, 255})
	return ebiten.NewImageFromImage(img)
}

func makeEnemySprite(et EnemyType, dir int, frame int) *ebiten.Image {
	var img *image.RGBA
	if et == EnemyDragonBoss {
		img = image.NewRGBA(image.Rect(0, 0, 64, 64))
	} else {
		img = image.NewRGBA(image.Rect(0, 0, 32, 32))
	}

	bobY := int(math.Sin(float64(frame)*math.Pi/4.0) * 1.2)
	swing := math.Sin(float64(frame) * math.Pi / 4.0)
	legOffset := int(swing * 2.0)

	if et == EnemyDragonBoss {
		// Massive red dragon boss shadow
		fillEllipse(img, 32, 54, 18, 6, color.RGBA{0, 0, 0, 90})

		dragonRed := color.RGBA{180, 30, 30, 255}
		dragonDark := color.RGBA{100, 10, 10, 255}
		dragonGold := color.RGBA{230, 160, 40, 255}
		wingCol := color.RGBA{130, 20, 20, 255}

		// Flapping wings
		flapY := int(math.Sin(float64(frame)*math.Pi/4.0) * 4.0)
		if dir == 1 {
			fillEllipse(img, 14, 25+flapY, 12, 16, wingCol)
			fillEllipse(img, 50, 25+flapY, 12, 16, wingCol)
		} else if dir == 2 {
			fillEllipse(img, 45, 25+flapY, 8, 14, wingCol)
		} else if dir == 3 {
			fillEllipse(img, 19, 25+flapY, 8, 14, wingCol)
		} else {
			fillEllipse(img, 12, 28+flapY, 12, 18, wingCol)
			fillEllipse(img, 52, 28+flapY, 12, 18, wingCol)
		}

		// Legs
		fillRect(img, 24, 40+bobY, 5, 10+legOffset, dragonDark)
		fillRect(img, 35, 40+bobY, 5, 10-legOffset, dragonDark)

		// Body
		fillCircle(img, 32, 32+bobY, 14, dragonRed)
		fillEllipse(img, 32, 36+bobY, 11, 13, dragonRed)
		if dir == 0 {
			fillEllipse(img, 32, 35+bobY, 6, 9, dragonGold)
		}

		// Tail
		tailOffset := int(math.Sin(float64(frame)*math.Pi/4.0) * 5.0)
		if dir == 2 {
			drawLine(img, 42, 38+bobY, 56, 42+bobY+tailOffset, dragonDark)
			fillCircle(img, 56, 42+bobY+tailOffset, 3, dragonRed)
		} else if dir == 3 {
			drawLine(img, 22, 38+bobY, 8, 42+bobY+tailOffset, dragonDark)
			fillCircle(img, 8, 42+bobY+tailOffset, 3, dragonRed)
		} else {
			drawLine(img, 32, 45+bobY, 32+tailOffset, 55+bobY, dragonDark)
			fillCircle(img, 32+tailOffset, 55+bobY, 3, dragonRed)
		}

		// Head
		headY := 18 + bobY
		fillCircle(img, 32, headY, 7, dragonRed)
		if dir == 2 {
			fillRect(img, 24, headY-2, 8, 5, dragonRed)
			img.Set(25, headY-1, color.RGBA{255, 230, 100, 255})
			drawLine(img, 34, headY-5, 40, headY-11, color.RGBA{230, 220, 200, 255})
		} else if dir == 3 {
			fillRect(img, 32, headY-2, 8, 5, dragonRed)
			img.Set(38, headY-1, color.RGBA{255, 230, 100, 255})
			drawLine(img, 30, headY-5, 24, headY-11, color.RGBA{230, 220, 200, 255})
		} else {
			if dir == 0 {
				fillRect(img, 29, headY, 6, 6, dragonRed)
				img.Set(29, headY-1, color.RGBA{255, 230, 100, 255})
				img.Set(35, headY-1, color.RGBA{255, 230, 100, 255})
				drawLine(img, 29, headY-5, 24, headY-11, color.RGBA{230, 220, 200, 255})
				drawLine(img, 35, headY-5, 40, headY-11, color.RGBA{230, 220, 200, 255})
				if frame%2 == 0 {
					img.Set(31, headY+6, color.RGBA{255, 120, 30, 255})
					img.Set(32, headY+7, color.RGBA{255, 210, 40, 255})
					img.Set(33, headY+6, color.RGBA{255, 120, 30, 255})
				}
			} else {
				drawLine(img, 29, headY-5, 24, headY-11, color.RGBA{230, 220, 200, 255})
				drawLine(img, 35, headY-5, 40, headY-11, color.RGBA{230, 220, 200, 255})
			}
		}

		applyOutline(img, color.RGBA{15, 10, 10, 255})
		return ebiten.NewImageFromImage(img)
	}

	// 32x32 Ground Shadow
	fillEllipse(img, 16, 26, 7, 3, color.RGBA{0, 0, 0, 80})

	switch et {
	case EnemyGoblin:
		skinCol := color.RGBA{80, 160, 70, 255}
		clothCol := color.RGBA{110, 70, 30, 255}

		fillRect(img, 14, 21+bobY, 1, 5+legOffset, skinCol)
		fillRect(img, 17, 21+bobY, 1, 5-legOffset, skinCol)

		fillRect(img, 13, 15+bobY, 6, 7, skinCol)
		fillRect(img, 13, 19+bobY, 6, 3, clothCol)

		fillCircle(img, 16, 11+bobY, 3, skinCol)
		if dir == 0 {
			img.Set(15, 11+bobY, color.RGBA{255, 230, 0, 255})
			img.Set(17, 11+bobY, color.RGBA{255, 230, 0, 255})
			img.Set(12, 10+bobY, skinCol)
			img.Set(11, 10+bobY, skinCol)
			img.Set(20, 10+bobY, skinCol)
			img.Set(21, 10+bobY, skinCol)
		} else if dir == 2 {
			img.Set(14, 11+bobY, color.RGBA{255, 230, 0, 255})
			img.Set(18, 10+bobY, skinCol)
			img.Set(19, 10+bobY, skinCol)
		} else if dir == 3 {
			img.Set(18, 11+bobY, color.RGBA{255, 230, 0, 255})
			img.Set(14, 10+bobY, skinCol)
			img.Set(13, 10+bobY, skinCol)
		}

		dagCol := color.RGBA{130, 130, 140, 255}
		if dir == 2 {
			drawLine(img, 10, 17+bobY, 7, 17+bobY, dagCol)
		} else {
			drawLine(img, 21, 17+bobY, 24, 17+bobY, dagCol)
		}

	case EnemyOrc:
		skinCol := color.RGBA{60, 110, 50, 255}
		armorCol := color.RGBA{160, 40, 40, 255}
		axeCol := color.RGBA{120, 120, 125, 255}

		fillRect(img, 13, 21+bobY, 2, 5+legOffset, skinCol)
		fillRect(img, 17, 21+bobY, 2, 5-legOffset, skinCol)

		fillRect(img, 11, 13+bobY, 10, 9, armorCol)
		if dir == 0 {
			fillRect(img, 13, 13+bobY, 6, 9, armorCol)
			img.Set(11, 12+bobY, color.RGBA{240, 240, 230, 255})
			img.Set(20, 12+bobY, color.RGBA{240, 240, 230, 255})
		}

		fillCircle(img, 16, 9+bobY, 4, skinCol)
		if dir == 0 {
			img.Set(14, 9+bobY, color.RGBA{20, 20, 20, 255})
			img.Set(18, 9+bobY, color.RGBA{20, 20, 20, 255})
			img.Set(15, 11+bobY, color.RGBA{255, 255, 255, 255})
			img.Set(17, 11+bobY, color.RGBA{255, 255, 255, 255})
		}

		if dir == 2 {
			drawLine(img, 10, 10+bobY, 10, 22+bobY, color.RGBA{100, 65, 30, 255})
			fillRect(img, 7, 8+bobY, 4, 3, axeCol)
		} else {
			drawLine(img, 21, 10+bobY, 21, 22+bobY, color.RGBA{100, 65, 30, 255})
			fillRect(img, 21, 8+bobY, 4, 3, axeCol)
		}

	case EnemyTroll:
		skinCol := color.RGBA{70, 110, 140, 255}
		mossCol := color.RGBA{50, 130, 60, 255}
		malletCol := color.RGBA{110, 110, 115, 255}

		fillRect(img, 12, 22+bobY, 3, 4+legOffset, skinCol)
		fillRect(img, 17, 22+bobY, 3, 4-legOffset, skinCol)

		fillRect(img, 10, 11+bobY, 12, 12, skinCol)
		if dir == 1 {
			fillRect(img, 10, 11+bobY, 12, 8, mossCol)
		} else if dir == 2 {
			fillRect(img, 18, 11+bobY, 4, 10, mossCol)
		} else if dir == 3 {
			fillRect(img, 10, 11+bobY, 4, 10, mossCol)
		}

		fillCircle(img, 16, 9+bobY, 4, skinCol)
		if dir == 0 {
			img.Set(14, 9+bobY, color.RGBA{255, 100, 30, 255})
			img.Set(18, 9+bobY, color.RGBA{255, 100, 30, 255})
		}

		if dir == 2 {
			drawLine(img, 8, 8+bobY, 8, 22+bobY, color.RGBA{100, 65, 30, 255})
			fillRect(img, 5, 6+bobY, 7, 5, malletCol)
		} else {
			drawLine(img, 23, 8+bobY, 23, 22+bobY, color.RGBA{100, 65, 30, 255})
			fillRect(img, 20, 6+bobY, 7, 5, malletCol)
		}

	case EnemyBat:
		batCol := color.RGBA{45, 35, 55, 255}
		wingCol := color.RGBA{30, 20, 35, 255}
		eyeCol := color.RGBA{220, 30, 30, 255}

		// Higher floating shadow
		img.Set(15, 26, color.RGBA{0, 0, 0, 45})
		img.Set(16, 26, color.RGBA{0, 0, 0, 45})
		img.Set(17, 26, color.RGBA{0, 0, 0, 45})

		fSpan := int(math.Sin(float64(frame)*math.Pi/4.0) * 3)

		fillEllipse(img, 9, 11+bobY-fSpan, 6, 8, wingCol)
		fillEllipse(img, 23, 11+bobY-fSpan, 6, 8, wingCol)

		fillCircle(img, 16, 12+bobY, 4, batCol)
		if dir != 1 {
			img.Set(14, 11+bobY, eyeCol)
			img.Set(18, 11+bobY, eyeCol)
		}

	case EnemyFireElem:
		coreCol := color.RGBA{35, 20, 45, 255}
		fillEllipse(img, 16, 26, 5, 2, color.RGBA{0, 0, 0, 50})

		for i := 0; i < 20; i++ {
			fX := 16 + int(math.Sin(float64(i+frame)*1.5)*6.0)
			fY := 21 - i + int(math.Sin(float64(i+frame*2))*1.2)
			r := 6 - i/3
			if r < 1 {
				r = 1
			}
			var col color.RGBA
			if i < 6 {
				col = color.RGBA{210, 30, 10, 255}
			} else if i < 13 {
				col = color.RGBA{255, 120, 20, 255}
			} else {
				col = color.RGBA{255, 230, 80, 255}
			}
			fillCircle(img, fX, fY, r, col)
		}
		fillCircle(img, 16, 13+bobY, 3, coreCol)
		fillCircle(img, 16, 13+bobY, 1, color.RGBA{255, 100, 35, 255})

	case EnemyAssassin:
		cloakCol := color.RGBA{25, 25, 35, 255}
		eyeCol := color.RGBA{230, 30, 30, 255}
		dagCol := color.RGBA{180, 185, 195, 255}

		fillRect(img, 13, 22+bobY, 2, 4+legOffset, cloakCol)
		fillRect(img, 17, 22+bobY, 2, 4-legOffset, cloakCol)

		fillRect(img, 12, 13+bobY, 8, 10, cloakCol)

		fillCircle(img, 16, 9+bobY, 4, cloakCol)
		if dir == 0 {
			fillRect(img, 14, 8+bobY, 4, 3, color.RGBA{10, 10, 15, 255})
			img.Set(14, 9+bobY, eyeCol)
			img.Set(17, 9+bobY, eyeCol)
		} else if dir == 2 {
			fillRect(img, 12, 8+bobY, 3, 3, color.RGBA{10, 10, 15, 255})
			img.Set(13, 9+bobY, eyeCol)
		} else if dir == 3 {
			fillRect(img, 17, 8+bobY, 3, 3, color.RGBA{10, 10, 15, 255})
			img.Set(18, 9+bobY, eyeCol)
		}

		if dir == 2 {
			drawLine(img, 9, 16+bobY, 6, 16+bobY, dagCol)
			drawLine(img, 11, 18+bobY, 11, 21+bobY, dagCol)
		} else {
			drawLine(img, 22, 16+bobY, 25, 16+bobY, dagCol)
			drawLine(img, 20, 18+bobY, 20, 21+bobY, dagCol)
		}

	case EnemyNecro:
		robeCol := color.RGBA{30, 30, 40, 255}
		skullCol := color.RGBA{240, 240, 225, 255}
		greenGlow := color.RGBA{50, 220, 100, 255}
		staffCol := color.RGBA{80, 50, 30, 255}

		fillRect(img, 14, 24+bobY, 1, 2, color.RGBA{15, 15, 20, 255})
		fillRect(img, 17, 24+bobY, 1, 2, color.RGBA{15, 15, 20, 255})

		fillRect(img, 12, 13+bobY, 8, 12, robeCol)
		if dir == 0 {
			fillRect(img, 15, 13+bobY, 2, 12, greenGlow)
		}

		fillCircle(img, 16, 9+bobY, 4, robeCol)
		if dir != 1 {
			fillRect(img, 14, 8+bobY, 4, 4, skullCol)
			img.Set(14, 9+bobY, color.RGBA{20, 20, 20, 255})
			img.Set(17, 9+bobY, color.RGBA{20, 20, 20, 255})
			img.Set(15, 10+bobY, color.RGBA{20, 20, 20, 255})
			img.Set(16, 10+bobY, color.RGBA{20, 20, 20, 255})
		}

		sX := 22
		if dir == 2 {
			sX = 10
		}
		drawLine(img, sX, 10+bobY, sX, 24+bobY, staffCol)
		fillCircle(img, sX, 8+bobY, 2, greenGlow)

	case EnemySkeleton:
		boneCol := color.RGBA{235, 230, 210, 255}
		helmetCol := color.RGBA{120, 120, 125, 255}
		eyeCol := color.RGBA{220, 30, 30, 255}

		drawLine(img, 14, 21+bobY, 14, 25+bobY+legOffset, boneCol)
		drawLine(img, 17, 21+bobY, 17, 25+bobY-legOffset, boneCol)

		drawLine(img, 16, 13+bobY, 16, 21+bobY, boneCol)
		fillRect(img, 13, 14+bobY, 7, 1, boneCol)
		fillRect(img, 13, 17+bobY, 7, 1, boneCol)
		fillRect(img, 14, 19+bobY, 5, 1, boneCol)

		fillCircle(img, 16, 9+bobY, 3, boneCol)
		for x := 13; x <= 19; x++ {
			img.Set(x, 7+bobY, helmetCol)
		}
		img.Set(16, 6+bobY, helmetCol)
		if dir == 0 {
			img.Set(14, 9+bobY, eyeCol)
			img.Set(17, 9+bobY, eyeCol)
		}

		shieldCol := color.RGBA{130, 60, 30, 255}
		if dir == 0 {
			fillCircle(img, 10, 17+bobY, 3, shieldCol)
			img.Set(10, 17+bobY, color.RGBA{100, 100, 105, 255})
			drawLine(img, 22, 17+bobY, 25, 11+bobY, helmetCol)
		} else if dir == 2 {
			fillCircle(img, 10, 17+bobY, 3, shieldCol)
		} else {
			drawLine(img, 22, 17+bobY, 26, 12+bobY, helmetCol)
		}
	}

	applyOutline(img, color.RGBA{15, 15, 20, 255})
	return ebiten.NewImageFromImage(img)
}

// --- PRIMITIVE DRAWING UTILITIES ---

func fillRect(img *image.RGBA, x, y, w, h int, col color.RGBA) {
	for dy := 0; dy < h; dy++ {
		for dx := 0; dx < w; dx++ {
			img.Set(x+dx, y+dy, col)
		}
	}
}

func fillCircle(img *image.RGBA, cx, cy, r int, col color.RGBA) {
	for y := cy - r; y <= cy+r; y++ {
		for x := cx - r; x <= cx+r; x++ {
			if (x-cx)*(x-cx)+(y-cy)*(y-cy) <= r*r {
				img.Set(x, y, col)
			}
		}
	}
}

func fillEllipse(img *image.RGBA, cx, cy, rx, ry int, col color.RGBA) {
	frx2 := float64(rx * rx)
	fry2 := float64(ry * ry)
	for y := cy - ry; y <= cy+ry; y++ {
		for x := cx - rx; x <= cx+rx; x++ {
			dx := float64(x - cx)
			dy := float64(y - cy)
			if (dx*dx)/frx2+(dy*dy)/fry2 <= 1.0 {
				img.Set(x, y, col)
			}
		}
	}
}

func drawLine(img *image.RGBA, x1, y1, x2, y2 int, col color.RGBA) {
	dx := x2 - x1
	dy := y2 - y1
	steps := int(math.Max(math.Abs(float64(dx)), math.Abs(float64(dy))))
	if steps == 0 {
		img.Set(x1, y1, col)
		return
	}
	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)
		x := int(float64(x1) + t*float64(dx) + 0.5)
		y := int(float64(y1) + t*float64(dy) + 0.5)
		img.Set(x, y, col)
	}
}

// applyOutline traces the perimeter of transparent-border pixels adjacent to non-transparent pixels.
func applyOutline(img *image.RGBA, outlineColor color.RGBA) {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	alpha := make([]byte, w*h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a > 0 {
				alpha[y*w+x] = 255
			}
		}
	}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if alpha[y*w+x] == 0 {
				hasNeighbor := false
				for dy := -1; dy <= 1; dy++ {
					for dx := -1; dx <= 1; dx++ {
						if dx == 0 && dy == 0 {
							continue
						}
						nx, ny := x+dx, y+dy
						if nx >= 0 && nx < w && ny >= 0 && ny < h {
							if alpha[ny*w+nx] > 0 {
								hasNeighbor = true
								break
							}
						}
					}
					if hasNeighbor {
						break
					}
				}
				if hasNeighbor {
					img.Set(x, y, outlineColor)
				}
			}
		}
	}
}
