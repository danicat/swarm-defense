package main

import (
	"math"
	"math/rand"
)

type UnitType int

const (
	UnitKnight UnitType = iota
	UnitArcher
	UnitMage
	UnitHealer
)

type Unit struct {
	ID             int
	Type           UnitType
	X, Y           float64 // Pixel coordinates
	GridX, GridY   int
	Health         float64
	MaxHealth      float64
	Speed          float64
	Damage         float64
	Range          float64
	AttackCooldown float64
	AnimFrame      int
	AnimTick       int
	Direction      int // 0: Down, 1: Up, 2: Left, 3: Right
	Target         *Enemy
}

type BuildingType int

const (
	BldArrow BuildingType = iota
	BldCannon
	BldIce
	BldGoldMine
)

type Building struct {
	ID         int
	Type       BuildingType
	GridX, GridY int
	X, Y       float64
	Health     float64
	MaxHealth  float64
	Range      float64
	Damage     float64
	ShootTimer float64
	MineTimer  float64
}

type EnemyType int

const (
	EnemyGoblin EnemyType = iota
	EnemyOrc
	EnemyTroll
	EnemyBat
	EnemyFireElem
	EnemyAssassin
	EnemyNecro
	EnemySkeleton
	EnemyDragonBoss
)

type Enemy struct {
	ID             int
	Type           EnemyType
	X, Y           float64 // Pixel coordinates
	Health         float64
	MaxHealth      float64
	Speed          float64
	Damage         float64 // Damage to player base
	GoldReward     int
	Path           [][2]int // Grid coordinates of the path
	PathIndex      int
	AnimFrame      int
	AnimTick       int
	Direction      int // 0: Down, 1: Up, 2: Left, 3: Right
	SlowDuration   float64
	NecroTimer     float64
	FireTrailTimer float64
	IsStealth      bool
}

type ProjectileType int

const (
	ProjArrow ProjectileType = iota
	ProjCannonBall
	ProjIceShard
	ProjFireBall
	ProjMagicBall
	ProjHealBeam
)

type Projectile struct {
	Type       ProjectileType
	X, Y       float64
	TargetX, TargetY float64
	Speed      float64
	Damage     float64
	IsSplash   bool
	IsSlow     bool
	Target     *Enemy
	TargetUnit *Unit // For boss fireballs or healers
}

type Particle struct {
	X, Y    float64
	VX, VY  float64
	Life    float64 // Lifetime in frames
	MaxLife float64
	Color   uint32
}

var entityIDCounter int

func nextID() int {
	entityIDCounter++
	return entityIDCounter
}

func NewUnit(t UnitType, gridX, gridY int) *Unit {
	u := &Unit{
		ID:    nextID(),
		Type:  t,
		GridX: gridX,
		GridY: gridY,
		X:     float64(gridX*32 + 16),
		Y:     float64(gridY*32 + 16),
	}
	switch t {
	case UnitKnight:
		u.MaxHealth, u.Damage, u.Speed, u.Range = 200, 20, 1.0, 40
	case UnitArcher:
		u.MaxHealth, u.Damage, u.Speed, u.Range = 100, 15, 2.0, 120
	case UnitMage:
		u.MaxHealth, u.Damage, u.Speed, u.Range = 80, 25, 1.5, 100
	case UnitHealer:
		u.MaxHealth, u.Damage, u.Speed, u.Range = 80, -10, 1.5, 120 // Negative damage for healing
	}
	u.Health = u.MaxHealth
	return u
}

func NewBuilding(t BuildingType, gridX, gridY int) *Building {
	b := &Building{
		ID:        nextID(),
		Type:      t,
		GridX:     gridX,
		GridY:     gridY,
		X:         float64(gridX*32 + 16),
		Y:         float64(gridY*32 + 16),
		MaxHealth: 300,
	}
	b.Health = b.MaxHealth
	switch t {
	case BldArrow:
		b.Damage, b.Range = 10, 150
	case BldCannon:
		b.Damage, b.Range = 30, 120
	case BldIce:
		b.Damage, b.Range = 5, 100
	case BldGoldMine:
		b.Damage, b.Range = 0, 0
	}
	return b
}

