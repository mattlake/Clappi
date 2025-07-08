package main

import (
	"Clappi/translations"
	"Clappi/translations/HttpMethods"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/rivo/tview"
	"sort"
)

type vimTreeView struct {
	*tview.TreeView
	lastKey rune
}

type EndpointsBox struct {
	View     tview.Primitive
	TreeView *vimTreeView
}

func createEndpointsBox(paths *v3.Paths) *EndpointsBox {
	tree := &vimTreeView{
		TreeView: tview.NewTreeView(),
		lastKey:  0,
	}
	root := tview.NewTreeNode(translations.Endpoints).
		SetColor(tcell.ColorYellow)
	tree.SetRoot(root)
	tree.SetBorder(true).SetTitle(translations.Endpoints)

	// We sort the names
	pathItems := paths.PathItems.FromOldest()
	pathNames := make([]string, 0)
	for pathName := range pathItems {
		pathNames = append(pathNames, pathName)
	}
	sort.Strings(pathNames)

	// Add the paths to the tree
	for _, pathName := range pathNames {
		pathItem, _ := paths.PathItems.Get(pathName)
		pathNode := tview.NewTreeNode(pathName).SetColor(tcell.ColorWhite).SetExpanded(false)
		root.AddChild(pathNode)

		// Add HTTP methods as children
		if pathItem.Get != nil {
			addMethodNode(pathNode, HttpMethods.Get, pathItem.Get)
		}
		if pathItem.Post != nil {
			addMethodNode(pathNode, HttpMethods.Post, pathItem.Post)
		}
		if pathItem.Put != nil {
			addMethodNode(pathNode, HttpMethods.Put, pathItem.Put)
		}
		if pathItem.Delete != nil {
			addMethodNode(pathNode, HttpMethods.Delete, pathItem.Delete)
		}
		if pathItem.Patch != nil {
			addMethodNode(pathNode, HttpMethods.Patch, pathItem.Patch)
		}
		if pathItem.Options != nil {
			addMethodNode(pathNode, HttpMethods.Options, pathItem.Options)
		}
		if pathItem.Head != nil {
			addMethodNode(pathNode, HttpMethods.Head, pathItem.Head)
		}
		if pathItem.Trace != nil {
			addMethodNode(pathNode, HttpMethods.Trace, pathItem.Trace)
		}
	}

	// We expand the root node by default
	root.SetExpanded(true)

	// Add vim style key bindings
	tree.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		node := tree.GetCurrentNode()
		switch event.Rune() {
		case 'j':
			return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
		case 'k':
			return tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
		case 'h':
			if node != nil {
				node.SetExpanded(false)
			}
			return nil
		case 'l':
			if node != nil {
				node.SetExpanded(true)
			}
			return nil
		case 'g':
			if tree.lastKey == 'g' {
				tree.lastKey = 0
				tree.SetCurrentNode(root)
				return nil
			}
		case 'G':
			if node != nil {
				lastNode := getLastVisibleNode(root)
				tree.SetCurrentNode(lastNode)
			}
			return nil
		}
		// todo add search
		return event
	})

	return &EndpointsBox{
		View:     tree,
		TreeView: tree,
	}
}

func addMethodNode(parent *tview.TreeNode, method string, operation *v3.Operation) {
	var description string
	if operation.Summary != "" {
		description = operation.Summary
	} else if operation.Description != "" {
		description = operation.Description
	}

	methodText := method
	if description != "" {
		methodText = fmt.Sprintf("%s - %s", method, description)
	}

	methodNode := tview.NewTreeNode(methodText)

	switch method {
	case HttpMethods.Get:
		methodNode.SetColor(tcell.ColorGreen)
	case HttpMethods.Post:
		methodNode.SetColor(tcell.ColorYellow)
	case HttpMethods.Put:
		methodNode.SetColor(tcell.ColorBlue)
	case HttpMethods.Delete:
		methodNode.SetColor(tcell.ColorRed)
	case HttpMethods.Patch:
		methodNode.SetColor(tcell.ColorPurple)
	default:
		methodNode.SetColor(tcell.ColorWhite)
	}

	methodNode.SetExpanded(false)
	parent.AddChild(methodNode)

	// Add parameters to the children if they exist
	// todo not sure that I like this info here, maybe it should just be in the main window
	if operation.Parameters != nil {
		paramsNode := tview.NewTreeNode(translations.Parameters).SetColor(tcell.ColorGray)
		methodNode.AddChild(paramsNode)

		for _, param := range operation.Parameters {
			paramText := fmt.Sprintf("%s (%s)", param.Name, param.In)
			if param.Required != nil && *param.Required {
				paramText += translations.RequiredParam
			}
			paramNode := tview.NewTreeNode(paramText).SetColor(tcell.ColorLightGray)
			paramsNode.AddChild(paramNode)
		}
	}
}

func getLastVisibleNode(node *tview.TreeNode) *tview.TreeNode {
	if node == nil {
		return nil
	}

	if !node.IsExpanded() || len(node.GetChildren()) == 0 {
		return node
	}

	children := node.GetChildren()
	return getLastVisibleNode(children[len(children)-1])
}
