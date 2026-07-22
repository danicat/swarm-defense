package main

import (
	"fmt"
	"image/color"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	colorBg        = color.RGBA{0x10, 0x10, 0x18, 0xff}
	colorText      = color.RGBA{0xe0, 0xe0, 0xe0, 0xff}
	colorHighlight = color.RGBA{0x00, 0xff, 0xff, 0xff}
	colorRed       = color.RGBA{0xff, 0x40, 0x40, 0xff}
	colorGreen     = color.RGBA{0x40, 0xff, 0x40, 0xff}
	colorGold      = color.RGBA{0xff, 0xd7, 0x00, 0xff}
	colorPanel     = color.RGBA{0x20, 0x20, 0x30, 0xff}
	colorGrid      = color.RGBA{0x40, 0x40, 0x50, 0x80}
)

var pixelFont = map[rune][]string{
	'A': {".###.", "#...#", "#####", "#...#", "#...#"},
	'B': {"####.", "#...#", "####.", "#...#", "####."},
	'C': {".###.", "#...#", "#....", "#...#", ".###."},
	'D': {"####.", "#...#", "#...#", "#...#", "####."},
	'E': {"#####", "#....", "####.", "#....", "#####"},
	'F': {"#####", "#....", "####.", "#....", "#...."},
	'G': {".###.", "#....", "#.###", "#...#", ".###."},
	'H': {"#...#", "#...#", "#####", "#...#", "#...#"},
	'I': {"#####", "..#..", "..#..", "..#..", "#####"},
	'J': {"..###", "...#.", "...#.", "#..#.", ".##.."},
	'K': {"#...#", "#..#.", "###..", "#..#.", "#...#"},
	'L': {"#....", "#....", "#....", "#....", "#####"},
	'M': {"#...#", "##.##", "#.#.#", "#...#", "#...#"},
	'N': {"#...#", "##..#", "#.#.#", "#..##", "#...#"},
	'O': {".###.", "#...#", "#...#", "#...#", ".###."},
	'P': {"####.", "#...#", "####.", "#....", "#...."},
	'Q': {".###.", "#...#", "#.#.#", "#..#.", ".##.#"},
	'R': {"####.", "#...#", "####.", "#..#.", "#...#"},
	'S': {".####", "#....", ".###.", "....#", "####."},
	'T': {"#####", "..#..", "..#..", "..#..", "..#.."},
	'U': {"#...#", "#...#", "#...#", "#...#", ".###."},
	'V': {"#...#", "#...#", "#...#", ".#.#.", "..#.."},
	'W': {"#...#", "#...#", "#.#.#", "##.##", "#...#"},
	'X': {"#...#", ".#.#.", "..#..", ".#.#.", "#...#"},
	'Y': {"#...#", ".#.#.", "..#..", "..#..", "..#.."},
	'Z': {"#####", "...#.", "..#..", ".#...", "#####"},
	'0': {".###.", "#..##", "#.#.#", "##..#", ".###."},
	'1': {"..#..", ".##..", "..#..", "..#..", ".###."},
	'2': {".###.", "#...#", "...#.", "..#..", "#####"},
	'3': {"#####", "...#.", ".###.", "...#.", "#####"},
	'4': {"#...#", "#...#", "#####", "...#.", "...#."},
	'5': {"#####", "#....", "####.", "....#", "####."},
	'6': {".###.", "#....", "####.", "#...#", ".###."},
	'7': {"#####", "....#", "...#.", "..#..", ".#..."},
	'8': {".###.", "#...#", ".###.", "#...#", ".###."},
	'9': {".###.", "#...#", ".####", "....#", ".###."},
	' ': {".....", ".....", ".....", ".....", "....."},
	'.': {".....", ".....", ".....", ".....", "..#.."},
	',': {".....", ".....", ".....", "..#..", ".#..."},
	'!': {"..#..", "..#..", "..#..", ".....", "..#.."},
	'?': {".###.", "#...#", "...#.", "..#..", "..#.."},
	':': {".....", "..#..", ".....", "..#..", "....."},
	'-': {".....", ".....", "#####", ".....", "....."},
	'+': {".....", "..#..", "#####", "..#..", "....."},
	'/': {"....#", "...#.", "..#..", ".#...", "#...."},
	'(': {"..#..", ".#...", ".#...", ".#...", "..#.."},
	')': {"..#..", "...#.", "...#.", "...#.", "..#.."},
}

