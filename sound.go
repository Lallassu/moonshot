package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

type Sound struct {
	sounds map[string]*Snd
}

func (s *Sound) Init() {
	s.sounds = make(map[string]*Snd)
	sr := beep.SampleRate(48000)
	speaker.Init(sr, sr.N(time.Second/100))
}

type Snd struct {
	buffer *beep.Buffer
	ctrl   *beep.Ctrl
	vol    *effects.Volume
}

func (s *Sound) Load(file, name string) {
	f, _ := os.Open(file)
	if strings.Contains(file, "mp3") {
		streamer, format, err := mp3.Decode(f)
		if err != nil {
			panic("Failed to load sound")
		}

		fmt.Printf("Loading sound: %v....", name)
		buff := beep.NewBuffer(format)
		buff.Append(streamer)
		streamer.Close()
		fmt.Printf("Done.\n")

		s.sounds[name] = &Snd{buffer: buff}
	}
}

func (s *Sound) Play(name string, vol float64) {
	if vol > 1 {
		vol = 1
	}
	sound := s.sounds[name].buffer.Streamer(0, s.sounds[name].buffer.Len())
	ctrl := &beep.Ctrl{Streamer: beep.Loop(1, sound), Paused: false}
	volume := &effects.Volume{
		Streamer: ctrl,
		Base:     1,
		Volume:   vol,
		Silent:   false,
	}
	sn := s.sounds[name]
	sn.ctrl = ctrl
	sn.vol = volume
	speaker.Play(volume)
}

func (s *Sound) Stop(name string) {
	if s.sounds[name].ctrl == nil || s.sounds[name].ctrl.Paused {
		return
	}
	s.sounds[name].ctrl.Paused = true
}

func (s *Sound) Volume(name string, vol float64) {
	if vol > 1 {
		vol = 1
	}
	speaker.Lock()
	s.sounds[name].vol.Volume = vol
	speaker.Unlock()
}
