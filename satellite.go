package main

import (
	"math"
	"math/rand"
)

type Satellite struct {
	Object
	released bool
	boost    float64
	origX    float64
	origY    float64
}

func (s *Satellite) Update(dt float64) {
	if s.removed {
		return
	}

	s.active = false
	s.rotation += dt / (10000 + rand.Float64()*1000)
	s.x = s.origX + 1000*math.Cos(s.rotation)
	s.y = s.origY + 1000*math.Sin(s.rotation)

	for i := 0; i < 1; i++ {
		cr := float32(0xFFFFFF)
		cg := float32(0xFFFFFF)
		cb := float32(0xFFFFFF)
		ca := float32(0xFFF + rand.Intn(0xFFFFFF))
		particles.NewParticle(particle{
			r:    cr * 2,
			g:    cg,
			b:    cb,
			a:    ca,
			size: float64(1 + rand.Intn(2)),
			Phys: Phys{
				x:           5 - rand.Float64()*10 + s.chunk.x + (float64(s.chunk.sizex) * float64(s.chunk.scale) / 2),
				y:           5 - rand.Float64()*10 + s.chunk.y + (float64(s.chunk.sizey) * float64(s.chunk.scale) / 2),
				vy:          1 - rand.Float64()*2,
				vx:          1 - rand.Float64()*2,
				fx:          1 - rand.Float64()*2,
				fy:          1 - rand.Float64()*2,
				life:        rand.Float64() / 3,
				mass:        1,
				restitution: -0.2,
				active:      true,
			},
		})
	}
	s.Object.Update(dt)
}

func (s *Satellite) Hit(x, y int, objType ObjectType) {
	if objType == ObjectMoon {
		return
	}
	s.removed = true
	s.Explode(int(s.x), int(s.y))
}
