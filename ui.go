package main

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type TopicUI struct {
	root    *tview.Flex
	content *tview.TextView
	author  *tview.TextView
	side    *tview.List
}

type Components struct {
	app         *tview.Application
	root        *tview.Flex
	header      *tview.TextView
	topicPages  *tview.Pages
	topics      []*TopicUI
	footerTopic *tview.TextView
	footerHelp  *tview.TextView
}

var UI Components

var headerTpl = `
    [red]♫  ♪ ♫  ♪ [yellow]%s

    [green]顺序循环  (%s)  %s

`

var footerTpl = `TAB [green]切换圈子[white] ENTER [green]播放[white] SPACE [green]暂停[white] CTRL-N/CTRL-P [green]下/上一首[white]  `
var normalTpl = "[yellow]% 3d) \t[-:-:-]%s"
var playingTpl = "[purple]% 3d) \t->[::b] %s"

func init() {
	UI = Components{
		app:        tview.NewApplication(),
		root:       tview.NewFlex().SetDirection(tview.FlexRow),
		header:     tview.NewTextView().SetDynamicColors(true),
		topicPages:		tview.NewPages(),
		footerTopic: tview.NewTextView().
			SetDynamicColors(true).SetRegions(true).SetWrap(false).SetTextAlign(tview.AlignCenter),
		footerHelp: tview.NewTextView().
			SetDynamicColors(true).SetRegions(true).SetWrap(false).SetTextAlign(tview.AlignRight).
			SetText(footerTpl),
	}
	UI.root.SetRect(0, 0, 110, 40)
	UI.header.SetBorder(true).SetTitle(" 即刻电台 " + version)

	footerContainer := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(UI.footerTopic, 44, 1, false).
		AddItem(UI.footerHelp, 0, 1, false)

	footerContainer.SetBackgroundColor(tcell.ColorGrey)
	UI.footerHelp.SetBackgroundColor(tcell.ColorGrey)
	UI.footerTopic.SetBackgroundColor(tcell.ColorGrey)

	UI.root.
		AddItem(UI.header, 7, 1, false).
		AddItem(UI.topicPages, 0, 1, false).
		AddItem(footerContainer, 1, 1, false)
	UI.app.SetRoot(UI.root, false)
}

func addTopic(topicId string) *TopicUI {
	topicUI := &TopicUI{
		root: tview.NewFlex().SetDirection(tview.FlexColumn),
		side: tview.NewList().SetSelectedBackgroundColor(tcell.ColorDimGray),
		content: tview.NewTextView().SetDynamicColors(true),
		author: tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignRight),
	}
	topicUI.side.SetBorder(true).SetBorderPadding(1, 0, 0, 1)
	mainContainer := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(topicUI.content, 0, 1, false).
		AddItem(topicUI.author, 2, 1, false)
	mainContainer.SetBorder(true).SetBorderPadding(1, 0, 2, 2)

	topicUI.root.
		AddItem(topicUI.side, 44, 1, true).
		AddItem(mainContainer, 0, 1, false)

	UI.topics = append(UI.topics, topicUI)
	UI.topicPages.AddPage(topicId, topicUI.root, true, false)
	return topicUI
}

func normalText(index int, title string) string {
	return fmt.Sprintf(normalTpl, index, title)
}

func playingText(index int, title string) string {
	return fmt.Sprintf(playingTpl, index, title)
}

func run() error {
	return UI.app.Run()
}
