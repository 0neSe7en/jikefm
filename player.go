package main

import (
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"io"
	"time"
)

var targetFormat = beep.Format{
	SampleRate:  44100,
	NumChannels: 2,
	Precision:   2,
}

type IterFunc func() beep.Streamer

type Player struct {
	ctrl          *beep.Ctrl
	volume        *effects.Volume
	streamer      beep.Streamer
	currentStream beep.StreamSeeker
	onNext        IterFunc
}

func newPlayer(onNext IterFunc) *Player {
	p := &Player{}
	p.onNext = onNext
	return p.reset()
}

func (p *Player) reset() *Player {
	speaker.Clear()
	p.streamer = beep.Iterate(p.onNext)
	p.ctrl = &beep.Ctrl{Streamer: p.streamer}
	return p
}

func (p *Player) open() *Player {
	speaker.Play(p.ctrl)
	return p
}

func (p *Player) close() *Player {
	speaker.Clear()
	return p
}

func (p *Player) togglePlay() {
	speaker.Lock()
	p.ctrl.Paused = !p.ctrl.Paused
	speaker.Unlock()
}

func (p *Player) playMp3(file io.ReadCloser) beep.Streamer {
	streamer, _, err := mp3.Decode(file)
	if err != nil {
		return nil
	}
	p.currentStream = streamer
	return streamer
}

func (p *Player) currentPosition() string {
	speaker.Lock()
	pos := targetFormat.SampleRate.D(p.currentStream.Position())
	speaker.Unlock()
	minutes := pos / time.Minute
	second := pos % time.Minute / time.Second
	return fmt.Sprintf("%02d:%02d", minutes, second)
}
