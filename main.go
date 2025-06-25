package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	// Bootstrap the app
	app := tview.NewApplication()

	// Set global keybindings
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'q':
			app.Stop()
		}
		return event
	})

	// Create the main window
	mainWindow := tview.NewFrame(createMainWindow()).
		SetBorders(0, 0, 0, 1, 0, 0).
		AddText("Press q to quit.", false, tview.AlignLeft, tcell.ColorWhite)

	// Run that badboy
	if err := app.SetRoot(mainWindow, true).Run(); err != nil {
		panic(err)
	}
}

func createMainWindow() tview.Primitive {
	return tview.NewFlex().
		AddItem(createLeftSideBar(), 0, 1, false).
		AddItem(createRequestArea(), 0, 2, false).
		AddItem(createResponseArea(), 0, 2, false)

}

func createLeftSideBar() tview.Primitive {
	return tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewBox().SetBorder(true).SetTitle("Environment"), 0, 1, false).
		AddItem(tview.NewBox().SetBorder(true).SetTitle("APIs"), 0, 1, false).
		AddItem(tview.NewBox().SetBorder(true).SetTitle("Endpoints"), 0, 3, false)
}

func createRequestArea() tview.Primitive {
	return tview.NewBox().SetBorder(true).SetTitle("Request")
}

func createResponseArea() tview.Primitive {
	return tview.NewBox().SetBorder(true).SetTitle("Response")
}