func NewEnemy(t EnemyType, startX, startY float64) *Enemy {
	e := &Enemy{
		ID:   nextID(),
		Type: t,
		X:    startX,
		Y:    startY,
	}
	switch t {
	case EnemyGoblin:
		e.MaxHealth, e.Speed, e.Damage, e.GoldReward = 50, 2.0, 5, 15
	case EnemyOrc:
		e.MaxHealth, e.Speed, e.Damage, e.GoldReward = 100, 1.0, 10, 25
	case EnemyTroll:
		e.MaxHealth, e.Speed, e.Damage, e.GoldReward = 300, 0.5, 20, 40
	case EnemyBat:
		e.MaxHealth, e.Speed, e.Damage, e.GoldReward = 40, 1.5, 5, 30
	case EnemyFireElem:
		e.MaxHealth, e.Speed, e.Damage, e.GoldReward = 80, 1.2, 10, 50
	case EnemyAssassin:
		e.MaxHealth, e.Speed, e.Damage, e.GoldReward = 60, 2.5, 15, 60
		e.IsStealth = true
	case EnemyNecro:
		e.MaxHealth, e.Speed, e.Damage, e.GoldReward = 120, 0.8, 15, 75
	case EnemySkeleton:
		e.MaxHealth, e.Speed, e.Damage, e.GoldReward = 30, 1.5, 5, 5
	case EnemyDragonBoss:
		e.MaxHealth, e.Speed, e.Damage, e.GoldReward = 1500, 0.6, 50, 500
	}
	e.Health = e.MaxHealth
	return e
}

func NewProjectile(t ProjectileType, x, y, tx, ty float64, target *Enemy, targetUnit *Unit, damage float64) *Projectile {
	p := &Projectile{
		Type:       t,
		X:          x,
		Y:          y,
		TargetX:    tx,
		TargetY:    ty,
		Target:     target,
		TargetUnit: targetUnit,
		Damage:     damage,
	}
	switch t {
	case ProjArrow:
		p.Speed = 5.0
	case ProjCannonBall:
		p.Speed, p.IsSplash = 3.0, true
	case ProjIceShard:
		p.Speed, p.IsSlow = 4.0, true
	case ProjFireBall:
		p.Speed, p.IsSplash = 4.0, true
	case ProjMagicBall:
		p.Speed, p.IsSplash = 4.0, true
	case ProjHealBeam:
		p.Speed = 6.0
	}
	return p
}

func NewParticle(x, y float64, color uint32) *Particle {
	angle := rand.Float64() * 2 * math.Pi
	speed := rand.Float64()*2 + 1
	return &Particle{
		X:       x,
		Y:       y,
		VX:      math.Cos(angle) * speed,
		VY:      math.Sin(angle) * speed,
		Life:    0,
		MaxLife: float64(20 + rand.Intn(20)),
		Color:   color,
	}
}

// Distance helper
func Distance(x1, y1, x2, y2 float64) float64 {
	dx, dy := x1-x2, y1-y2
	return math.Sqrt(dx*dx + dy*dy)
}

// A* Node
type node struct {
	x, y   int
	g, h, f int
	parent *node
}

// Simple BFS / Pathfinding on 20x12 grid (0=Empty, 1=Path, 2=Obstacle, 3=Bld, 4=Spawn, 5=Base)
func FindPath(grid *[20][12]int, startX, startY, endX, endY int) [][2]int {
	// Directions: up, right, down, left
	dx := []int{0, 1, 0, -1}
	dy := []int{-1, 0, 1, 0}

	visited := make([][]bool, 20)
	for i := range visited {
		visited[i] = make([]bool, 12)
	}

	queue := [][2]int{{startX, startY}}
	visited[startX][startY] = true
	parent := make(map[[2]int][2]int)

	found := false
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if curr[0] == endX && curr[1] == endY {
			found = true
			break
		}

		for i := 0; i < 4; i++ {
			nx, ny := curr[0]+dx[i], curr[1]+dy[i]
			if nx >= 0 && nx < 20 && ny >= 0 && ny < 12 && !visited[nx][ny] {
				// Can walk on Empty (0), Path (1), Spawn (4), Base (5)
				// Cannot walk on Obstacle (2) or Building (3)
				val := grid[nx][ny]
				if val == 0 || val == 1 || val == 4 || val == 5 {
					visited[nx][ny] = true
					queue = append(queue, [2]int{nx, ny})
					parent[[2]int{nx, ny}] = curr
				}
			}
		}
	}

	if !found {
		return nil
	}

	path := [][2]int{}
	curr := [2]int{endX, endY}
	for curr != [2]int{startX, startY} {
		path = append([][2]int{curr}, path...) // Prepend
		curr = parent[curr]
	}
	return path
}

