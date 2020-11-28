package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/go-gl/mathgl/mgl32"
)

type ShotType int

const (
	ShotRandom ShotType = iota
	ShotStraight
	ShotSeeking
)

type Shot struct {
	Object
	Type    ShotType
	Speed   float64
	MaxTime float64
	StartX  float64
	StartY  float64
}

func (s *Shot) Init(x, y, z float64, img string, objType ObjectType) {
	sName := ""
	switch s.Type {
	case ShotRandom:
		sName = "shot1"
	case ShotStraight:
		sName = "shot2"
	case ShotSeeking:
		sName = "shot3"
	default:
		sName = "shot1"
	}

	s.Object.Init(x, y, z, conf.Assets[sName], objType)
}

func (s *Shot) Update(dt float64) {
	if s.removed {
		return
	}

	// TBD: Configureable
	if time.Since(s.started).Seconds() > s.MaxTime {
		Explode(s.x, s.y, 30)
		s.removed = true
	}

	for i := 0; i < 5; i++ {
		cr := float32(0xFFFFF)
		cg := float32(0xFFFFF)
		cb := float32(0xFFFFF)
		ca := float32(0xFFFFF)
		particles.NewParticle(particle{
			r:    cr,
			g:    cg,
			b:    cb,
			a:    ca,
			size: float64(1 + rand.Intn(2)),
			Phys: Phys{
				x:           2 - rand.Float64()*4 + s.chunk.x + float64(s.chunk.sizex)/4,
				y:           2 - rand.Float64()*4 + s.chunk.y + float64(s.chunk.sizey)/4,
				vy:          1 - rand.Float64()*2,
				vx:          1 - rand.Float64()*2,
				fx:          1 - rand.Float64()*2,
				fy:          1 - rand.Float64()*2,
				life:        rand.Float64(),
				mass:        1,
				restitution: -0.2,
				active:      true,
			},
		})
	}
	s.Object.Update(dt)

	if rocket.removed && !rocket.hasReleased {
		s.Explode(int(s.x), int(s.y))
		s.removed = true
		s.active = true
	} else {
		s.Lerp(dt)
	}
}

func (s *Shot) Hit(x, y int, objType ObjectType) {
	if objType == ObjectAlien || objType == ObjectShot {
		return
	}
	s.removed = true
	Explode(s.x, s.y, 30)
	sound.Play(fmt.Sprintf("explosion%d", 2+rand.Intn(2)), 0.5)
}

func (s *Shot) Lerp(dt float64) {
	if !rocket.removed && s.Type == ShotSeeking {
		dist := math.Sqrt(math.Pow(rocket.x-s.chunk.x, 2) + math.Pow(rocket.y-s.chunk.y, 2))
		v1 := mgl32.Vec3{float32(s.chunk.x), float32(s.chunk.y), 0}
		v2 := mgl32.Vec3{float32(rocket.x), float32(rocket.y), 0}
		q1 := mgl32.Quat{W: 0, V: v1}
		q2 := mgl32.Quat{W: 0, V: v2}
		q3 := mgl32.QuatLerp(q1, q2, float32(s.Speed)/float32(dist))
		s.x = float64(q3.X())
		s.y = float64(q3.Y())
		s.active = false
	} else if !rocket.removed && s.Type == ShotStraight {
		// Initiate at first only
		if s.StartX == 0 || s.StartY == 0 {
			s.StartX = rocket.chunk.x
			s.StartY = rocket.chunk.y
		}
		s.active = false
		dist := math.Sqrt(math.Pow(s.StartX-s.chunk.x, 2) + math.Pow(s.StartY-s.chunk.y, 2))
		v1 := mgl32.Vec3{float32(s.chunk.x), float32(s.chunk.y), 0}
		v2 := mgl32.Vec3{float32(s.StartX), float32(s.StartY), 0}
		q1 := mgl32.Quat{W: 0, V: v1}
		q2 := mgl32.Quat{W: 0, V: v2}
		q3 := mgl32.QuatLerp(q1, q2, float32(s.Speed)/float32(dist))
		s.x = float64(q3.X())
		s.y = float64(q3.Y())
	} else if !rocket.removed && s.Type == ShotRandom {
		s.rotation += dt / (rand.Float64() * 5000)
		s.x += math.Cos(s.rotation)
		s.y += math.Sin(s.rotation)

	} else {
		s.active = true
		s.rotation += dt / (rand.Float64() * 5000)
		s.x += math.Cos(s.rotation) * 2
		s.y += math.Sin(s.rotation) * 2
	}

}
