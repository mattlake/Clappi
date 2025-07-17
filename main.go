package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var defaultBorderColor = tcell.ColorWhite
var focusedBorderColor = tcell.ColorYellow

func main() {
	app := tview.NewApplication()

	sidebarTop := tview.NewTextView().SetText("Sidebar top").
		SetTextAlign(tview.AlignCenter).
		SetBorder(true).
		SetTitle("Panel 1")
	sidebarMiddle := tview.NewTextView().SetText("Sidebar middle").
		SetTextAlign(tview.AlignCenter).
		SetBorder(true).
		SetTitle("Panel 2")
	sidebarBottom := tview.NewTextView().SetText("Sidebar bottom").
		SetTextAlign(tview.AlignCenter).
		SetBorder(true).
		SetTitle("Panel 3")

	mainTop := tview.NewTextView().SetText("Main top").
		SetTextAlign(tview.AlignCenter).
		SetBorder(true).
		SetTitle("Main 1")
	mainBottom := tview.NewTextView().SetText("Main bottom").
		SetTextAlign(tview.AlignCenter).
		SetBorder(true).
		SetTitle("Main bottom")

	sidebar := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(sidebarTop, 0, 1, false).
		AddItem(sidebarMiddle, 0, 1, false).
		AddItem(sidebarBottom, 0, 1, false)
	mainContent := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(mainTop, 0, 1, false).
		AddItem(mainBottom, 0, 1, false)

	root := tview.NewFlex().
		AddItem(sidebar, 30, 1, false).
		AddItem(mainContent, 0, 2, false)

	panels := []tview.Primitive{
		sidebarTop, sidebarMiddle, sidebarBottom,
		mainTop, mainBottom,
	}

	currentPanel := 0

	setFocus := func(index int) {
		// Reset the border colors
		for _, p := range panels {
			p.(*tview.Box).SetBorderColor(defaultBorderColor)
		}

		// highlight the focused panel
		currentPanel = index
		if currentPanel >= len(panels) {
			currentPanel = 0
		} else if currentPanel < 0 {
			currentPanel = len(panels) - 1
		}
		panels[currentPanel].(*tview.Box).SetBorderColor(focusedBorderColor)
		app.SetFocus(panels[currentPanel])
	}

	setFocus(0)

	// Key bindings
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		// Standard keys, I may bin these tbh
		case tcell.KeyEscape:
			app.Stop()
		case tcell.KeyTab:
			setFocus(currentPanel + 1)
			return nil
		case tcell.KeyBacktab:
			setFocus(currentPanel - 1)
			return nil
		}

		// Vim Keybindings
		switch event.Rune() {
		case 'j':
			if currentPanel < len(panels)-1 {
				setFocus(currentPanel + 1)
			}
			return nil
		case 'k':
			if currentPanel > 0 {
				setFocus(currentPanel - 1)
			}
			return nil
		case 'h':
			if currentPanel >= 3 {
				setFocus(currentPanel - 3)
			}
			return nil
		case 'l':
			if currentPanel < 3 {
				setFocus(currentPanel + 3)
			}
			return nil
		case 'q':
			app.Stop()
			return nil
		}

		return event
	})

	if err := app.SetRoot(root, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
