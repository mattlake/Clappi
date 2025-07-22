package main

import (
	"Clappi/api"
	"Clappi/constants"
	"fmt"
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
	apiManager   *api.APIManager
}

func newClappiTUI(specPath string) *ClappiTUI {
	tui := &ClappiTUI{
		app:          tview.NewApplication(),
		currentPanel: 0,
		sidebarLists: make(map[string]*tview.List),
		apiManager:   api.NewAPIManager(specPath),
	}
	tui.setupPanels()
	return tui
}

func (tui *ClappiTUI) loadAPIs() {
	apiList := tui.sidebarLists[constants.ApisPanelTitle]
	apiList.Clear()

	if err := tui.apiManager.LoadSpecs(); err != nil {
		if mainPanel, ok := tui.panels[mainPanelsStartIndex].(*tview.TextView); ok {
			mainPanel.SetText(fmt.Sprintf("Error loading APIs: %s", err))
		}
		return
	}

	for _, api := range tui.apiManager.GetAPIs() {
		secondaryText := "Valid OpenAPi spec"
		if api.LoadError != nil {
			secondaryText = fmt.Sprintf("Error loading spec: %s", api.LoadError)
		}

		apiList.AddItem(api.Name, secondaryText, 0, func() {
			tui.handleApiSelection(api)
		})
	}
}

func (tui *ClappiTUI) handleApiSelection(api *api.API) {
	if api.LoadError != nil {
		if mainPanel, ok := tui.panels[mainPanelsStartIndex].(*tview.TextView); ok {
			mainPanel.SetText(fmt.Sprintf("Error loading spec: %s", api.LoadError))
		}
		return
	}

	if mainPanel, ok := tui.panels[mainPanelsStartIndex].(*tview.TextView); ok {
		info := api.Model.Model.Info
		details := fmt.Sprintf("API: %s\nVersion: %s\nDescription: %s\nFile: %s",
			info.Title,
			info.Version,
			info.Description,
			api.FilePath)
		mainPanel.SetText(details)
	}
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

func (tui *ClappiTUI) setupPanels() {
	var sideBarViews []tview.Primitive

	for _, config := range sidebarPanelConfigs {
		l := createList(config.title)
		tui.sidebarLists[config.title] = l
		sideBarViews = append(sideBarViews, l)
		tui.panels = append(tui.panels, l)
	}

	tui.loadAPIs()

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
	specPath := "specs" //TODO this will be in user data somewhere someday
	tui := newClappiTUI(specPath)
	tui.setupKeyBindings()
	if err := tui.run(); err != nil {
		panic(err)
	}
}
