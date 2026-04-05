package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/theoneandonlyvabo/grimoire/core"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case tea.KeyMsg:
		if m.EditMode {
			return m.handleEditKey(msg)
		}
		return m.handleNavKey(msg)
	}

	return m, nil
}

func (m Model) handleNavKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "Q":
		if m.Pane == PaneSidebar {
			if m.Dirty {
				// TODO: confirm quit dialog
				return m, tea.Quit
			}
			return m, tea.Quit
		}

	case "ctrl+c":
		return m, tea.Quit

	case "ctrl+s":
		if !m.ReadOnly {
			core.Save(m.Grimoire)
			m.Dirty = false
		}

	case "tab":
		m.Pane = PaneSidebar

	case "left":
		m.Pane = PaneSidebar

	case "right":
		if m.Pane == PaneSidebar {
			node := m.Tree[m.ActiveIndex]
			if node.IsFolder {
				if !node.Expanded {
					m.Tree[m.ActiveIndex].Expanded = true
					m.Tree = rebuildVisibleTree(m.Grimoire, m.Tree)
				}
			} else {
				m.Pane = PaneEditor
				m.ActiveField = 0
			}
		}

	case "up":
		if m.Pane == PaneSidebar {
			if m.ActiveIndex > 0 {
				m.ActiveIndex--
			}
		} else {
			if m.ActiveField > 0 {
				m.ActiveField--
			}
		}

	case "down":
		if m.Pane == PaneSidebar {
			if m.ActiveIndex < len(m.Tree)-1 {
				m.ActiveIndex++
			}
		} else {
			doc := m.ActiveDoc()
			if doc != nil && m.ActiveField < len(doc.Functions) {
				m.ActiveField++
			}
		}

	case "enter":
		if m.Pane == PaneSidebar {
			node := m.Tree[m.ActiveIndex]
			if node.IsFolder {
				m.Tree[m.ActiveIndex].Expanded = !m.Tree[m.ActiveIndex].Expanded
				m.Tree = rebuildVisibleTree(m.Grimoire, m.Tree)
			} else {
				m.Pane = PaneEditor
				m.ActiveField = 0
			}
		} else {
			if !m.ReadOnly {
				m = m.enterEditMode()
			}
		}

	case "esc":
		if m.Pane == PaneEditor {
			m.Pane = PaneSidebar
		}
	}

	return m, nil
}

func (m Model) handleEditKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.EditMode = false
		m.EditBuffer = ""
		m.EditTarget = ""

	case "enter":
		m = m.commitEdit()

	case "backspace":
		if len(m.EditBuffer) > 0 {
			runes := []rune(m.EditBuffer)
			m.EditBuffer = string(runes[:len(runes)-1])
		}

	default:
		if msg.Type == tea.KeyRunes {
			m.EditBuffer += msg.String()
		}
	}

	return m, nil
}

func (m Model) enterEditMode() Model {
	doc := m.ActiveDoc()
	if doc == nil {
		return m
	}

	if m.ActiveField == 0 {
		m.EditTarget = "description"
		m.EditBuffer = doc.Description
	} else {
		fnIndex := m.ActiveField - 1
		if fnIndex < len(doc.Functions) {
			m.EditTarget = EditTarget("fn:" + doc.Functions[fnIndex].Name)
			m.EditBuffer = doc.Functions[fnIndex].Notes
		}
	}

	m.EditMode = true
	return m
}

func (m Model) commitEdit() Model {
	doc := m.ActiveDoc()
	if doc == nil {
		return m
	}

	author := core.GetUserName()
	now := time.Now().Format("2006-01-02 15:04")

	if m.EditTarget == "description" {
		doc.Description = m.EditBuffer
		doc.Author = author
		doc.UpdatedAt = now
	} else if len(m.EditTarget) > 3 && m.EditTarget[:3] == "fn:" {
		fnName := string(m.EditTarget[3:])
		for i := range doc.Functions {
			if doc.Functions[i].Name == fnName {
				doc.Functions[i].Notes = m.EditBuffer
				doc.Functions[i].Author = author
				doc.Functions[i].UpdatedAt = now
				break
			}
		}
	}

	m.EditMode = false
	m.EditBuffer = ""
	m.EditTarget = ""
	m.Dirty = true
	return m
}
