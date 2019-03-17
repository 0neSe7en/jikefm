package main

import (
	"fmt"
	"github.com/0neSe7en/jikefm/jike"
	"github.com/0neSe7en/jikefm/musicapi"
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"time"
	"unicode"
)

var topic = "55483ddee4b03ccbae843925"
var CurrentSession *jike.Session

var targetFormat = beep.Format{
	SampleRate:  44100,
	NumChannels: 2,
	Precision:   2,
}

type Music struct {
	url    string
	seeker beep.StreamSeeker
	index  int
}

type Player struct {
	playlist     []jike.Message
	currentMusic Music
	streamer     beep.Streamer
	ctrl         *beep.Ctrl
	volume       *effects.Volume
	skip         string
}

var UI = struct {
	app    *tview.Application
	root   *tview.Flex
	header *tview.TextView
	side   *tview.List
	main   *tview.TextView
}{
	app:    tview.NewApplication(),
	root:   tview.NewFlex().SetDirection(tview.FlexRow),
	header: tview.NewTextView().SetDynamicColors(true),
	main:   tview.NewTextView().SetDynamicColors(true),
	side:   tview.NewList(),
}

var player = newPlayer()

func init() {
	UI.root.SetRect(0, 0, 120, 40)
	UI.header.SetBorder(true).SetTitle("JIKE FM")
	UI.side.SetBorder(true)
	UI.main.SetBorder(true)
	UI.root.
		AddItem(UI.header, 8, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(UI.side, 40, 1, true).
			AddItem(UI.main, 0, 1, false),
			0, 1, false)
}

func newPlayer() *Player {
	p := &Player{}
	return p
}

func (p *Player) feed() {
	res, next, _ := jike.FetchMoreSelectedFM(CurrentSession, topic, p.skip)
	for _, msg := range res {
		p.playlist = append(p.playlist, msg)
	}
	p.skip = next
}

func (p *Player) play() {
	p.currentMusic.index = -1
	p.streamer = beep.Iterate(player.iter)
	p.ctrl = &beep.Ctrl{Streamer: player.streamer}
	speaker.Play(p.ctrl)
}

func (p *Player) playIndex(next int) beep.Streamer {
	mp3Url := musicapi.NeteaseUrlToMp3(p.playlist[next].LinkInfo.LinkUrl)
	f, err := musicapi.NeteaseDownload(mp3Url)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	streamer, _, err := mp3.Decode(f)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	p.currentMusic = Music{
		seeker: streamer,
		url:    mp3Url,
		index:  next,
	}
	return streamer
}

func (p *Player) add() {

}

func (p *Player) close() {

}

func (p *Player) onSelectChange(index int, _ string, _ string, _ rune) {
	i := index
	if index < 0 {
		i = len(player.playlist) + index
	}
	if index >= len(player.playlist) {
		i = index - len(player.playlist)
	}
	msg := player.playlist[i]
	content := fmt.Sprintf("%s\n [yellow]by %s", msg.Content, msg.User.ScreenName)
	UI.main.SetText(content)
}

func (p *Player) onEnterPlay(index int, mainText string, secondaryText string, shortcut rune) {
	speaker.Clear()
	p.currentMusic.index = index - 1
	p.streamer = beep.Iterate(player.iter)
	p.ctrl = &beep.Ctrl{Streamer: player.streamer}
	speaker.Play(p.ctrl)
}

func (p *Player) drawHeader() {
	UI.app.QueueUpdateDraw(func() {
		speaker.Lock()
		position := targetFormat.SampleRate.D(player.currentMusic.seeker.Position())
		text := fmt.Sprintf("  %s [green]%v", p.playlist[p.currentMusic.index].LinkInfo.Title, position.Round(time.Second))
		speaker.Unlock()
		UI.header.SetText(text)
	})
}

func (p *Player) iter() beep.Streamer {
	if len(p.playlist) == 0 {
		return nil
	}
	current := p.currentMusic.index
	var next int
	if next = current + 1; next >= len(p.playlist) {
		next = 0
	}
	go p.drawHeader()
	return p.playIndex(next)
}

func (p *Player) handle(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyESC:
		UI.app.Stop()
	case tcell.KeyRune:
		switch unicode.ToLower(event.Rune()) {
		case ' ':
			speaker.Lock()
			p.ctrl.Paused = !p.ctrl.Paused
			speaker.Unlock()
		}
	}
	return event
}

func main() {
	CurrentSession = jike.NewSession()
	_ = speaker.Init(targetFormat.SampleRate, targetFormat.SampleRate.N(time.Second/30))
	defer player.close()

	for {
		player.feed()
		if len(player.playlist) > 20 {
			break
		}
	}

	player.onSelectChange(0, "", "", 0)

	for index, msg := range player.playlist {
		var shortcut rune
		if index > 9 {
			shortcut = 0
		} else {
			shortcut = rune(index + '0')
		}
		UI.side.AddItem(msg.LinkInfo.Title, "", shortcut, nil)
	}
	UI.side.
		SetChangedFunc(player.onSelectChange).
		SetSelectedFunc(player.onEnterPlay).
		ShowSecondaryText(false)
	UI.app.SetInputCapture(player.handle)
	player.play()

	go func() {
		for {
			time.Sleep(time.Second)
			player.drawHeader()
		}
	}()

	if err := UI.app.SetRoot(UI.root, false).SetFocus(UI.side).Run(); err != nil {
		panic(err)
	}
}
