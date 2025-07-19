package main

import (
	"Clappi/constants"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Panel interface {
	tview.Primitive
	SetBorderColor(color tcell.Color) *tview.Box
}

type PanelConfig struct {
	title   string
	content string
}

type ClappiTUI struct {
	app          *tview.Application
	panels       []Panel
	currentPanel int
	sidebarLists map[string]*tview.List
}

const (
	mainPanelsStartIndex = 3
	defaultBorderColor   = tcell.ColorWhite
	focusedBorderColor   = tcell.ColorYellow
)

// Define panel configurations
var (
	sidebarPanelConfigs = []PanelConfig{
		{title: constants.EnvironmentsPanelTitle, content: ""},
		{title: constants.ApisPanelTitle, content: ""},
		{title: constants.EndpointsPanelTitle, content: ""},
	}

	mainPanelConfigs = []PanelConfig{
		{title: constants.RequestPanelTitle, content: "Main top"},
		{title: constants.ResponsePanelTitle, content: "Main Bottom"},
	}
)

func createList(title string) *tview.List {
	// This needs to be done is stages otherwise the list is converted to a box
	list := tview.NewList()
	list.SetTitle(title).SetBorder(true)
	return list
}

func createTextPanel(config PanelConfig) *tview.TextView {
	// This needs to be done is stages otherwise the list is converted to a box
	tv := tview.NewTextView()
	tv.SetText(config.content).
		SetTextAlign(tview.AlignCenter).
		SetBorder(true).
		SetTitle(config.title)
	return tv
}

func newClappiTUI() *ClappiTUI {
	tui := &ClappiTUI{
		app:          tview.NewApplication(),
		currentPanel: 0,
		sidebarLists: make(map[string]*tview.List),
	}
	tui.setupPanels()
	return tui
}

func (tui *ClappiTUI) setupPanels() {
	var sideBarViews []tview.Primitive

	for _, config := range sidebarPanelConfigs {
		l := createList(config.title)
		tui.sidebarLists[config.title] = l
		sideBarViews = append(sideBarViews, l)
		tui.panels = append(tui.panels, l)
	}

	tui.hydratePanels()

	var mainViews []tview.Primitive
	for _, config := range mainPanelConfigs {
		panel := createTextPanel(config)
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

func (tui *ClappiTUI) hydratePanels() {
	// todo replace this later with dynamic functionality
	environmentslist := tui.sidebarLists[constants.EnvironmentsPanelTitle]
	environmentslist.AddItem("Default", "", 0, nil).
		AddItem("Local", "", 0, nil).
		AddItem("Development", "", 0, nil).
		AddItem("Production", "", 0, nil)

	apiList := tui.sidebarLists[constants.ApisPanelTitle]
	apiList.AddItem("Official Joke API", "", 0, nil)

	endpointsList := tui.sidebarLists[constants.EndpointsPanelTitle]
	endpointsList.AddItem("https://official-joke-api.appspot.com/jokes/random", "", 0, nil).
		AddItem("https://official-joke-api.appspot.com/types", "", 0, nil).
		AddItem("https://official-joke-api.appspot.com/jokes/ten", "", 0, nil).
		AddItem("https://official-joke-api.appspot.com/jokes/programming/random", "", 0, nil).
		AddItem("https://official-joke-api.appspot.com/jokes/:id", "", 0, nil)
}

func (tui *ClappiTUI) createVerticalFlex(items ...tview.Primitive) *tview.Flex {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	for _, item := range items {
		flex.AddItem(item, 0, 1, false)
	}

	return flex
}
func (tui *ClappiTUI) setFocus(index int) {
	tui.resetBorderColors()
	tui.updateCurrentPanel(index)
	tui.highlightCurrentPanel()
	tui.app.SetFocus(tui.panels[tui.currentPanel])
}

func (tui *ClappiTUI) highlightCurrentPanel() {
	tui.panels[tui.currentPanel].SetBorderColor(focusedBorderColor)
}

func (tui *ClappiTUI) updateCurrentPanel(index int) {
	tui.currentPanel = index
	if tui.currentPanel >= len(tui.panels) {
		tui.currentPanel = 0
	} else if tui.currentPanel < 0 {
		tui.currentPanel = len(tui.panels) - 1
	}
}

func (tui *ClappiTUI) resetBorderColors() {
	for _, p := range tui.panels {
		p.SetBorderColor(defaultBorderColor)
	}
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
