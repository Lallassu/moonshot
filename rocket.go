package main

import (
	"math/rand"
	"time"

	"github.com/nullboundary/glfont"
)

type Rocket struct {
	Object
	released    bool
	hasReleased bool
	failed      bool
	boost       float64
	landed      bool
	landX       float64
	landY       float64
	boostFont   *glfont.Font
	maxBoost    float64
	initSleepMs int64
	startTS     time.Time
}

func (r *Rocket) Init(x, y, z float64, img string, objType ObjectType) {
	r.startTS = time.Now()
	font, err := glfont.LoadFont(conf.Assets["menuFont"], int32(50), screenWidth, screenHeight)
	if err != nil {
		panic(err)
	}
	r.boostFont = font

	r.Object.Init(x, y, z, img, objType)
}

func (r *Rocket) Update(dt float64) {
	if r.removed || time.Since(r.startTS).Milliseconds() < r.initSleepMs {
		return
	}

	if r.landed {
		r.active = false
		r.x = moon.x - r.landX
		r.y = moon.y - r.landY
	} else {
		if r.vy < -5 && r.hasReleased {
			r.chunk.rotationDeg += float32(dt)
			sound.Stop("liftoff")
			r.failed = true
			r.mass = 0.5
		} else {
			Smoke(
				r.chunk.x+(float64(r.chunk.sizex)*float64(r.chunk.scale)/2),
				r.chunk.y-10,
				10,
			)
			// smoke
			for i := 0; i < int(r.boost)*40; i++ {
				color := rand.Float32() * 0xFFFF
				particles.NewParticle(particle{
					r:    color,
					g:    color,
					b:    color,
					a:    color,
					size: float64(1 + rand.Intn(2)),
					Phys: Phys{
						x:           r.chunk.x + (float64(r.chunk.sizex) * float64(r.chunk.scale) / 2),
						y:           r.chunk.y + 2,
						vy:          r.boost/2 - rand.Float64()*r.boost,
						vx:          r.boost*2 - rand.Float64()*r.boost*4,
						fx:          r.boost*2 - rand.Float64()*r.boost*4,
						fy:          r.boost/2 - rand.Float64()*r.boost,
						life:        rand.Float64() * 1,
						mass:        -0.1,
						restitution: 0,
						active:      true,
					},
				})
			}
			if r.boost > 0 || r.hasReleased {
				for i := 0; i < 50; i++ {
					// Fire jet
					cr := float32(rand.Intn(0xFFFFFF))
					cg := float32(rand.Intn(0x33555))
					cb := float32(0)
					ca := float32(0xFFF + rand.Intn(0xFFFFFF))
					particles.NewParticle(particle{
						r:    cr * 2,
						g:    cg,
						b:    cb,
						a:    ca,
						size: float64(1 + rand.Intn(2)),
						Phys: Phys{
							x:           r.chunk.x + (float64(r.chunk.sizex) * float64(r.chunk.scale) / 2),
							y:           r.chunk.y - 3,
							vy:          2 - rand.Float64()*4,
							vx:          2 - rand.Float64()*4,
							fx:          2 - rand.Float64()*4,
							fy:          2 - rand.Float64()*4,
							life:        rand.Float64(),
							mass:        1,
							restitution: -0.2,
							active:      true,
						},
					})

					// Blue intensive jet
					cr = 0
					cg = float32(rand.Intn(0x33555))
					cb = float32(rand.Intn(0xFFFFFF))
					ca = float32(0xFFF + rand.Intn(0xFFFFFF))
					particles.NewParticle(particle{
						r:    cr * 2,
						g:    cg,
						b:    cb,
						a:    ca,
						size: float64(1 + rand.Intn(2)),
						Phys: Phys{
							x:           r.chunk.x + (float64(r.chunk.sizex) * float64(r.chunk.scale) / 2),
							y:           r.chunk.y,
							vy:          2 - rand.Float64()*4,
							vx:          2 - rand.Float64()*4,
							fx:          2 - rand.Float64()*4,
							fy:          2 - rand.Float64()*4,
							life:        rand.Float64() / 4,
							mass:        1,
							restitution: -0.2,
							active:      true,
						},
					})
				}
			}
		}

		if r.released {
			sound.Play("liftoff", 0.5)
			r.fy = r.boost
			r.vy = r.boost
			r.released = false
			r.boost = 0
			r.vx = gameMap.Wind
		}

		if (r.hit || r.y < 1) && r.failed && r.hasReleased {
			Explode(r.x, r.y, 100)
			r.Hit(int(r.x), int(r.y), ObjectWorld)
		}
	}

	r.Object.Update(dt)

	// Don't fall out of the map.
	if r.y < 0 {
		r.y = 0
		r.active = false
	}
}

func (r *Rocket) Release() {
	if r.boost > 0 {
		sound.Stop("thrusters")
		r.released = true
		r.hasReleased = true
		r.active = true
		// Just so we don't explode at start if at bottom
		r.y += 1
		r.vy = 0
	}
}

func (r *Rocket) Boost() {
	if !r.removed && !r.hasReleased && time.Since(r.startTS).Milliseconds() > r.initSleepMs {
		if r.boost == 0 {
			sound.Play("thrusters", 0.1)
		}
		if r.boost < r.maxBoost {
			sound.Volume("thrusters", r.boost*1/r.maxBoost)
			r.boost += 0.2
		}
		if r.boost > r.maxBoost {
			r.boost = r.maxBoost
		}
	}
}

func (r *Rocket) Hit(x, y int, objType ObjectType) {
	if r.landed {
		return
	}
	if objType == ObjectMoon {
		// Check if bottom of rocket, then land!
		r.landX = moon.x - float64(x)
		r.landY = moon.y - float64(y)
		r.landed = true
		gameMap.Landed()
		return
	}
	r.removed = true
	r.Explode(int(r.x), int(r.y))
	sound.Play("explosion1", 0.6)
	sound.Stop("liftoff")
	gameMap.Reset()
}
