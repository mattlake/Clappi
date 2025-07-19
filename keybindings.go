package main

import "github.com/gdamore/tcell/v2"

type KeyBinding struct {
	key     tcell.Key
	handler func()
}

type RuneBinding struct {
	rune    rune
	handler func()
}

func (tui *ClappiTUI) setupKeyBindings() {
	tui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if handled := tui.handleSpecialKeys(event); handled {
			return nil
		}

		if handled := tui.handleVimKeys(event); handled {
			return nil
		}

		return event
	})
}

func (tui *ClappiTUI) handleSpecialKeys(event *tcell.EventKey) bool {

	keyBindings := map[tcell.Key]func(){
		tcell.KeyEscape: tui.app.Stop,
	}

	if handler, exists := keyBindings[event.Key()]; exists {
		handler()
		return true
	}
	return false
}

func (tui *ClappiTUI) handleVimKeys(event *tcell.EventKey) bool {
	vimBindings := map[rune]func(){
		'j': tui.moveFocusDown,
		'k': tui.moveFocusUp,
		'h': tui.moveFocusLeft,
		'l': tui.moveFocusRight,
		'q': tui.app.Stop,
	}

	if handler, exists := vimBindings[event.Rune()]; exists {
		handler()
		return true
	}

	return false
}

func (tui *ClappiTUI) moveFocusDown() {
	if tui.currentPanel < len(tui.panels)-1 {
		tui.setFocus(tui.currentPanel + 1)
	}
}

func (tui *ClappiTUI) moveFocusUp() {
	if tui.currentPanel > 0 {
		tui.setFocus(tui.currentPanel - 1)
	}
}

func (tui *ClappiTUI) moveFocusLeft() {
	if tui.currentPanel >= mainPanelsStartIndex {
		tui.setFocus(tui.currentPanel - mainPanelsStartIndex)
	}
}

func (tui *ClappiTUI) moveFocusRight() {
	if tui.currentPanel < mainPanelsStartIndex {
		tui.setFocus(tui.currentPanel + mainPanelsStartIndex)
	}
}
