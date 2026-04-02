package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/theoneandonlyvabo/grimoire/core"
)

type TreeNode struct {
	Name     string
	Path     string
	IsFolder bool
	Expanded bool
	Depth    int
	Doc      *core.Doc
}

type AppState struct {
	Grimoire    *core.Grimoire
	Tree        []TreeNode
	ActiveIndex int
	ActivePane  int
	ActiveField int
	ActiveDoc   *core.Doc
	ReadOnly    bool
	Dirty       bool
}

func Start(grimoire *core.Grimoire, readOnly bool) error {
	screen, err := tcell.NewScreen()
	if err != nil {
		return err
	}
	if err := screen.Init(); err != nil {
		return err
	}
	defer screen.Fini()

	state := &AppState{
		Grimoire:    grimoire,
		Tree:        buildTree(grimoire),
		ActiveIndex: 0,
		ActivePane:  0,
		ActiveField: 0,
		ReadOnly:    readOnly,
		Dirty:       false,
	}

	state.ActiveDoc = findFirstDoc(state)

	for {
		screen.Clear()
		render(screen, state)
		screen.Show()

		event := screen.PollEvent()
		switch ev := event.(type) {
		case *tcell.EventResize:
			screen.Sync()
		case *tcell.EventKey:
			if done := handleKey(screen, state, ev); done {
				return nil
			}
		}
	}
}

func buildTree(grimoire *core.Grimoire) []TreeNode {
	folderMap := map[string]bool{}
	var nodes []TreeNode

	for _, doc := range grimoire.Document {
		parts := splitPath(doc.LinkedFile)
		if len(parts) > 1 {
			folder := parts[0]
			if !folderMap[folder] {
				folderMap[folder] = true
				nodes = append(nodes, TreeNode{
					Name:     folder,
					Path:     folder,
					IsFolder: true,
					Expanded: false,
					Depth:    0,
				})
			}
		}
	}

	return nodes
}

func splitPath(path string) []string {
	var parts []string
	current := ""
	for _, ch := range path {
		if ch == '/' || ch == '\\' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(ch)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}

func findFirstDoc(state *AppState) *core.Doc {
	if len(state.Grimoire.Document) > 0 {
		return &state.Grimoire.Document[0]
	}
	return nil
}
