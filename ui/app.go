package ui

import (
	"sort"
	"strings"

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
		Grimoire: grimoire,
		Tree:     buildTree(grimoire),
		ReadOnly: readOnly,
	}

	state.ActiveDoc = findReadmeOrFirst(state)

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
	folderFiles := map[string][]string{}
	rootFiles := []string{}

	for _, doc := range grimoire.Document {
		parts := splitPath(doc.LinkedFile)
		if len(parts) > 1 {
			folder := parts[0]
			folderFiles[folder] = append(folderFiles[folder], doc.LinkedFile)
		} else {
			rootFiles = append(rootFiles, doc.LinkedFile)
		}
	}

	var folders []string
	for f := range folderFiles {
		folders = append(folders, f)
	}
	sort.Strings(folders)

	var nodes []TreeNode
	for _, folder := range folders {
		nodes = append(nodes, TreeNode{
			Name:     folder,
			Path:     folder,
			IsFolder: true,
			Expanded: false,
			Depth:    0,
		})
	}

	for _, file := range rootFiles {
		doc := findDoc(grimoire, file)
		nodes = append(nodes, TreeNode{
			Name:  file,
			Path:  file,
			Depth: 0,
			Doc:   doc,
		})
	}

	return nodes
}

func rebuildVisibleTree(state *AppState) []TreeNode {
	folderFiles := map[string][]string{}
	rootFiles := []string{}

	for _, doc := range state.Grimoire.Document {
		parts := splitPath(doc.LinkedFile)
		if len(parts) > 1 {
			folder := parts[0]
			folderFiles[folder] = append(folderFiles[folder], doc.LinkedFile)
		} else {
			rootFiles = append(rootFiles, doc.LinkedFile)
		}
	}

	var folders []string
	for f := range folderFiles {
		folders = append(folders, f)
	}
	sort.Strings(folders)

	expandedState := map[string]bool{}
	for _, node := range state.Tree {
		if node.IsFolder {
			expandedState[node.Name] = node.Expanded
		}
	}

	var nodes []TreeNode
	for _, folder := range folders {
		isExpanded := expandedState[folder]
		nodes = append(nodes, TreeNode{
			Name:     folder,
			Path:     folder,
			IsFolder: true,
			Expanded: isExpanded,
			Depth:    0,
		})

		if isExpanded {
			files := folderFiles[folder]
			sort.Strings(files)
			for _, file := range files {
				doc := findDoc(state.Grimoire, file)
				parts := splitPath(file)
				name := parts[len(parts)-1]
				nodes = append(nodes, TreeNode{
					Name:  name,
					Path:  file,
					Depth: 1,
					Doc:   doc,
				})
			}
		}
	}

	for _, file := range rootFiles {
		doc := findDoc(state.Grimoire, file)
		nodes = append(nodes, TreeNode{
			Name:  file,
			Path:  file,
			Depth: 0,
			Doc:   doc,
		})
	}

	return nodes
}

func findDoc(grimoire *core.Grimoire, path string) *core.Doc {
	for i := range grimoire.Document {
		if grimoire.Document[i].LinkedFile == path {
			return &grimoire.Document[i]
		}
	}
	return nil
}

func findReadmeOrFirst(state *AppState) *core.Doc {
	for i := range state.Grimoire.Document {
		name := strings.ToLower(state.Grimoire.Document[i].LinkedFile)
		if name == "readme.md" {
			return &state.Grimoire.Document[i]
		}
	}
	if len(state.Grimoire.Document) > 0 {
		return &state.Grimoire.Document[0]
	}
	return nil
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
