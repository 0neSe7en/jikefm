package main

import (
	"fmt"
	"github.com/rivo/tview"
)

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

var headerTemplate = `
    [red]♫  ♪ ♫  ♪ [yellow]%s

    [green]顺序循环  (%s)

`

var normalTemplate = "[yellow]% 3d) \t[-:-:-]%s"
var playingTemplate = "[purple]% 3d) \t->[::b] %s"

func init() {
	UI.root.SetRect(0, 0, 100, 40)
	UI.header.SetBorder(true).SetTitle("JIKE FM")
	UI.side.SetBorder(true).SetBorderPadding(1, 0, 0, 1)
	UI.main.SetBorder(true).SetBorderPadding(1, 0, 2, 2)
	UI.root.
		AddItem(UI.header, 7, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(UI.side, 40, 1, true).
			AddItem(UI.main, 0, 1, false),
			0, 1, false)
	UI.app.SetRoot(UI.root, false).SetFocus(UI.side)
}

func normalText(index int, title string) string {
	return fmt.Sprintf(normalTemplate, index + 1, title)
}

func playingText(index int, title string) string {
	return fmt.Sprintf(playingTemplate, index + 1, title)
}

func run() error {
	return UI.app.Run()
}
