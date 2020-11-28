package main

import (
	"math/rand"
)

// Phys is a generic physics simulation for both particles
// and chunks.
type Phys struct {
	x            float64
	y            float64
	borderPixels []float64
	rotation     float64

	restitution  float64
	vx           float64
	vy           float64
	prevX        float64
	prevY        float64
	mdt          float64
	life         float64
	mass         float64
	active       bool
	fx           float64
	fy           float64
	keepAlive    bool
	hit          bool
	explodeOnHit int
}

func (p *Phys) Update(dt float64) {
	if !p.active {
		return
	}

	if p.life <= 0 && !p.keepAlive {
		p.active = false
		return
	}

	dt /= 1000

	p.life -= dt
	ax := p.fx * dt * p.vx * p.mass
	ay := p.fy * dt * p.vy * p.mass

	p.prevX = p.x
	p.prevY = p.y

	hit := false
	p.hit = false
	if p.explodeOnHit > 0 {
		for i := 0; i < 2; i++ {
			color := rand.Float32() * 0xFFFF
			particles.NewParticle(particle{
				r:    color,
				g:    color,
				b:    color,
				a:    color,
				size: float64(1 + rand.Intn(2)),
				Phys: Phys{
					x:           p.x,
					y:           p.y,
					vy:          0.5 - rand.Float64(),
					vx:          0.5 - rand.Float64(),
					fx:          0.5 - rand.Float64(),
					fy:          0.5 - rand.Float64(),
					life:        rand.Float64() / 10,
					mass:        -0.3,
					restitution: 0,
					active:      true,
				},
			})
		}
	}
	if len(p.borderPixels) > 0 {
		for i := 0; i < len(p.borderPixels); i += 2 {
			x := p.x + p.borderPixels[i]
			y := p.y + p.borderPixels[i+1]

			if world.IsActive(int(x+ax), int(y+ay)) {
				hit = true
				p.hit = true
				// TBD: p.rotation = ?
				break
			}
		}
	} else {
		if world.IsActive(int(p.x+ax), int(p.y+ay)) {
			hit = true
			p.hit = true
			if p.explodeOnHit > 0 {
				world.Explode(int(p.x+ax), int(p.y+ay), p.explodeOnHit)
				p.active = false
			}
		}
	}
	if hit {
		if p.vy < 0 {
			p.vy *= p.restitution * rand.Float64()
		} else {
			p.vx *= p.restitution * rand.Float64()
			p.vy *= p.restitution * rand.Float64()
		}
	} else {
		p.x += ax
		p.y += ay
	}

	if !world.IsActive(int(p.x), int(p.y-1)) {
		p.vy -= dt * p.fy
		p.fx += dt * gravity * p.mass
		p.fy += dt * gravity * p.mass
	}

	if p.prevX-p.x == 0 && p.prevY-p.y == 0 {
		p.mdt += dt
	} else {
		p.mdt = 0
	}
}
