package main

import (
	"fmt"
	"github.com/0neSe7en/jikefm/jike"
	"github.com/0neSe7en/jikefm/musicapi"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/gdamore/tcell"
	"strings"
	"time"
	"unicode"
)

var topics = map[string]string{
	"55483ddee4b03ccbae843925": "晚安电台",
	"5a1ccf936e6e7c0011037480": "即友在听什么歌",
}

var CurrentSession *jike.Session
var fm = newFm()

type Music struct {
	url    string
	index  int
}

type JikeFm struct {
	playlist     []jike.Message
	player *Player
	currentMusic Music
	currentTopic string
	nextMusicIndex int
	skip         string
}

func newFm() *JikeFm {
	p := &JikeFm{
		currentTopic:"55483ddee4b03ccbae843925",
	}
	p.player = newPlayer(p.iter)
	p.nextMusicIndex = 0
	return p
}

func (p *JikeFm) feed() {
	res, next, _ := jike.FetchMoreSelectedFM(CurrentSession, p.currentTopic, p.skip)
	for _, msg := range res {
		p.playlist = append(p.playlist, msg)
	}
	p.skip = next
}

func (p *JikeFm) play() {
	p.player.open()
}

func (p *JikeFm) playIndex(next int) beep.Streamer {
	mp3Url := musicapi.NeteaseUrlToMp3(p.playlist[next].LinkInfo.LinkUrl)
	f, err := musicapi.NeteaseDownload(mp3Url)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	current := p.currentMusic.index
	p.currentMusic = Music{
		url:    mp3Url,
		index:  next,
	}
	UI.app.QueueUpdateDraw(func () {
		p.drawHeader()
		p.changeSong(current, next)
	})
	return p.player.playMp3(f)
}

func (p *JikeFm) onSelectChange(index int, _ string, _ string, _ rune) {
	i := index
	if index < 0 {
		i = len(p.playlist) + index
	}
	if index >= len(p.playlist) {
		i = index - len(p.playlist)
	}
	msg := p.playlist[i]
	UI.main.SetText(msg.Content)
	UI.mainAuthor.SetText("[green]@" + msg.User.ScreenName)
}

func (p *JikeFm) onEnterPlay(index int, mainText string, secondaryText string, shortcut rune) {
	p.nextMusicIndex = index
	p.player.reset().open()
}

func (p *JikeFm) drawHeader() {
	text := fmt.Sprintf(headerTpl,
		p.playlist[p.currentMusic.index].GetTitle(),
		p.player.currentPosition(),
	)
	UI.header.SetText(text)
}

func (p *JikeFm) changeSong(from int, target int) {
	if from >= 0 {
		UI.side.SetItemText(
			from,
			normalText(from, p.playlist[from].GetTitle()),
			"",
		)
	}
	UI.side.SetItemText(
		target,
		playingText(target, p.playlist[target].GetTitle()),
		"",
	)
}

func (p *JikeFm) iter() beep.Streamer {
	if len(p.playlist) == 0 {
		return nil
	}
	stream := p.playIndex(p.nextMusicIndex)
	var next int
	if next = p.nextMusicIndex + 1; next >= len(p.playlist) {
		next = 0
	}
	p.nextMusicIndex = next

	return stream
}

func (p *JikeFm) handle(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyESC:
		UI.app.Stop()
	case tcell.KeyRune:
		switch unicode.ToLower(event.Rune()) {
		case ' ':
			p.player.togglePlay()
		}
	}
	return event
}

func main() {
	CurrentSession = jike.NewSession()
	_ = speaker.Init(targetFormat.SampleRate, targetFormat.SampleRate.N(time.Second/30))

	defer fm.player.close()

	for {
		fm.feed()
		if len(fm.playlist) > 20 {
			break
		}
	}

	fm.onSelectChange(0, "", "", 0)

	for index, msg := range fm.playlist {
		UI.side.AddItem(normalText(index, msg.GetTitle()), "", 0, nil)
	}

	UI.side.
		SetChangedFunc(fm.onSelectChange).
		SetSelectedFunc(fm.onEnterPlay).
		ShowSecondaryText(false)
	UI.app.SetInputCapture(fm.handle)

	var s []string
	for id, topicName := range topics {
		s = append(s, fmt.Sprintf(`["%s"]%s[""]`, id, topicName))
	}
	UI.footerTopic.SetText("| " + strings.Join(s, " | ") + " |")

	fm.play()

	go func() {
		for {
			time.Sleep(time.Second)
			go UI.app.QueueUpdateDraw(func () {
				fm.drawHeader()
			})
		}
	}()

	if err := run(); err != nil {
		panic(err)
	}
}
