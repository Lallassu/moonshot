package main

import (
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Chunk struct {
	rotationDeg  float32
	vertices     []float32
	indices      []uint32
	scale        float32
	sizex        int
	sizey        int
	x            float64
	y            float64
	z            float64
	dirty        bool
	vbo          uint32
	vao          uint32
	ebo          uint32
	shader       *Shader
	blocks       [][]block
	triangles    int
	activeBlocks int
	static       bool
}

type block struct {
	bType int
	r     float32
	g     float32
	b     float32
	a     float32
	drawn bool
	used  bool
}

func (c *Chunk) Init(sizex, sizey int, posx, posy, posz float64, static bool) {
	c.static = static
	c.x = posx
	c.y = posy
	c.z = posz
	c.sizex = sizex
	c.sizey = sizey
	c.scale = 1.0

	c.vertices = []float32{0}
	c.indices = []uint32{0}

	c.blocks = make([][]block, sizex)
	for i := 0; i < sizex; i++ {
		c.blocks[i] = make([]block, sizey)
	}

	for x := 0; x < c.sizex; x++ {
		for y := 0; y < c.sizey; y++ {
			c.blocks[x][y] = block{}
		}
	}

	// Blocks
	gl.GenVertexArrays(1, &c.vao)
	gl.GenBuffers(1, &c.vbo)
	gl.GenBuffers(1, &c.ebo)

	gl.BindVertexArray(c.vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, c.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(c.vertices)*4, gl.Ptr(c.vertices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 7*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 4, gl.FLOAT, false, 7*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, c.ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(c.indices)*6, gl.Ptr(c.indices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
}

func (c *Chunk) Clear() {
	gl.DeleteVertexArrays(1, &c.vao)
	gl.DeleteBuffers(1, &c.vbo)
	gl.DeleteBuffers(1, &c.ebo)
}

func (c *Chunk) Draw(dt float64) {
	c.shader.Use()

	translate := mgl32.Translate3D(float32(c.x), float32(c.y), float32(c.z))
	trans := translate
	scale := mgl32.Scale3D(c.scale, c.scale, c.scale)

	if !c.static {
		rot := mgl32.HomogRotate3D(float32(mgl32.DegToRad(float32(c.rotationDeg))), mgl32.Vec3{0.0, 0.0, 1.0})
		trans = translate.Mul4(scale).Mul4(rot)
	} else {
		trans = translate.Mul4(scale)
	}

	if err := c.shader.SetUniformMatrixName("model", false, trans); err != nil {
		panic(err)
	}

	if err := c.shader.SetUniformMatrixName("projection", false, projection); err != nil {
		panic(err)
	}

	if err := c.shader.SetUniformMatrixName("view", false, view); err != nil {
		panic(err)
	}

	c.Update(dt)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	gl.BindVertexArray(c.vao)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, c.ebo)
	gl.DrawElements(gl.TRIANGLES, int32(len(c.indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))
	gl.BindVertexArray(0)
}

func (c *Chunk) Update(dt float64) {
	if !c.dirty {
		return
	}

	c.triangles = 0
	c.activeBlocks = 0
	c.vertices = []float32{}
	c.indices = []uint32{}

	// Clear drawn
	for x := 0; x < c.sizex; x++ {
		for y := 0; y < c.sizey; y++ {
			c.blocks[x][y].drawn = false
		}
	}

	px := 1
	py := 1
	max := 0
	tx := 0
	tmp := 0
	maxX := 0
	r := float32(0)
	g := float32(0)
	b := float32(0)
	a := float32(0)

	for xx := 0; xx < c.sizex; xx++ {
		for yy := 0; yy < c.sizey; yy++ {
			if !c.blocks[xx][yy].used || c.blocks[xx][yy].drawn {
				continue
			}

			r = c.blocks[xx][yy].r
			g = c.blocks[xx][yy].g
			b = c.blocks[xx][yy].b
			a = c.blocks[xx][yy].a
			px = 1
			py = 1
			max = 0
			maxX = 0
			for y := 0; y < c.sizey-yy; y++ {
				tx = 0
				if maxX == 0 {
					maxX = c.sizex - xx
				}
				for x := 0; x < maxX; x++ {
					if !c.blocks[xx+x][yy+y].used ||
						c.blocks[xx+x][yy+y].drawn ||
						c.blocks[xx+x][yy+y].r != r ||
						c.blocks[xx+x][yy+y].g != g ||
						c.blocks[xx+x][yy+y].b != b ||
						c.blocks[xx+x][yy+y].a != a {
						maxX = x - 1
						break
					}
					tx++
				}

				tmp = tx * (1 + y)
				if tmp > max {
					max = tmp
					px = tx
					py = y + 1
				}
			}

			for i := xx; i < xx+px; i++ {
				for j := yy; j < yy+py; j++ {
					c.activeBlocks++
					c.blocks[i][j].drawn = true
				}
			}

			x1 := float32(xx)
			x2 := float32(xx + px*blockSize)
			y1 := float32(yy)
			y2 := float32(yy + py*blockSize)

			c.triangles += 2
			c.vertices = append(c.vertices, []float32{
				x1, y1, 0,
				c.blocks[xx][yy].r,
				c.blocks[xx][yy].g,
				c.blocks[xx][yy].b,
				c.blocks[xx][yy].a,
				x2, y1, 0,
				c.blocks[xx][yy].r,
				c.blocks[xx][yy].g,
				c.blocks[xx][yy].b,
				c.blocks[xx][yy].a,
				x1, y2, 0,
				c.blocks[xx][yy].r,
				c.blocks[xx][yy].g,
				c.blocks[xx][yy].b,
				c.blocks[xx][yy].a,
				x2, y2, 0,
				c.blocks[xx][yy].r,
				c.blocks[xx][yy].g,
				c.blocks[xx][yy].b,
				c.blocks[xx][yy].a,
			}...)

			l := uint32(len(c.indices))
			if l != 0 {
				l -= uint32(2 * len(c.indices) / 6)
			}
			c.indices = append(c.indices, []uint32{
				l, l + 1, l + 2,
				l + 2, l + 1, l + 3,
			}...)
		}
	}

	if c.triangles > 0 {
		gl.BindBuffer(gl.ARRAY_BUFFER, c.vbo)
		gl.BufferData(gl.ARRAY_BUFFER, len(c.vertices)*4, gl.Ptr(c.vertices), gl.STATIC_DRAW)

		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, c.ebo)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(c.indices)*6, gl.Ptr(c.indices), gl.STATIC_DRAW)
		c.dirty = false
	}
}

func (c *Chunk) Add(x, y int, r, g, b, a float32) {
	if x < 0 || y < 0 || x > chunkSize || y > chunkSize || a == 0 {
		return
	}
	c.blocks[x][y].r = r
	c.blocks[x][y].g = g
	c.blocks[x][y].b = b
	c.blocks[x][y].a = a
	c.blocks[x][y].used = true
	c.dirty = true
}

func (c *Chunk) Remove(x, y int, world bool) {
	if world {
		x = x - int(c.x)
		y = y - int(c.y)
	}
	if x < 0 || y < 0 || x >= c.sizex || y >= c.sizey {
		return
	}
	c.blocks[x][y].used = false
	c.dirty = true
}

func (c *Chunk) GetBlock(x, y int, world bool) block {
	if world {
		x = x - int(c.x)
		y = y - int(c.y)
	}

	if x < 0 || y < 0 || x >= c.sizex || y >= c.sizey {
		return block{}
	}
	return c.blocks[x][y]
}

func (c *Chunk) IsActive(x, y int, world bool) bool {
	if world {
		x = x - int(c.x)
		y = y - int(c.y)
	}

	if x < 0 || y < 0 || x >= c.sizex || y >= c.sizey {
		return false
	}
	return c.blocks[x][y].used
}

func (c *Chunk) GetBorderPixels() []float64 {
	pixels := []float64{}

	for xx := 0; xx < c.sizex; xx++ {
		for yy := 0; yy < c.sizey; yy++ {
			if !c.blocks[xx][yy].used {
				continue
			}

			adj := 0
			for x := 0; x <= 1; x++ {
				for y := 0; y <= 1; y++ {
					if xx+x < c.sizex && yy+y < c.sizey && xx-x >= 0 && yy-y >= 0 {
						if !c.blocks[xx+x][yy+y].used || !c.blocks[xx-x][yy-y].used {
							adj++
						}
					}
				}
			}
			if adj >= 1 {
				pixels = append(pixels, []float64{float64(xx), float64(yy)}...)
			}
		}
	}

	return pixels
}
