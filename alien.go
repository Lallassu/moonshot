package main

import (
	"math"
	"math/rand"
)

type Alien struct {
	Object
	AmmoType         ShotType
	AmmoSpeed        float64
	AmmoFreq         int
	AmmoMaxTime      float64
	WaitForRelease   bool
	CanBeHitByDebris bool
}

func (a *Alien) Update(dt float64) {
	if a.removed {
		return
	}

	a.active = false

	a.rotation += dt / (rand.Float64() * 50000)
	a.x += math.Cos(a.rotation)
	a.y += math.Sin(a.rotation)

	for i := 0; i < 2; i++ {
		particles.NewParticle(particle{
			r:    0xFFFFFF,
			g:    0xFFFFFF,
			b:    0xFFFFFF,
			a:    0xFFFFFF,
			size: float64(1 + rand.Intn(2)),
			Phys: Phys{
				x:           10 - rand.Float64()*20 + a.chunk.x + float64(a.chunk.sizex)/4,
				y:           10 - rand.Float64()*20 + a.chunk.y + float64(a.chunk.sizey)/4,
				vy:          1 - rand.Float64()*6,
				vx:          1 - rand.Float64()*6,
				fx:          1 - rand.Float64()*6,
				fy:          1 - rand.Float64()*6,
				life:        rand.Float64() / 2,
				mass:        1,
				restitution: -0.2,
				active:      true,
			},
		})
	}
	a.Object.Update(dt)

	if !rocket.removed && rand.Intn(100) > (100-a.AmmoFreq) {
		if (a.WaitForRelease && rocket.hasReleased) || !a.WaitForRelease {
			s := &Shot{Type: a.AmmoType, Speed: a.AmmoSpeed, MaxTime: a.AmmoMaxTime}
			s.Init(a.chunk.x, a.chunk.y, 2, "", ObjectShot)
			objects = append(objects, s)
		}
	}

}

func (a *Alien) Hit(x, y int, objType ObjectType) {
	if (a.CanBeHitByDebris && objType == ObjectDebris) || objType == ObjectRocket {
		a.Explode(int(a.x), int(a.y))
		a.removed = true
	}
}
