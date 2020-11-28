package main

import (
	"image"
	"image/jpeg"
	"image/png"
	"os"
)

func LoadTexture(file string) (img image.Image, width, height, size float64, err error) {
	width = 0.0
	height = 0.0
	size = 0.0
	img = nil
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	image.RegisterFormat("jpg", "jpg", jpeg.Decode, jpeg.DecodeConfig)

	imgfile, err := os.Open(file)
	if err != nil {
		return
	}

	defer imgfile.Close()

	imgCfg, _, err := image.DecodeConfig(imgfile)
	if err != nil {
		return
	}

	imgfile.Seek(0, 0)
	img, _, _ = image.Decode(imgfile)

	height = float64(imgCfg.Height)
	width = float64(imgCfg.Width)
	size = width
	if width < height {
		size = height
	}
	return
}

func LoadObject(file string, shader *Shader, z float64) Chunk {
	img, width, height, _, err := LoadTexture(file)
	if err != nil {
		panic(err)
	}

	c := Chunk{shader: shader}
	c.Init(int(width+1), int(height+1), 100, 100, z, false)

	for x := 0; x <= int(width); x++ {
		for y := 0; y <= int(height); y++ {
			r, g, b, a := img.At(x, int(height)-y).RGBA()
			c.Add(x, y, float32(r), float32(g), float32(b), float32(a))
		}
	}

	return c
}

func LoadMap(file string) error {
	img, width, height, _, err := LoadTexture(file)
	if err != nil {
		return err
	}

	for x := 0; x <= int(width); x++ {
		for y := 0; y <= int(height); y++ {
			r, g, b, a := img.At(x, int(height)-y).RGBA()
			world.Add(x, y, float32(r), float32(g), float32(b), float32(a))
		}
	}

	return nil
}