// Update Units
func (u *Unit) Update(enemies []*Enemy, buildings []*Building, units []*Unit, spawnProjectile func(*Projectile)) {
	if u.AttackCooldown > 0 {
		u.AttackCooldown--
	}
	
	// Find target
	var target *Enemy
	minDist := u.Range
	for _, e := range enemies {
		d := Distance(u.X, u.Y, e.X, e.Y)
		
		// Stealth check for Assassin
		if e.IsStealth && d > 64 { // 2 tiles = 64 pixels
			continue 
		}
		
		if d <= minDist {
			minDist = d
			target = e
		}
	}
	u.Target = target

	if u.Target != nil && u.AttackCooldown <= 0 {
		if u.Type == UnitKnight || u.Type == UnitArcher || u.Type == UnitMage {
			// Attack
			var projType ProjectileType
			if u.Type == UnitArcher {
				projType = ProjArrow
			} else if u.Type == UnitMage {
				projType = ProjMagicBall
			}
			
			if u.Type == UnitKnight {
				u.Target.Health -= u.Damage
			} else {
				spawnProjectile(NewProjectile(projType, u.X, u.Y, u.Target.X, u.Target.Y, u.Target, nil, u.Damage))
			}
			u.AttackCooldown = 60 // 1 sec cooldown at 60 TPS
		}
	}
	
	if u.Type == UnitHealer && u.AttackCooldown <= 0 {
		// Heal closest damaged building or unit
		var healTargetUnit *Unit
		var healTargetBld *Building
		hMinDist := u.Range
		for _, ou := range units {
			if ou.Health < ou.MaxHealth {
				d := Distance(u.X, u.Y, ou.X, ou.Y)
				if d <= hMinDist {
					hMinDist = d
					healTargetUnit = ou
				}
			}
		}
		for _, b := range buildings {
			if b.Health < b.MaxHealth {
				d := Distance(u.X, u.Y, b.X, b.Y)
				if d <= hMinDist {
					hMinDist = d
					healTargetBld = b
					healTargetUnit = nil
				}
			}
		}
		if healTargetUnit != nil {
			spawnProjectile(NewProjectile(ProjHealBeam, u.X, u.Y, healTargetUnit.X, healTargetUnit.Y, nil, healTargetUnit, u.Damage))
			u.AttackCooldown = 60
		} else if healTargetBld != nil {
			// Instant heal building
			healTargetBld.Health -= u.Damage // negative damage heals
			if healTargetBld.Health > healTargetBld.MaxHealth {
				healTargetBld.Health = healTargetBld.MaxHealth
			}
			u.AttackCooldown = 60
		}
	}

	u.AnimTick++
	if u.AnimTick > 10 {
		u.AnimTick = 0
		u.AnimFrame = (u.AnimFrame + 1) % 4
	}
}

// Update Building
func (b *Building) Update(enemies []*Enemy, units []*Unit, addGold func(int), spawnProjectile func(*Projectile)) {
	if b.Type == BldGoldMine {
		b.MineTimer++
		if b.MineTimer >= 300 { // 5 seconds at 60 TPS
			addGold(10)
			b.MineTimer = 0
		}
		return
	}

	if b.ShootTimer > 0 {
		b.ShootTimer--
	}

	if b.ShootTimer <= 0 {
		var target *Enemy
		minDist := b.Range
		for _, e := range enemies {
			d := Distance(b.X, b.Y, e.X, e.Y)
			if e.IsStealth && d > 64 {
				// Towers can't see stealth outside 2 tiles
				// Actually rules say untargetable until within 2 tiles of a *Unit*, but we'll use distance to tower or unit logic
				continue
			}
			if d <= minDist {
				minDist = d
				target = e
			}
		}
		
		if target != nil {
			var pType ProjectileType
			if b.Type == BldArrow {
				pType = ProjArrow
				b.ShootTimer = 30 // fast
			} else if b.Type == BldCannon {
				pType = ProjCannonBall
				b.ShootTimer = 90 // slow
			} else if b.Type == BldIce {
				pType = ProjIceShard
				b.ShootTimer = 60
			}
			spawnProjectile(NewProjectile(pType, b.X, b.Y, target.X, target.Y, target, nil, b.Damage))
		}
	}
}

