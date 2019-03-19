package main

import (
	"fmt"
	"github.com/0neSe7en/jikefm/jike"
	"github.com/0neSe7en/jikefm/musicapi"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/gdamore/tcell"
	"strconv"
	"strings"
	"time"
	"unicode"
)

var CurrentSession *jike.Session
var fm = newFm()

type Music struct {
	url     string
	topicId string
	current int
	next    int
}

type JikeFm struct {
	player            *Player
	music             *Music
	viewingTopicIndex int
	timer             *time.Timer
}

func newFm() *JikeFm {
	p := &JikeFm{
		viewingTopicIndex: 0,
	}
	p.player = newPlayer(p.iter)
	p.music = &Music{
		topicId: topicOrder[0],
		current: 0,
		next: 0,
	}
	return p
}

func (p *JikeFm) playingTopic() *Topic {
	return topics[p.music.topicId]
}

func (p *JikeFm) viewingTopic() *Topic {
	return topics[topicOrder[p.viewingTopicIndex]]
}

func (p *JikeFm) playingList() []jike.Message {
	return p.playingTopic().playlist
}

func (p *JikeFm) setTimer(d time.Duration) {
	if p.timer != nil {
		p.stopTimer()
	}
	p.timer = time.AfterFunc(d, func() {
		UI.app.Stop()
	})
}

func (p *JikeFm) stopTimer() {
	if p.timer != nil {
		p.timer.Stop()
		p.timer = nil
	}
}

func (p *JikeFm) play() {
	p.player.open()
}

func (p *JikeFm) onEnterPlay(topicId string, index int) {
	prevTopic := topics[p.music.topicId]
	prevIndex := p.music.current
	p.music = &Music{
		topicId: topicId,
		next:    index,
		current: 0,
	}
	UI.app.QueueUpdateDraw(func() {
		prevTopic.ChangeSong(prevIndex, -1)
	})
	p.player.reset().open()
}

func (p *JikeFm) drawHeader() {
	text := fmt.Sprintf(headerTpl,
		p.playingList()[p.music.current].GetTitle(),
		p.player.currentPosition(),
		"",
	)
	UI.header.SetText(text)
}

func (p *JikeFm) calcNextIndex() int {
	next := p.music.current + 1
	if next > len(p.playingList())-1 {
		p.playingTopic().queueMore()
	}
	if next >= len(p.playingList()) {
		next = 0
	}
	return next
}

func (p *JikeFm) iter() beep.Streamer {
	if len(p.playingList()) == 0 {
		return nil
	}
	mp3Url := p.playingList()[p.music.next].Mp3Url
	f, err := musicapi.NeteaseDownload(mp3Url)
	if err != nil {
		return nil
	}
	c := p.music.current
	p.music.url = mp3Url
	p.music.current = p.music.next
	p.music.next = p.calcNextIndex()

	UI.app.QueueUpdateDraw(func() {
		p.drawHeader()
		topics[p.music.topicId].ChangeSong(c, p.music.current)
	})
	return p.player.playMp3(f)
}

func (p *JikeFm) handle(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyCtrlN:
		p.music.next = p.calcNextIndex()
		p.player.reset().open()
		return nil
	case tcell.KeyCtrlP:
		n := p.music.current - 1
		if n < 0 {
			n = 0
		}
		p.music.next = n
		p.player.reset().open()
		return nil
	case tcell.KeyTab:
		p.viewingTopicIndex += 1
		if p.viewingTopicIndex == len(topicOrder) {
			p.viewingTopicIndex = 0
		}
		p.switchTopic()
		return nil
	case tcell.KeyRune:
		switch unicode.ToLower(event.Rune()) {
		case ' ':
			p.player.togglePlay()
			return nil
		case 't':
			return nil
		}
	}
	return event
}

func (p *JikeFm) switchTopic() {
	UI.footerTopic.Highlight(strconv.Itoa(fm.viewingTopicIndex))
	UI.topicPages.SwitchToPage(topicOrder[fm.viewingTopicIndex])
	UI.app.SetFocus(UI.topics[fm.viewingTopicIndex].side)
}

func (p *JikeFm) highlight() {
}

func main() {
	CurrentSession = jike.NewSession()
	_ = speaker.Init(targetFormat.SampleRate, targetFormat.SampleRate.N(time.Second/30))

	var s []string
	for index, topicId := range topicOrder {
		t := topics[topicId]
		t.BindUI(addTopic(topicId)).SetSelectedFunc(fm.onEnterPlay)
		s = append(s, fmt.Sprintf(`["%d"]%s[""]`, index, t.name))
		t.AddToPlaylist(t.Feed())
		t.onSelect(0, "", "", 0)
	}

	defer fm.player.close()

	UI.app.SetInputCapture(fm.handle)
	UI.footerTopic.SetText("| " + strings.Join(s, " | ") + " |")

	fm.switchTopic()
	fm.play()

	go func() {
		for {
			time.Sleep(time.Second)
			go UI.app.QueueUpdateDraw(func() {
				fm.drawHeader()
			})
		}
	}()

	if err := run(); err != nil {
		panic(err)
	}
}
