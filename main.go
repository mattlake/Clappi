package main

import (
	"Clappi/constants"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ClappiTUI struct {
	app          *tview.Application
	panels       []tview.Primitive
	currentPanel int
}

const (
	mainPanelsStartIndex = 3
	defaultBorderColor   = tcell.ColorWhite
	focusedBorderColor   = tcell.ColorYellow
)

func createTextPanelBase(title, content string) *tview.Box {
	return tview.NewTextView().
		SetText(content).
		SetTextAlign(tview.AlignCenter).
		SetBorder(true).
		SetTitle(title)
}

func newClappiTUI() *ClappiTUI {
	tui := &ClappiTUI{
		app:          tview.NewApplication(),
		currentPanel: 0,
	}
	tui.setupPanels()
	return tui
}

func (tui *ClappiTUI) setupPanels() {
	sidebarPanels := []struct {
		title, content string
	}{
		{constants.EnvironmentsPanelTitle, "Sidebar top"},
		{constants.ApisPanelTitle, "Sidebar middle"},
		{constants.EndpointsPanelTitle, "Sidebar bottom"},
	}

	mainPanels := []struct {
		title, content string
	}{
		{constants.RequestPanelTitle, "Main top"},
		{constants.ResponsePanelTitle, "Main Bottom"},
	}

	var sideBarViews []tview.Primitive
	for _, p := range sidebarPanels {
		panel := createTextPanelBase(p.title, p.content)
		sideBarViews = append(sideBarViews, panel)
		tui.panels = append(tui.panels, panel)
	}
	var mainViews []tview.Primitive
	for _, p := range mainPanels {
		panel := createTextPanelBase(p.title, p.content)
		mainViews = append(mainViews, panel)
		tui.panels = append(tui.panels, panel)
	}

	sidebar := tui.createVerticalFlex(sideBarViews...)
	mainContent := tui.createVerticalFlex(mainViews...)

	root := tview.NewFlex().
		AddItem(sidebar, 0, 1, false).
		AddItem(mainContent, 0, 2, false)

	tui.app.SetRoot(root, true).EnableMouse(true)
}

func (tui *ClappiTUI) createVerticalFlex(items ...tview.Primitive) *tview.Flex {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	for _, item := range items {
		flex.AddItem(item, 0, 1, false)
	}

	return flex
}
func (tui *ClappiTUI) setFocus(index int) {
	// Reset the border colors
	for _, p := range tui.panels {
		p.(*tview.Box).SetBorderColor(defaultBorderColor)
	}

	// highlight the focused panel
	tui.currentPanel = index
	if tui.currentPanel >= len(tui.panels) {
		tui.currentPanel = 0
	} else if tui.currentPanel < 0 {
		tui.currentPanel = len(tui.panels) - 1
	}
	tui.panels[tui.currentPanel].(*tview.Box).SetBorderColor(focusedBorderColor)
	tui.app.SetFocus(tui.panels[tui.currentPanel])
}
func (tui *ClappiTUI) setupKeyBindings() {
	tui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		// Standard keys, I may bin these tbh
		case tcell.KeyEscape:
			tui.app.Stop()
		case tcell.KeyTab:
			tui.setFocus(tui.currentPanel + 1)
			return nil
		case tcell.KeyBacktab:
			tui.setFocus(tui.currentPanel - 1)
			return nil
		}

		// Vim Keybindings
		switch event.Rune() {
		case 'j':
			if tui.currentPanel < len(tui.panels)-1 {
				tui.setFocus(tui.currentPanel + 1)
			}
			return nil
		case 'k':
			if tui.currentPanel > 0 {
				tui.setFocus(tui.currentPanel - 1)
			}
			return nil
		case 'h':
			if tui.currentPanel >= mainPanelsStartIndex {
				tui.setFocus(tui.currentPanel - mainPanelsStartIndex)
			}
			return nil
		case 'l':
			if tui.currentPanel < mainPanelsStartIndex {
				tui.setFocus(tui.currentPanel + mainPanelsStartIndex)
			}
			return nil
		case 'q':
			tui.app.Stop()
			return nil
		}

		return event
	})
}

func (tui *ClappiTUI) run() error {
	tui.setFocus(0)
	return tui.app.Run()
}

func main() {
	tui := newClappiTUI()
	tui.setupKeyBindings()
	if err := tui.run(); err != nil {
		panic(err)
	}
}