// DrawText draws retro pixel text at (x, y). scale is size multiplier.
func DrawText(screen *ebiten.Image, text string, x, y, scale int, c color.Color) {
	text = strings.ToUpper(text)
	startX := float32(x)
	startY := float32(y)
	s := float32(scale)

	for _, char := range text {
		matrix, ok := pixelFont[char]
		if !ok {
			matrix = pixelFont[' ']
		}
		for row := 0; row < 5; row++ {
			for col := 0; col < 5; col++ {
				if matrix[row][col] == '#' {
					vector.DrawFilledRect(screen, startX+float32(col)*s, startY+float32(row)*s, s, s, c, true)
				}
			}
		}
		startX += 6 * s
	}
}

func MeasureText(text string, scale int) (width, height int) {
	return len(text) * 6 * scale, 5 * scale
}

// DrawTextCentered centers text horizontally around cx.
func DrawTextCentered(screen *ebiten.Image, text string, cx, y, scale int, c color.Color) {
	w, _ := MeasureText(text, scale)
	DrawText(screen, text, cx-w/2, y, scale, c)
}

func DrawUI(screen *ebiten.Image, g *Game) {
	switch g.state {
	case StateIntro:
		drawIntro(screen, g)
	case StateTitle:
		drawTitle(screen, g)
	case StateBuild:
		drawGameEntities(screen, g)
		drawBuildHUD(screen, g)
		drawGridOverlay(screen, g)
	case StateAttack:
		drawGameEntities(screen, g)
		drawAttackHUD(screen, g)
		drawGridOverlay(screen, g)
	case StateWin:
		drawWin(screen, g)
	case StateOver:
		drawGameOver(screen, g)
	}
}

func drawIntro(screen *ebiten.Image, g *Game) {
	screen.Fill(colorBg)
	elapsed := time.Since(g.introTimer).Seconds()
	scrollY := 480 - int(elapsed*40)

	DrawTextCentered(screen, "EPIC TOWER DEFENSE", 320, scrollY, 3, colorHighlight)
	DrawTextCentered(screen, "A SWARM OF ENEMIES APPROACHES", 320, scrollY+60, 2, colorText)
	DrawTextCentered(screen, "DEFEND YOUR BASE", 320, scrollY+100, 2, colorText)
	DrawTextCentered(screen, "BUILD TOWERS TO SURVIVE", 320, scrollY+140, 2, colorText)

	if scrollY < -200 {
		DrawTextCentered(screen, "PRESS ENTER TO SKIP", 320, 400, 2, colorHighlight)
	}
}

func drawTitle(screen *ebiten.Image, g *Game) {
	screen.Fill(colorBg)
	
	glow := uint8(128 + 127*time.Now().UnixNano()%1000000000/1000000000)
	titleColor := color.RGBA{0, glow, 255, 255}
	
	DrawTextCentered(screen, "SWARM DEFENSE", 320, 100, 4, titleColor)
	DrawTextCentered(screen, fmt.Sprintf("HIGH SCORE: %d", g.highScore), 320, 250, 2, colorGold)

	if g.titleTicks%60 < 30 {
		DrawTextCentered(screen, "PRESS ENTER TO START", 320, 350, 2, colorHighlight)
	}
}

func drawBuildHUD(screen *ebiten.Image, g *Game) {
	// Draw bottom UI panel
	vector.DrawFilledRect(screen, 0, 384, 640, 96, colorPanel, true)
	vector.StrokeRect(screen, 0, 384, 640, 96, 2, colorHighlight, true)

	// Panel title and stats
	DrawText(screen, "BUILD PHASE", 10, 394, 2, colorHighlight)
	DrawText(screen, fmt.Sprintf("GOLD: %d", g.gold), 10, 420, 2, colorGold)
	DrawText(screen, fmt.Sprintf("LIVES: %d", g.lives), 10, 440, 2, colorRed)
	DrawText(screen, fmt.Sprintf("WAVE: %d", g.wave), 10, 460, 2, colorText)

	catalogX := 160
	catalogY := 400
	itemSize := 32

	// Render the 8 Catalog Items
	for i := 0; i < 8; i++ {
		x := catalogX + i*(itemSize+12)
		y := catalogY

		// Draw slot background box
		vector.DrawFilledRect(screen, float32(x), float32(y), float32(itemSize), float32(itemSize), color.RGBA{30, 30, 45, 255}, true)
		vector.StrokeRect(screen, float32(x), float32(y), float32(itemSize), float32(itemSize), 1, color.RGBA{60, 60, 80, 255}, true)

		// Draw actual sprite inside the slot
		var spriteImg *ebiten.Image
		if i <= 3 {
			// Building type (0-3)
			spriteImg = g.sprites.Buildings[BuildingType(i)]
		} else {
			// Unit type (0-3)
			spriteImg = g.sprites.Units[UnitType(i-4)][0][0] // Down direction, frame 0
		}

		if spriteImg != nil {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x), float64(y))
			screen.DrawImage(spriteImg, op)
		}

		// Highlight if selected
		isSelected := false
		if g.isPlacingBuilding && i == g.selectedBuildType {
			isSelected = true
		} else if g.isPlacingUnit && (i-4) == g.selectedUnitType {
			isSelected = true
		}

		if isSelected {
			vector.StrokeRect(screen, float32(x-2), float32(y-2), float32(itemSize+4), float32(itemSize+4), 2, colorHighlight, true)
		}

		// Draw Item Name and Cost
		itemName := Catalog[i].Name
		itemCost := Catalog[i].Cost
		DrawTextCentered(screen, itemName, x+itemSize/2, y+itemSize+4, 1, colorText)
		DrawTextCentered(screen, fmt.Sprintf("%d", itemCost), x+itemSize/2, y+itemSize+13, 1, colorGold)
	}

	// Instructions on the right side
	DrawText(screen, "ENTER: START", 545, 400, 1, colorGreen)
	DrawText(screen, "SPACE: READY", 545, 415, 1, colorText)
	DrawText(screen, "CLICK HUD SLOT", 545, 435, 1, colorHighlight)
	DrawText(screen, "CLICK MAP GRD", 545, 450, 1, colorHighlight)
}

