package main

import (
	"math"
)

type Moon struct {
	Object
	released bool
	boost    float64
	speed    float64
}

func (m *Moon) Update(dt float64) {
	m.active = false
	m.rotation += dt * m.speed / 1000
	m.x = 600 + 300*math.Cos(m.rotation)
	m.y = 800 + 100*math.Sin(m.rotation)
	m.Object.Update(dt)
}

func (m *Moon) Hit(x, y int, objType ObjectType) {
	if objType == ObjectSatellite || objType == ObjectAlien || objType == ObjectDebris {
		return
	}
	if objType == ObjectRocket && rocket.landed {
		return
	}
	m.Explode(x, y)
}
