package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/rivo/tview"
	"os"
)

var openApiSpec *libopenapi.DocumentModel[v3.Document]

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

	// Load the openApiSpec
	openApiSpec = loadOpenApiSpec()

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
		AddItem(createEndpointsBox(openApiSpec.Model.Paths), 0, 3, false)
}
func createEndpointsBox(paths *v3.Paths) tview.Primitive {
	endpointList := tview.NewList()
	for pathName, _ := range paths.PathItems.FromOldest() {
		endpointList.AddItem(pathName, "", 0, nil)
	}
	endpointList.SetBorder(true).SetTitle("Endpoints")
	return endpointList
}

func createRequestArea() tview.Primitive {
	return tview.NewBox().SetBorder(true).SetTitle("Request")
}

func createResponseArea() tview.Primitive {
	return tview.NewBox().SetBorder(true).SetTitle("Response")
}

func loadOpenApiSpec() *libopenapi.DocumentModel[v3.Document] {
	// Read the api spec from file
	apiSpec, _ := os.ReadFile("assets/petstorev3.json")

	// Load the spec into a document
	doc, err := libopenapi.NewDocument(apiSpec)
	if err != nil {
		panic(fmt.Sprintf("Could not create document: %e", err))
	}

	docModel, errors := doc.BuildV3Model()
	if len(errors) > 0 {
		for i := range errors {
			fmt.Printf("Error: %s\n", errors[i])
		}
		panic(fmt.Sprintf("Could not create v3 model, %d errors reported", len(errors)))
	}

	return docModel
}