func drawAttackHUD(screen *ebiten.Image, g *Game) {
	vector.DrawFilledRect(screen, 0, 384, 640, 96, colorPanel, true)
	vector.StrokeRect(screen, 0, 384, 640, 96, 2, colorRed, true)

	DrawText(screen, "ATTACK PHASE", 10, 394, 2, colorRed)
	DrawText(screen, fmt.Sprintf("GOLD: %d", g.gold), 10, 420, 2, colorGold)
	DrawText(screen, fmt.Sprintf("LIVES: %d", g.lives), 10, 440, 2, colorRed)
	DrawText(screen, fmt.Sprintf("WAVE: %d", g.wave), 10, 460, 2, colorText)
	DrawText(screen, fmt.Sprintf("SCORE: %d", g.score), 200, 394, 2, colorText)

	DrawText(screen, "BASE HP", 200, 420, 2, colorText)
	hpPercent := float32(g.baseHealth) / 100.0
	if hpPercent < 0 { hpPercent = 0 }
	vector.DrawFilledRect(screen, 200, 440, 150, 15, color.RGBA{50, 0, 0, 255}, true)
	vector.DrawFilledRect(screen, 200, 440, 150*hpPercent, 15, colorRed, true)

	DrawText(screen, "WAVE PROG", 400, 420, 2, colorText)
	wavePercent := float32(0)
	if len(g.spawnQueue) > 0 || g.waveTimer > 0 {
		wavePercent = 1.0 - float32(len(g.spawnQueue))/20.0
	}
	if wavePercent < 0 { wavePercent = 0 }
	if wavePercent > 1 { wavePercent = 1 }
	vector.DrawFilledRect(screen, 400, 440, 150, 15, color.RGBA{0, 50, 0, 255}, true)
	vector.DrawFilledRect(screen, 400, 440, 150*wavePercent, 15, colorGreen, true)
}

func drawGridOverlay(screen *ebiten.Image, g *Game) {
	for x := 0; x <= 20; x++ {
		vector.StrokeLine(screen, float32(x*32), 0, float32(x*32), 384, 1, colorGrid, true)
	}
	for y := 0; y <= 12; y++ {
		vector.StrokeLine(screen, 0, float32(y*32), 640, float32(y*32), 1, colorGrid, true)
	}

	if g.state == StateBuild {
		cx, cy := float32(g.cursorX*32), float32(g.cursorY*32)
		vector.StrokeRect(screen, cx, cy, 32, 32, 2, colorHighlight, true)
		
		if g.isPlacingBuilding {
			ranges := []float32{3 * 32, 4 * 32, 2 * 32, 0}
			if g.selectedBuildType >= 0 && g.selectedBuildType < len(ranges) {
				r := ranges[g.selectedBuildType]
				if r > 0 {
					vector.StrokeCircle(screen, cx+16, cy+16, r, 1, color.RGBA{0, 255, 255, 100}, true)
				}
			}
		}
	}
}

