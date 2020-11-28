package main

import (
	"fmt"
	"time"

	"github.com/nullboundary/glfont"
)

const (
	maxFPSProbes = 60
)

type Stats struct {
	currFPS         int
	fpsProbes       []int
	probeCount      int
	avgFPS          int
	avgDrawTime     time.Duration
	avgChunkRebuild time.Duration
	prevTime        time.Time
	font            *glfont.Font
}

func (s *Stats) Init() {
	s.fpsProbes = make([]int, maxFPSProbes)

	font, err := glfont.LoadFont(conf.Assets["statsFont"], int32(12), screenWidth, screenHeight)
	if err != nil {
		panic(err)
	}
	s.font = font
}

func (s *Stats) Draw() {
	strs := []string{}

	tri := world.Triangles()
	ab := world.ActiveBlocks()
	geff := 100 - ((float64(tri) / float64(ab*2)) * 100.0)
	tot := world.totalBlocks
	if tot == 0 {
		tot = 1
	}

	strs = append(strs, []string{
		fmt.Sprintf("FPS: %d", s.currFPS),
		fmt.Sprintf("Avg. FPS: %d", s.avgFPS),
		fmt.Sprintf("Chunks: %d", world.totalChunks),
		fmt.Sprintf("Dirty Chunks: %d", world.DirtyChunks()),
		fmt.Sprintf("Triangles: %d", tri+particles.triangles),
		fmt.Sprintf("Active Blocks: %d (%d)", ab, int((ab/tot)*100)),
		fmt.Sprintf("Total Blocks: %d", world.totalBlocks),
		fmt.Sprintf("Greedy Efficiency: %f", geff),
		fmt.Sprintf("Particles: %d/%d", particles.activeParticles, maxParticles),
	}...)

	s.font.SetColor(1.0, 1.0, 1.0, 1.0)
	for i, l := range strs {
		s.font.Printf(5, 20+float32(i*20), 1.0, l)
	}
}

func (s *Stats) Update() {
	s.currFPS = int(1000000 / time.Since(s.prevTime).Microseconds())
	s.fpsProbes[s.probeCount] = s.currFPS

	if s.probeCount >= maxFPSProbes-1 {
		s.avgFPS = 0
		for _, v := range s.fpsProbes {
			s.avgFPS += v
		}
		s.avgFPS /= maxFPSProbes
		s.probeCount = 0
	} else {
		s.probeCount++
	}

	s.prevTime = time.Now()
}
