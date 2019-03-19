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

type Topic struct {
	id   string
	name string
}

var topics = []Topic{
	{id: "55483ddee4b03ccbae843925", name: "晚安电台"},
	{id: "5a1ccf936e6e7c0011037480", name: "即友在听什么歌"},
}

var CurrentSession *jike.Session
var fm = newFm()

type Music struct {
	url   string
	index int
}

type JikeFm struct {
	playlist          []jike.Message
	player            *Player
	currentMusic      Music
	currentTopicIndex int
	nextMusicIndex    int
	skip              string
	more              chan bool
}

func newFm() *JikeFm {
	p := &JikeFm{
		currentTopicIndex: 0,
		more:              make(chan bool, 1),
	}
	p.player = newPlayer(p.iter)
	p.nextMusicIndex = 0
	go p.fetchMore()
	return p
}

func (p *JikeFm) fetchMore() {
	for {
		<-p.more
		msgs := p.feed()

		UI.app.QueueUpdateDraw(func() {
			currentLen := len(p.playlist)
			p.addToPlaylist(msgs)
			updateTotalSong(len(p.playlist))
			if p.nextMusicIndex == 0 {
				p.nextMusicIndex = currentLen
			}
		})
	}
}

func (p *JikeFm) feed() []jike.Message {
	res, next, _ := jike.FetchMoreSelectedFM(CurrentSession, topics[p.currentTopicIndex].id, p.skip)
	p.skip = next
	return res
}

func (p *JikeFm) addToPlaylist(messages []jike.Message) {
	for _, msg := range messages {
		p.playlist = append(p.playlist, msg)
		UI.side.AddItem(normalText(len(p.playlist), msg.GetTitle()), "", 0, nil)
	}
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
		url:   mp3Url,
		index: next,
	}
	UI.app.QueueUpdateDraw(func() {
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
	if index == len(p.playlist)-1 {
		p.queueMore()
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
			normalText(from+1, p.playlist[from].GetTitle()),
			"",
		)
	}
	UI.side.SetItemText(
		target,
		playingText(target+1, p.playlist[target].GetTitle()),
		"",
	)
}

func (p *JikeFm) queueMore() {
	select {
	case p.more <- true:
	default:
	}
}

func (p *JikeFm) calcNextIndex() int {
	next := p.currentMusic.index + 1
	if next > len(p.playlist)-1 {
		p.queueMore()
	}
	if next >= len(p.playlist) {
		next = 0
	}
	return next
}

func (p *JikeFm) iter() beep.Streamer {
	if len(p.playlist) == 0 {
		return nil
	}
	stream := p.playIndex(p.nextMusicIndex)
	p.nextMusicIndex = p.calcNextIndex()
	return stream
}

func (p *JikeFm) handle(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyESC:
		UI.app.Stop()
	case tcell.KeyCtrlN:
		p.nextMusicIndex = p.calcNextIndex()
		p.player.reset().open()
		return nil
	case tcell.KeyCtrlP:
		var n int
		if n := p.currentMusic.index - 1; n < 0 {
			n = 0
		}
		p.nextMusicIndex = n
		p.player.reset().open()
		return nil
	case tcell.KeyTab:
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

func main() {
	CurrentSession = jike.NewSession()
	_ = speaker.Init(targetFormat.SampleRate, targetFormat.SampleRate.N(time.Second/30))

	defer fm.player.close()

	fm.addToPlaylist(fm.feed())
	updateTotalSong(len(fm.playlist))

	fm.onSelectChange(0, "", "", 0)

	UI.side.
		SetChangedFunc(fm.onSelectChange).
		SetSelectedFunc(fm.onEnterPlay).
		ShowSecondaryText(false)
	UI.app.SetInputCapture(fm.handle)

	var s []string
	for index, topic := range topics {
		s = append(s, fmt.Sprintf(`["%d"]%s[""]`, index, topic.name))
	}
	UI.footerTopic.SetText("| " + strings.Join(s, " | ") + " |")
	UI.footerTopic.Highlight(strconv.Itoa(fm.currentTopicIndex))

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
