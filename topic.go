package main

import (
	"fmt"
	"github.com/0neSe7en/jikefm/jike"
)

type Topic struct {
	id   string
	name string
	playing int
	playlist []jike.Message
	skip              string
	more              chan bool
	ui *TopicUI
}

type TopicList map[string]*Topic

var topicOrder = []string{
	"55483ddee4b03ccbae843925",
	"5a1ccf936e6e7c0011037480",
}

var topics = TopicList{
	"55483ddee4b03ccbae843925": {
		id: "55483ddee4b03ccbae843925",
		name: "晚安电台",
		playing: -1,
	},
	"5a1ccf936e6e7c0011037480": {
		id: "5a1ccf936e6e7c0011037480",
		name: "即友在听什么歌",
		playing: -1,
	},
}

func init() {
	for _, topic := range topics {
		topic.more = make(chan bool)
		go topic.FetchMore()
	}
}

func (t *Topic) FetchMore() {
	for {
		<- t.more
		msgs := t.Feed()
		UI.app.QueueUpdateDraw(func() {
			t.AddToPlaylist(msgs)
		})
	}
}

func (t *Topic) ChangeSong(from int, target int) {
	if from >= 0 {
		t.ui.side.SetItemText(
			from,
			normalText(from + 1, t.playlist[from].GetTitle()),
			"",
		)
	}

	if target >= 0 {
		t.ui.side.SetItemText(
			target,
			playingText(target + 1, t.playlist[target].GetTitle()),
			"",
		)
	}
}

func (t *Topic) AddToPlaylist(messages []jike.Message) {
	for _, msg := range messages {
		t.playlist = append(t.playlist, msg)
		t.ui.side.AddItem(normalText(len(t.playlist), msg.GetTitle()), "", 0, nil)
	}
	t.ui.side.SetTitle(fmt.Sprintf(" 歌曲数: %d", len(t.playlist)))
}

func (t *Topic) BindUI(ui *TopicUI) *Topic {
	t.ui = ui
	t.ui.side.
		SetChangedFunc(t.onSelect).
		ShowSecondaryText(false)
	return t
}

func (t *Topic) SetSelectedFunc(handler func (topicId string, index int)) *Topic {
	if t.ui == nil { return t}
	t.ui.side.SetSelectedFunc(func (index int, _ string, _ string, _ rune) {
		handler(t.id, index)
	})
	return t
}

func (t *Topic) fetch() []jike.Message {
	res, next, _ := jike.FetchMoreSelectedFM(CurrentSession, t.id, t.skip)
	t.skip = next
	return res
}

func (t *Topic) Feed() []jike.Message {
	msgs := t.fetch()
	for len(msgs) < 20 {
		for _, msg := range t.fetch() {
			msgs = append(msgs, msg)
		}
	}
	return msgs
}

func (t *Topic) onSelect(index int, _ string, _ string, _ rune) {
	i := index
	if index < 0 {
		i = len(t.playlist) + index
	}
	if index >= len(t.playlist) {
		i = index - len(t.playlist)
	}
	if index == len(t.playlist)-1 {
		t.queueMore()
	}
	msg := t.playlist[i]
	t.ui.content.SetText(msg.Content)
	t.ui.author.SetText("[green]@" + msg.User.ScreenName)
}

func (t *Topic) queueMore() {
	select {
	case t.more <- true:
	default:
	}
}
