package main

import (
	"math/rand"
)

func Smoke(x, y, power float64) {
	for i := 0; i < 50; i++ {
		// smoke
		color := rand.Float32() * 0xFFFF
		particles.NewParticle(particle{
			r:    color,
			g:    color,
			b:    color,
			a:    color,
			size: float64(1 + rand.Intn(2)),
			Phys: Phys{
				x:           x,
				y:           y,
				vy:          power/2 - rand.Float64()*power,
				vx:          power/2 - rand.Float64()*power,
				fx:          power/2 - rand.Float64()*power,
				fy:          power/2 - rand.Float64()*power,
				life:        rand.Float64() * 3,
				mass:        -0.2,
				restitution: 0,
				active:      true,
			},
		})
	}
}

func Explode(x, y, power float64) {
	world.Explode(int(x), int(y), int(power))
	for i := 0; i < int(power)*50; i++ {
		// smoke
		color := rand.Float32() * 0xFFFF

		particles.NewParticle(particle{
			r:    color,
			g:    color,
			b:    color,
			a:    color,
			size: float64(1 + rand.Intn(2)),
			Phys: Phys{
				x:           x,
				y:           y,
				vy:          power/6 - rand.Float64()*power/3,
				vx:          power/6 - rand.Float64()*power/3,
				fx:          power/6 - rand.Float64()*power/3,
				fy:          power/6 - rand.Float64()*power/3,
				life:        rand.Float64() * 2,
				mass:        -0.1,
				restitution: 0,
				active:      true,
			},
		})
		// Fire
		cr := float32(rand.Intn(0xFFFFFF))
		cg := float32(rand.Intn(0x33555))
		cb := float32(0)
		ca := float32(0xFFF + rand.Intn(0xFFFFFF))

		// Some random exploding parts.
		life := rand.Float64()
		expHit := 0
		if power > 3 {
			if rand.Intn(100) > 97 {
				life += 2
				expHit = 2 + rand.Intn(3)
			}
		}

		particles.NewParticle(particle{
			r:    cr * 2,
			g:    cg,
			b:    cb,
			a:    ca,
			size: float64(1 + rand.Intn(3)),
			Phys: Phys{
				explodeOnHit: expHit,
				x:            x,
				y:            y,
				vy:           power/4 - rand.Float64()*power/2,
				vx:           power/4 - rand.Float64()*power/2,
				fx:           power/4 - rand.Float64()*power/2,
				fy:           power/4 - rand.Float64()*power/2,
				life:         life,
				mass:         1,
				restitution:  -0.1,
				active:       true,
			},
		})
	}
}

func Star() {
	particles.NewParticle(particle{
		r:    0xFFFFFF,
		g:    0xFFFFFF,
		b:    0xFFFFFF,
		a:    0xFFFF,
		z:    0,
		size: rand.Float64() * 5,
		Phys: Phys{
			x:           screenWidth * rand.Float64(),
			y:           screenHeight - rand.Float64()*screenHeight/3,
			vy:          0,
			vx:          rand.Float64() * 10,
			fx:          rand.Float64() * 10,
			fy:          0,
			life:        3 + rand.Float64()*2,
			mass:        0,
			restitution: 0,
			active:      true,
		},
	})
}