// Update Enemy
func (e *Enemy) Update(baseX, baseY float64, spawnEnemy func(*Enemy), spawnProjectile func(*Projectile), units []*Unit, spawnParticle func(*Particle)) bool {
	// Returns true if reached base
	if e.Type == EnemyTroll {
		// Regenerate 2 HP/s = 2/60 per tick
		e.Health += 2.0 / 60.0
		if e.Health > e.MaxHealth {
			e.Health = e.MaxHealth
		}
	}

	if e.Type == EnemyNecro {
		e.NecroTimer++
		if e.NecroTimer >= 240 { // 4 seconds
			skel := NewEnemy(EnemySkeleton, e.X, e.Y)
			skel.Path = e.Path // follow same path
			skel.PathIndex = e.PathIndex
			spawnEnemy(skel)
			e.NecroTimer = 0
		}
	}

	if e.Type == EnemyFireElem {
		e.FireTrailTimer++
		if e.FireTrailTimer >= 15 {
			spawnParticle(NewParticle(e.X, e.Y, 0xFF0000FF)) // Red particle
			e.FireTrailTimer = 0
		}
	}

	if e.Type == EnemyDragonBoss {
		// Randomly shoot fireball at units
		if rand.Float64() < 0.01 && len(units) > 0 {
			targetUnit := units[rand.Intn(len(units))]
			spawnProjectile(NewProjectile(ProjFireBall, e.X, e.Y, targetUnit.X, targetUnit.Y, nil, targetUnit, 30))
		}
	}

	speed := e.Speed
	if e.SlowDuration > 0 {
		speed *= 0.5
		e.SlowDuration--
	}

	// Move
	if e.Type == EnemyBat {
		// Fly straight to base
		dx, dy := baseX-e.X, baseY-e.Y
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist < speed {
			e.X = baseX
			e.Y = baseY
			return true
		}
		e.X += (dx / dist) * speed
		e.Y += (dy / dist) * speed
	} else {
		// Follow path
		if e.PathIndex < len(e.Path) {
			targetX := float64(e.Path[e.PathIndex][0]*32 + 16)
			targetY := float64(e.Path[e.PathIndex][1]*32 + 16)
			
			dx, dy := targetX-e.X, targetY-e.Y
			dist := math.Sqrt(dx*dx + dy*dy)
			
			if dist < speed {
				e.X = targetX
				e.Y = targetY
				e.PathIndex++
			} else {
				e.X += (dx / dist) * speed
				e.Y += (dy / dist) * speed
			}
		} else {
			return true // Reached end of path (Base)
		}
	}

	e.AnimTick++
	if e.AnimTick > 10 {
		e.AnimTick = 0
		e.AnimFrame = (e.AnimFrame + 1) % 4
	}
	return false
}

// Update Projectile
func (p *Projectile) Update(enemies []*Enemy, units []*Unit, spawnParticle func(*Particle)) bool {
	// Returns true if hit
	var targetX, targetY float64
	if p.Target != nil {
		targetX, targetY = p.Target.X, p.Target.Y
	} else if p.TargetUnit != nil {
		targetX, targetY = p.TargetUnit.X, p.TargetUnit.Y
	} else {
		targetX, targetY = p.TargetX, p.TargetY
	}

	dx, dy := targetX-p.X, targetY-p.Y
	dist := math.Sqrt(dx*dx + dy*dy)

	if dist < p.Speed {
		// Hit!
		if p.Target != nil {
			p.Target.Health -= p.Damage
			if p.IsSlow {
				p.Target.SlowDuration = 120 // 2 seconds
			}
			if p.IsSplash {
				for _, e := range enemies {
					if Distance(p.Target.X, p.Target.Y, e.X, e.Y) < 64 {
						e.Health -= p.Damage * 0.5 // 50% splash damage
					}
				}
				// Spawn explosion particles
				for i:=0; i<10; i++ {
					spawnParticle(NewParticle(p.Target.X, p.Target.Y, 0xFFA500FF)) // Orange
				}
			}
		} else if p.TargetUnit != nil {
			if p.Type == ProjHealBeam {
				p.TargetUnit.Health -= p.Damage // negative damage heals
				if p.TargetUnit.Health > p.TargetUnit.MaxHealth {
					p.TargetUnit.Health = p.TargetUnit.MaxHealth
				}
				spawnParticle(NewParticle(p.TargetUnit.X, p.TargetUnit.Y, 0x00FF00FF)) // Green
			} else {
				p.TargetUnit.Health -= p.Damage
				if p.IsSplash {
					for _, u := range units {
						if Distance(p.TargetUnit.X, p.TargetUnit.Y, u.X, u.Y) < 64 {
							u.Health -= p.Damage * 0.5
						}
					}
				}
			}
		}
		return true
	}

	p.X += (dx / dist) * p.Speed
	p.Y += (dy / dist) * p.Speed
	return false
}

// Update Particle
func (p *Particle) Update() bool {
	// Returns true if dead
	p.X += p.VX
	p.Y += p.VY
	p.Life++
	return p.Life >= p.MaxLife
}

