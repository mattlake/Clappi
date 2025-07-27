package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (tui *ClappiTUI) setupKeyBindings() {
	tui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		focusedPanel := tui.app.GetFocus()

		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'h':
				tui.focusPrevPanel()
			case 'j':
				tui.nextListItem(focusedPanel)
			case 'k':
				tui.prevListItem(focusedPanel)
			case 'l':
				tui.focusNext()
			case 'q':
				tui.app.Stop()
			}
		}

		return event
	})
}

func (tui *ClappiTUI) nextListItem(focusedPanel tview.Primitive) {
	if list, ok := focusedPanel.(*tview.List); ok {
		currentIndex := list.GetCurrentItem()
		if currentIndex < list.GetItemCount()-1 {
			list.SetCurrentItem(currentIndex + 1)
		}
	}
}

func (tui *ClappiTUI) prevListItem(focusedPanel tview.Primitive) {
	if list, ok := focusedPanel.(*tview.List); ok {
		currentIndex := list.GetCurrentItem()
		if currentIndex > 0 {
			list.SetCurrentItem(currentIndex - 1)
		}
	}
}

func (tui *ClappiTUI) focusPrevPanel() {
	if tui.currentPanel > 0 {
		tui.setFocus(tui.currentPanel - 1)
	}
}

func (tui *ClappiTUI) focusNext() {
	if tui.currentPanel < mainPanelsStartIndex-1 {
		tui.setFocus(tui.currentPanel + 1)
	}
}