func drawWin(screen *ebiten.Image, g *Game) {
	screen.Fill(colorBg)
	DrawTextCentered(screen, "VICTORY", 320, 100, 6, colorGold)
	DrawTextCentered(screen, fmt.Sprintf("FINAL SCORE: %d", g.score), 320, 200, 3, colorText)
	DrawTextCentered(screen, fmt.Sprintf("WAVE REACHED: %d", g.wave), 320, 250, 2, colorText)
	DrawTextCentered(screen, fmt.Sprintf("GOLD EARNED: %d", g.gold), 320, 300, 2, colorText)
	
	if g.winTicks%60 < 30 {
		DrawTextCentered(screen, "PRESS ENTER TO RETURN", 320, 400, 2, colorHighlight)
	}
}

func drawGameOver(screen *ebiten.Image, g *Game) {
	screen.Fill(color.RGBA{50, 0, 0, 255})
	DrawTextCentered(screen, "GAME OVER", 320, 100, 6, colorRed)
	DrawTextCentered(screen, "THE BASE WAS DESTROYED", 320, 180, 2, colorText)
	DrawTextCentered(screen, fmt.Sprintf("FINAL SCORE: %d", g.score), 320, 250, 3, colorText)
	DrawTextCentered(screen, fmt.Sprintf("WAVE REACHED: %d", g.wave), 320, 300, 2, colorText)
	
	if g.overTicks%60 < 30 {
		DrawTextCentered(screen, "PRESS ENTER TO RETRY", 320, 400, 2, colorHighlight)
	}
}

func drawGameEntities(screen *ebiten.Image, g *Game) {
	// 1. Draw Grid Tiles
	for x := 0; x < 20; x++ {
		for y := 0; y < 12; y++ {
			tileIdx := g.grid[x][y]
			var img *ebiten.Image
			if tileIdx == 3 {
				img = g.sprites.Tiles[3]
			} else if tileIdx >= 0 && tileIdx < 10 {
				img = g.sprites.Tiles[tileIdx]
			}
			if img != nil {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(x*32), float64(y*32))
				screen.DrawImage(img, op)
			}
		}
	}

	// 2. Draw Buildings
	for _, b := range g.buildings {
		img := g.sprites.Buildings[b.Type]
		if img != nil {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(b.GridX*32), float64(b.GridY*32))
			screen.DrawImage(img, op)
		}
		if b.Health < b.MaxHealth {
			drawHealthBar(screen, float64(b.GridX*32+4), float64(b.GridY*32-6), 24, 3, b.Health/b.MaxHealth)
		}
	}

	// 3. Draw Units
	for _, u := range g.units {
		frame := u.AnimFrame % 8
		img := g.sprites.Units[u.Type][u.Direction][frame]
		if img != nil {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(u.X-16, u.Y-16)
			screen.DrawImage(img, op)
		}
		if u.Health < u.MaxHealth {
			drawHealthBar(screen, u.X-12, u.Y-18, 24, 3, u.Health/u.MaxHealth)
		}
	}

	// 4. Draw Enemies
	for _, e := range g.enemies {
		frame := e.AnimFrame % 8
		img := g.sprites.Enemies[e.Type][e.Direction][frame]
		if img != nil {
			op := &ebiten.DrawImageOptions{}
			if e.Type == EnemyDragonBoss {
				op.GeoM.Translate(e.X-32, e.Y-32)
			} else {
				op.GeoM.Translate(e.X-16, e.Y-16)
			}
			screen.DrawImage(img, op)
		}
		if e.Health < e.MaxHealth {
			hY := e.Y - 18
			if e.Type == EnemyDragonBoss {
				hY = e.Y - 36
			}
			drawHealthBar(screen, e.X-12, hY, 24, 3, e.Health/e.MaxHealth)
		}
	}

	// 5. Draw Projectiles
	for _, p := range g.projectiles {
		img := g.sprites.Projectiles[p.Type]
		if img != nil {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(p.X-16, p.Y-16)
			screen.DrawImage(img, op)
		}
	}

	// 6. Draw Particles
	for _, part := range g.particles {
		c := color.RGBA{
			R: uint8((part.Color >> 24) & 0xFF),
			G: uint8((part.Color >> 16) & 0xFF),
			B: uint8((part.Color >> 8) & 0xFF),
			A: uint8(part.Color & 0xFF),
		}
		vector.DrawFilledRect(screen, float32(part.X-1), float32(part.Y-1), 2, 2, c, true)
	}
}

func drawHealthBar(screen *ebiten.Image, x, y, w, h float64, percent float64) {
	if percent < 0 { percent = 0 }
	if percent > 1 { percent = 1 }
	vector.DrawFilledRect(screen, float32(x), float32(y), float32(w), float32(h), color.RGBA{50, 0, 0, 255}, true)
	vector.DrawFilledRect(screen, float32(x), float32(y), float32(w*percent), float32(h), color.RGBA{40, 220, 40, 255}, true)
}
