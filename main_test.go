package main

import (
	"Clappi/api"
	"Clappi/constants"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"testing"
)

func TestClappiTUI_setupPanels(t *testing.T) {
	tui := newClappiTUI("testdata")

	// test panel creation
	if len(tui.panels) != 5 {
		t.Errorf("Expected 5 panels, got %d", len(tui.panels))
	}

	// Test Sidebar list creation
	expectedLists := []string{constants.EnvironmentsPanelTitle, constants.ApisPanelTitle, constants.EndpointsPanelTitle}
	for _, title := range expectedLists {
		if _, ok := tui.sidebarLists[title]; !ok {
			t.Errorf("Expected sidebar list %s, got none", title)
		}
	}
}

func TestClappiTUI_handleApiSelection(t *testing.T) {

	t.Run("Test API with load error", func(t *testing.T) {
		tui := newClappiTUI("testdata")
		testAPI := &api.API{
			Name:     "Test API",
			FilePath: "testdata/valid_spec.yaml",
		}

		// Test API with load error
		testAPI.LoadError = fmt.Errorf("test error")
		tui.handleApiSelection(testAPI)

		mainPanel, ok := tui.panels[mainPanelsStartIndex].(*tview.TextView)
		if !ok {
			t.Error("Failed to get main panel")
			return
		}

		if mainPanel.GetText(false) == "" {
			t.Error("Expected error message in main panel")
		}
	})
}
func TestClappiTUI_handleApiSelection_Endpoints(t *testing.T) {
	tui := newClappiTUI("testdata")

	for _, test := range tui.apiManager.GetAPIs() {
		// We only want to test the valid spec here
		if test.FilePath != "testdata/valid_spec.yaml" {
			return
		}
	}

	endpointsList := tui.sidebarLists[constants.EndpointsPanelTitle]
	if endpointsList.GetItemCount() == 0 {
		t.Error("Expected endpoints list to have items")
	}
}

func TestClappiTUI_setFocus(t *testing.T) {
	tui := newClappiTUI("testdata")

	tests := []struct {
		name          string
		index         int
		wantedColor   tcell.Color
		unwantedColor tcell.Color
	}{
		{
			name:          "Focus first panel",
			index:         0,
			wantedColor:   focusedBorderColor,
			unwantedColor: defaultBorderColor,
		},
		{
			name:          "Focus last panel",
			index:         len(tui.panels) - 1,
			wantedColor:   focusedBorderColor,
			unwantedColor: defaultBorderColor,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tui.setFocus(tt.index)

			if tui.panels[tt.index].GetBorderColor() != tt.wantedColor {
				t.Errorf("Expected border color %v, got %v", tt.wantedColor, tui.panels[tt.index].GetBorderColor())
			}

			// Check other panels have the default color
			for i, panel := range tui.panels {
				if i != tt.index && panel.GetBorderColor() != tt.unwantedColor {
					t.Errorf("Panel %d should have color %v, got %v", i, tt.unwantedColor, panel.GetBorderColor())
				}
			}
		})
	}
}
