package main

import (
	"github.com/go-gl/gl/all-core/gl"
)

const (
	maxParticles = 100000
)

type ParticlePool struct {
	particles       []particle
	idx             int
	shader          *Shader
	triangles       int
	activeParticles int
	pbo             uint32
	cbo             uint32
	vao             uint32
	positions       []float32
	colors          []float32
}

func (pp *ParticlePool) Init() {
	pp.particles = make([]particle, maxParticles)

	for i := 0; i < maxParticles; i++ {
		p := particle{}
		pp.particles = append(pp.particles, p)
	}
	pp.idx = 0

	shader, err := NewShader(conf.Shaders["particle_vs"], conf.Shaders["particle_fs"])
	if err != nil {
		panic(err)
	}

	pp.shader = shader

	gl.GenVertexArrays(1, &pp.vao)
	gl.BindVertexArray(pp.vao)

	gl.GenBuffers(1, &pp.pbo)
	gl.GenBuffers(1, &pp.cbo)

	gl.BindBuffer(gl.ARRAY_BUFFER, pp.pbo)
	gl.BufferData(gl.ARRAY_BUFFER, maxParticles*4*4, nil, gl.STREAM_DRAW)
	gl.VertexAttribPointer(0, 4, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.BindBuffer(gl.ARRAY_BUFFER, pp.cbo)
	gl.BufferData(gl.ARRAY_BUFFER, maxParticles*4*4, nil, gl.STREAM_DRAW)
	gl.VertexAttribPointer(1, 4, gl.FLOAT, true, 0, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(1)

	gl.VertexAttribDivisor(0, 1)
	gl.VertexAttribDivisor(1, 1)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
}

func (pp *ParticlePool) NewParticle(p particle) {
	// Don't draw off-screen
	if p.x > screenWidth+30 || p.x < -30 || p.y > screenHeight+30 || p.y < -30 {
		return
	}
	pp.idx++
	if pp.idx >= maxParticles {
		pp.idx = 0
	}
	newp := pp.particles[pp.idx]

	if p.size <= 0 {
		p.size = 1
	}

	newp = p
	if p.z == 0 {
		newp.z = 2
	}
	newp.active = true

	pp.particles[pp.idx] = newp
}

func (pp *ParticlePool) Update(dt float64) {
	pp.activeParticles = 0
	pp.positions = []float32{}
	pp.colors = []float32{}

	for i := range pp.particles {
		if pp.particles[i].active {

			pp.particles[i].Update(dt)
			pp.activeParticles++
			pp.positions = append(pp.positions, []float32{
				float32(pp.particles[i].x),
				float32(pp.particles[i].y),
				float32(pp.particles[i].z),
				float32(pp.particles[i].size),
			}...)
			pp.colors = append(pp.colors, []float32{
				float32(pp.particles[i].r),
				float32(pp.particles[i].g),
				float32(pp.particles[i].b),
				float32(pp.particles[i].a),
			}...)
		}
	}
	if pp.activeParticles > 0 {
		gl.BindBuffer(gl.ARRAY_BUFFER, pp.pbo)
		gl.BufferData(gl.ARRAY_BUFFER, maxParticles*4*4, nil, gl.STREAM_DRAW)
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, pp.activeParticles*4*4, gl.Ptr(pp.positions))

		gl.BindBuffer(gl.ARRAY_BUFFER, pp.cbo)
		gl.BufferData(gl.ARRAY_BUFFER, maxParticles*4*4, nil, gl.STREAM_DRAW)
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, pp.activeParticles*4*4, gl.Ptr(pp.colors))
	}
}

func (pp *ParticlePool) Draw(dt float64) {
	pp.Update(dt)
	pp.shader.Use()

	if err := pp.shader.SetUniformMatrixName("projection", false, projection); err != nil {
		panic(err)
	}

	if err := pp.shader.SetUniformMatrixName("view", false, view); err != nil {
		panic(err)
	}

	// gl.Enable(gl.BLEND)
	// gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	gl.BindVertexArray(pp.vao)
	gl.DrawArraysInstanced(gl.POINTS, 0, 1, int32(pp.activeParticles))
	gl.BindVertexArray(0)
}
