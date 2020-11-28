package main

import (
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Background struct {
	shader   *Shader
	vbo      uint32
	vao      uint32
	ebo      uint32
	vertices []float32
	indices  []uint32
	texture  Texture
}

func (b *Background) Clear() {
	gl.DeleteBuffers(1, &b.vbo)
	gl.DeleteVertexArrays(1, &b.vao)
	gl.DeleteBuffers(1, &b.ebo)
}

func (b *Background) Init(bg string) {
	b.shader = shaders[ObjectBackground]

	b.vertices = []float32{
		0.5, 0.5, 0.0, 1.0, 0.0, 0.0, 1.0, 1.0,
		0.5, -0.5, 0.0, 0.0, 1.0, 0.0, 1.0, 0.0,
		-0.5, -0.5, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0,
		-0.5, 0.5, 0.0, 1.0, 1.0, 0.0, 0.0, 1.0,
	}

	b.indices = []uint32{
		0, 1, 3,
		1, 2, 3,
	}

	gl.GenVertexArrays(1, &b.vao)
	gl.GenBuffers(1, &b.vbo)
	gl.GenBuffers(1, &b.ebo)

	gl.BindVertexArray(b.vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, b.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(b.vertices)*4, gl.Ptr(b.vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, b.ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(b.indices)*4, gl.Ptr(b.indices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 8*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 8*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 8*4, gl.PtrOffset(6*4))
	gl.EnableVertexAttribArray(2)

	b.texture = NewTexture2D()
	b.texture.Use()
	b.texture.SetParameter(gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	b.texture.SetParameter(gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	_, err := b.texture.Load(conf.Assets[bg], false, true)
	if err != nil {
		panic(err)
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
}

func (b *Background) Draw(dt float64) {
	if b.texture == nil {
		return
	}
	b.texture.Use()
	b.shader.Use()

	translate := mgl32.Translate3D(screenWidth/2, screenHeight/2, -1)
	scale := mgl32.Scale3D(screenWidth, screenHeight, 1)
	trans := translate.Mul4(scale)

	if err := b.shader.SetUniformMatrixName("model", false, trans); err != nil {
		panic(err)
	}

	if err := b.shader.SetUniformMatrixName("projection", false, projection); err != nil {
		panic(err)
	}

	if err := b.shader.SetUniformMatrixName("view", false, view); err != nil {
		panic(err)
	}
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	gl.BindVertexArray(b.vao)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, b.ebo)
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))
	gl.BindVertexArray(0)
}
