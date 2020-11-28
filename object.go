package main

import (
	"fmt"
	"math/rand"
	"time"
)

type CustomFunc func(*Object, float64)
type Object struct {
	Phys
	chunk    Chunk
	shader   *Shader
	removed  bool
	objType  ObjectType
	boost    float64
	started  time.Time
	released bool
}

func (r *Object) Init(x, y, z float64, img string, objType ObjectType) {
	r.started = time.Now()
	r.chunk.x = x
	r.chunk.y = y
	r.chunk.z = z
	r.active = true
	r.objType = objType

	if _, ok := shaders[objType]; !ok {
		panic(fmt.Sprintf("Shader no found for object: %v", objType))
	}
	r.shader = shaders[objType]

	r.chunk = LoadObject(img, r.shader, z)
	r.x = x
	r.y = y
	r.mass = 2
	r.keepAlive = true
	r.restitution = -0.1
	r.active = true

	r.borderPixels = r.chunk.GetBorderPixels()
}

func (r *Object) Clear() {

}

func (r *Object) GetBorderPixels() []float64 {
	return r.borderPixels
}

func (r *Object) GetX() float64 {
	return r.x
}

func (r *Object) GetObjType() ObjectType {
	return r.objType
}

func (r *Object) GetY() float64 {
	return r.y
}

func (r *Object) IsActive(x, y int) bool {
	return r.chunk.IsActive(x, y, true)
}

func (r *Object) Update(dt float64) {
	r.Phys.Update(dt)
	r.Draw(dt)
}

func (r *Object) Remove() {
	r.removed = true
}

func (r *Object) Draw(dt float64) {
	if r.removed {
		return
	}

	r.chunk.x = r.x
	r.chunk.y = r.y

	// Don't draw off-screen except rocket
	if r.objType != ObjectRocket {
		if r.x > screenWidth+30 || r.x < -30 || r.y > screenHeight+30 || r.y < -30 {
			return
		}
	}
	r.chunk.Draw(dt)
}

func (r *Object) IsRemoved() bool {
	return r.removed
}

func (r *Object) Hit(x, y int, objType ObjectType) {
}

func (r *Object) Explode(x, y int) {
	power := 10 // TBD
	pow := power * power
	for rx := x - power; rx <= x+power; rx++ {
		vx := (rx - x) * (rx - x)
		for ry := y - power; ry <= y+power; ry++ {
			if ry < 0 {
				continue
			}
			val := (ry-y)*(ry-y) + vx
			if val < pow {
				b := r.chunk.GetBlock(rx, ry, true)
				if !b.used {
					continue
				}
				r.chunk.Remove(rx, ry, true)

				life := 5 + rand.Float64()*3
				for i := 0; i < 5; i++ {
					particles.NewParticle(particle{
						r:    b.r,
						g:    b.g,
						b:    b.b,
						a:    b.a,
						size: float64(1 + rand.Intn(2)),
						Phys: Phys{
							x:           float64(rx),
							y:           float64(ry) - rand.Float64()*10,
							vy:          10 - rand.Float64()*20,
							vx:          10 - rand.Float64()*20,
							fx:          10 + rand.Float64()*20,
							fy:          10 + rand.Float64()*20,
							life:        life,
							mass:        1,
							restitution: -0.2,
							active:      true,
						},
					})
				}
			}
		}
	}

}
