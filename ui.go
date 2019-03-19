package main

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type Components struct {
	app         *tview.Application
	root        *tview.Flex
	header      *tview.TextView
	side        *tview.List
	main        *tview.TextView
	mainAuthor  *tview.TextView
	footerTopic *tview.TextView
	footerHelp  *tview.TextView
}

var UI Components

var headerTpl = `
    [red]♫  ♪ ♫  ♪ [yellow]%s

    [green]顺序循环  (%s)

`

var footerTpl = `ENTER [green]播放[white]  SPACE [green]暂停[white] CTRL-N/CTRL-P [green]下一首/上一首[white]  `
var normalTpl = "[yellow]% 3d) \t[-:-:-]%s"
var playingTpl = "[purple]% 3d) \t->[::b] %s"

func init() {
	UI = Components{
		app:        tview.NewApplication(),
		root:       tview.NewFlex().SetDirection(tview.FlexRow),
		header:     tview.NewTextView().SetDynamicColors(true),
		main:       tview.NewTextView().SetDynamicColors(true),
		mainAuthor: tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignRight),
		side:       tview.NewList().SetSelectedBackgroundColor(tcell.ColorDimGray),
		footerTopic: tview.NewTextView().
			SetDynamicColors(true).SetRegions(true).SetWrap(false).SetTextAlign(tview.AlignCenter),
		footerHelp: tview.NewTextView().
			SetDynamicColors(true).SetRegions(true).SetWrap(false).SetTextAlign(tview.AlignRight).
			SetText(footerTpl),
	}
	UI.root.SetRect(0, 0, 100, 40)
	UI.header.SetBorder(true).SetTitle("即刻电台 " + version)
	UI.side.SetBorder(true).SetBorderPadding(1, 0, 0, 1)

	mainContainer := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(UI.main, 0, 1, false).
		AddItem(UI.mainAuthor, 2, 1, false)
	mainContainer.SetBorder(true).SetBorderPadding(1, 0, 2, 2)

	footerContainer := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(UI.footerTopic, 40, 1, false).
		AddItem(UI.footerHelp, 0, 1, false)

	footerContainer.SetBackgroundColor(tcell.ColorGrey)
	UI.footerHelp.SetBackgroundColor(tcell.ColorGrey)
	UI.footerTopic.SetBackgroundColor(tcell.ColorGrey)

	UI.root.
		AddItem(UI.header, 7, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(UI.side, 40, 1, true).
			AddItem(mainContainer, 0, 1, false),
			0, 1, false).
		AddItem(footerContainer, 1, 1, false)
	UI.app.SetRoot(UI.root, false).SetFocus(UI.side)
}

func normalText(index int, title string) string {
	return fmt.Sprintf(normalTpl, index, title)
}

func playingText(index int, title string) string {
	return fmt.Sprintf(playingTpl, index, title)
}

func updateTotalSong(count int) {
	UI.side.SetTitle(fmt.Sprintf(" 歌曲数: %d", count))
}

func run() error {
	return UI.app.Run()
}
