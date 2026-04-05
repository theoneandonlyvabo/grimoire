package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/theoneandonlyvabo/grimoire/core"
)

type ActivePane int
type EditTarget string

const (
	PaneSidebar ActivePane = iota
	PaneEditor
)

type Model struct {
	// data utama
	Grimoire *core.Grimoire
	LiveMeta core.MetaData

	// sidebar
	Tree        []TreeNode
	ActiveIndex int
	Pane        ActivePane

	// editor
	ActiveField int

	// edit mode
	EditMode   bool
	EditBuffer string
	EditTarget EditTarget

	// ui
	Width    int
	Height   int
	ReadOnly bool
	Dirty    bool
}

func NewModel(grimoire *core.Grimoire, readOnly bool) Model {
	liveMeta, _ := core.GetMetadata()
	tree := buildCollapsedTree(grimoire)

	m := Model{
		Grimoire: grimoire,
		LiveMeta: liveMeta,
		Tree:     tree,
		ReadOnly: readOnly,
	}

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) ActiveDoc() *core.Doc {
	if m.ActiveIndex >= len(m.Tree) {
		return nil
	}
	node := m.Tree[m.ActiveIndex]
	if node.IsFolder || node.Doc == nil {
		return nil
	}
	return node.Doc
}
