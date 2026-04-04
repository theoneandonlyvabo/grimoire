package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/theoneandonlyvabo/grimoire/core"
)

func handleKey(screen tcell.Screen, state *AppState, ev *tcell.EventKey) bool {
	switch ev.Key() {

	case tcell.KeyEscape:
		if state.ActivePane == 1 {
			state.ActivePane = 0
			return false
		}
		if state.Dirty {
			return confirmQuit(screen, state)
		}
		return true

	case tcell.KeyTab:
		state.ActivePane = 0

	case tcell.KeyLeft:
		state.ActivePane = 0

	case tcell.KeyRight:
		if state.ActivePane == 0 {
			node := state.Tree[state.ActiveIndex]
			if node.IsFolder {
				if !node.Expanded {
					state.Tree[state.ActiveIndex].Expanded = true
					state.Tree = rebuildVisibleTree(state)
				}
			} else {
				updateActiveDoc(state)
				state.ActivePane = 1
			}
		}

	case tcell.KeyUp:
		if state.ActivePane == 0 {
			if state.ActiveIndex > 0 {
				state.ActiveIndex--
			}
		} else {
			if state.ActiveField > 0 {
				state.ActiveField--
			}
		}

	case tcell.KeyDown:
		if state.ActivePane == 0 {
			if state.ActiveIndex < len(state.Tree)-1 {
				state.ActiveIndex++
			}
		} else {
			if state.ActiveField < 1 {
				state.ActiveField++
			}
		}

	case tcell.KeyEnter:
		if state.ActivePane == 0 {
			node := state.Tree[state.ActiveIndex]
			if node.IsFolder {
				state.Tree[state.ActiveIndex].Expanded = !state.Tree[state.ActiveIndex].Expanded
				state.Tree = rebuildVisibleTree(state)
			} else {
				updateActiveDoc(state)
				state.ActivePane = 1
			}
		}

	case tcell.KeyCtrlS:
		if !state.ReadOnly {
			core.Save(state.Grimoire)
			state.Dirty = false
		}

	case tcell.KeyRune:
		switch ev.Rune() {
		case 'q', 'Q':
			if state.ActivePane == 0 {
				if state.Dirty {
					return confirmQuit(screen, state)
				}
				return true
			}
		}
	}

	return false
}

func updateActiveDoc(state *AppState) {
	if state.ActiveIndex >= len(state.Tree) {
		return
	}
	node := state.Tree[state.ActiveIndex]
	if !node.IsFolder && node.Doc != nil {
		state.ActiveDoc = node.Doc
	}
}

func confirmQuit(screen tcell.Screen, state *AppState) bool {
	w, h := screen.Size()
	msg := "  unsaved changes. quit anyway? (y/n)  "
	x := w/2 - len(msg)/2
	y := h / 2
	fill(screen, x-2, y-2, x+len(msg)+2, y+2, stDefault)
	drawBox(screen, x-2, y-2, x+len(msg)+2, y+2, stBorder)
	draw(screen, x, y, stDefault, msg)
	screen.Show()

	for {
		ev := screen.PollEvent()
		if key, ok := ev.(*tcell.EventKey); ok {
			if key.Rune() == 'y' || key.Rune() == 'Y' {
				return true
			}
			if key.Rune() == 'n' || key.Rune() == 'N' {
				return false
			}
		}
	}
}
