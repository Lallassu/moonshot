package main

import (
	"math"
)

type Debris struct {
	Object
	released bool
	boost    float64
	origX    float64
	origY    float64
}

func (s *Debris) Update(dt float64) {
	if s.removed {
		return
	}

	s.chunk.rotationDeg += float32(dt / s.origX)
	s.active = false
	s.rotation += dt / 1000
	s.x = (moon.chunk.x + float64(moon.chunk.sizex/2)) + s.origX*math.Cos(s.rotation)
	s.y = (moon.chunk.y + float64(moon.chunk.sizey/2)) + s.origY*math.Sin(s.rotation)

	s.Object.Update(dt)
}

func (s *Debris) Hit(x, y int, objType ObjectType) {
	if objType == ObjectMoon || objType == ObjectDebris {
		return
	}
	//s.removed = true
	s.Explode(int(s.x), int(s.y))
}
