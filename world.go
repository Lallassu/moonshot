package main

import ()

const (
	chunkSize = 128
	blockSize = 1
	gravity   = 9.82
)

type World struct {
	cx          int
	cy          int
	chunks      [][]*Chunk
	shader      *Shader
	totalBlocks int
	triangles   int
	totalChunks int
}

// Init
func (w *World) Init(sizex, sizey int) {
	w.cx = int(sizex / chunkSize)
	w.cy = int(sizey / chunkSize)

	w.chunks = make([][]*Chunk, w.cx)
	for i := 0; i < w.cx; i++ {
		w.chunks[i] = make([]*Chunk, w.cy)
	}

	for x := 0; x < w.cx; x++ {
		for y := 0; y < w.cy; y++ {
			w.totalChunks++
			w.chunks[x][y] = &Chunk{shader: w.shader}
			w.chunks[x][y].Init(chunkSize, chunkSize, float64(x*chunkSize), float64(y*chunkSize), 0, true)
			w.totalBlocks += chunkSize * chunkSize
		}
	}
}

func (w *World) RedrawAll() {
	for x := 0; x < w.cx; x++ {
		for y := 0; y < w.cy; y++ {
			w.chunks[x][y].dirty = true
		}
	}
}

func (w *World) Clear() {
	for x := 0; x < w.cx; x++ {
		for y := 0; y < w.cy; y++ {
			w.chunks[x][y].Clear()
		}
	}
}

func (w *World) Draw(dt float64) {
	for x := 0; x < w.cx; x++ {
		for y := 0; y < w.cy; y++ {
			w.chunks[x][y].Draw(dt)
		}
	}
}

// Add
func (w *World) Add(x, y int, r, g, b, a float32) {
	cix := int(x / chunkSize)
	ciy := int(y / chunkSize)

	if cix < 0 || ciy < 0 || cix >= len(w.chunks) || ciy >= len(w.chunks[cix]) {
		return
	}

	w.chunks[cix][ciy].Add(x-(cix*chunkSize), y-(ciy*chunkSize), r, g, b, a)
}

// Remove
func (w *World) Remove(x, y int) {
	cix := int(x / chunkSize)
	ciy := int(y / chunkSize)

	if cix < 0 || ciy < 0 || cix >= len(w.chunks) || ciy >= len(w.chunks[cix]) {
		return
	}

	w.chunks[cix][ciy].Remove(x-(cix*chunkSize), y-(ciy*chunkSize), false)
}

// IsWall checks if a block is active or not.
func (w *World) IsActive(x, y int) bool {
	cix := int(x / chunkSize)
	ciy := int(y / chunkSize)

	if cix < 0 || ciy < 0 || cix >= len(w.chunks) || ciy >= len(w.chunks[cix]) {
		return false
	}

	return w.chunks[cix][ciy].IsActive(x-(cix*chunkSize), y-(ciy*chunkSize), false)
}

func (w *World) DirtyChunks() int {
	tot := 0
	for x := 0; x < w.cx; x++ {
		for y := 0; y < w.cy; y++ {
			if w.chunks[x][y].dirty {
				tot++
			}
		}
	}
	return tot
}

func (w *World) ActiveBlocks() int {
	tot := 0
	for x := 0; x < w.cx; x++ {
		for y := 0; y < w.cy; y++ {
			tot += w.chunks[x][y].activeBlocks
		}
	}
	return tot
}

func (w *World) Triangles() int {
	tot := 0
	for x := 0; x < w.cx; x++ {
		for y := 0; y < w.cy; y++ {
			tot += w.chunks[x][y].triangles
		}
	}
	return tot
}

func (w *World) Explode(x, y, power int) {
	pow := power * power
	for rx := x - power; rx <= x+power; rx++ {
		vx := (rx - x) * (rx - x)
		for ry := y - power; ry <= y+power; ry++ {
			if ry < 0 {
				continue
			}
			val := (ry-y)*(ry-y) + vx
			if val < pow {
				w.Remove(rx, ry)
			}
		}
	}
}

func (w *World) Pencil(x, y, power int) {
	pow := power * power
	for rx := x - power; rx <= x+power; rx++ {
		vx := (rx - x) * (rx - x)
		for ry := y - power; ry <= y+power; ry++ {
			if ry < 0 {
				continue
			}
			val := (ry-y)*(ry-y) + vx
			if val < pow {
				w.Add(rx, ry, 0xFFFF, 0, 0xFFFF, 0xFFFFF)
			}
		}
	}
}
