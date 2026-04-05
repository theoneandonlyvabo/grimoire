package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	colAmber     = lipgloss.Color("#F0C070")
	colAmberBase = lipgloss.Color("#C9A96E")
	colTextPri   = lipgloss.Color("#E8E4D8")
	colTextMuted = lipgloss.Color("#646460")
	colTextDim   = lipgloss.Color("#323230")
	colBlue      = lipgloss.Color("#6E9EC9")
	colGreen     = lipgloss.Color("#6EC98A")
	colRed       = lipgloss.Color("#C96E6E")
	colBorder    = lipgloss.Color("#C9A96E")

	styleBase = lipgloss.NewStyle().
			Foreground(colTextPri)

	styleDim = lipgloss.NewStyle().
			Foreground(colTextMuted)

	styleVeryDim = lipgloss.NewStyle().
			Foreground(colTextDim)

	styleAmber = lipgloss.NewStyle().
			Foreground(colAmber)

	styleBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colBorder)

	styleBorderDim = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colTextDim)

	styleBorderActive = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colAmber)
)

func (m Model) View() string {
	if m.Width == 0 {
		return ""
	}

	header := m.viewHeader()
	body := m.viewBody()
	footer := m.viewFooter()

	outer := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colBorder).
		Width(m.Width-2).
		Padding(0, 1)

	content := lipgloss.JoinVertical(lipgloss.Left,
		header,
		strings.Repeat("─", m.Width-4),
		body,
		strings.Repeat("─", m.Width-4),
		footer,
	)

	return outer.Render(content)
}

func (m Model) viewHeader() string {
	meta := m.LiveMeta

	repo := meta.Repository
	repo = strings.TrimPrefix(repo, "https://github.com/")
	repo = strings.TrimPrefix(repo, "git@github.com:")
	repo = strings.TrimSuffix(repo, ".git")

	date := meta.LastCommitDate
	if len(date) > 10 {
		date = date[:10]
	}

	brand := styleAmber.Bold(true).Render("GRIMOIRE")
	branch := styleDim.Render("⎇ " + meta.Branch)
	repoStr := styleDim.Render(repo)
	commit := lipgloss.NewStyle().Foreground(colAmberBase).Render(
		fmt.Sprintf("Commit ↑ %s · %s", meta.LastCommit, date),
	)

	row1 := lipgloss.JoinHorizontal(lipgloss.Left,
		brand, "  ", branch, "  │  ", repoStr, "  │  ", commit,
	)

	msg := styleVeryDim.Render("»  " + meta.LastCommitMessage)

	return lipgloss.JoinVertical(lipgloss.Left, row1, msg)
}

func (m Model) viewBody() string {
	sidebarWidth := 24
	editorWidth := m.Width - sidebarWidth - 6

	sidebar := m.viewSidebar(sidebarWidth)
	editor := m.viewEditor(editorWidth)

	divider := lipgloss.NewStyle().
		Foreground(colBorder).
		Render(strings.Repeat("│\n", m.Height-8))

	return lipgloss.JoinHorizontal(lipgloss.Top,
		sidebar, divider, editor,
	)
}

func (m Model) viewSidebar(width int) string {
	var lines []string

	title := styleDim.Render("files")
	lines = append(lines, title)
	lines = append(lines, styleVeryDim.Render(strings.Repeat("─", width-2)))

	for i, node := range m.Tree {
		indent := strings.Repeat("  ", node.Depth)
		isFocused := i == m.ActiveIndex && m.Pane == PaneSidebar

		var line string
		if node.IsFolder {
			arrow := "▸ "
			if node.Expanded {
				arrow = "▾ "
			}
			text := indent + arrow + node.Name + "/"
			if isFocused {
				line = styleAmber.Render("▌") + styleAmber.Render(text)
			} else if node.Expanded {
				line = " " + lipgloss.NewStyle().Foreground(colAmberBase).Render(text)
			} else {
				line = " " + styleDim.Render(text)
			}
		} else {
			dotColor := colTextDim
			if node.Doc != nil {
				dotColor = statusColor(node.Doc.Status)
			}
			dot := lipgloss.NewStyle().Foreground(dotColor).Render("●")

			maxLen := width - 4 - node.Depth*2
			name := node.Name
			if len([]rune(name)) > maxLen {
				name = string([]rune(name)[:maxLen-1]) + "…"
			}

			var nameStr string
			if isFocused {
				nameStr = styleAmber.Render(name)
				line = styleAmber.Render("▌") + indent + dot + " " + nameStr
			} else {
				nameStr = styleDim.Render(name)
				line = " " + indent + dot + " " + nameStr
			}
		}

		lines = append(lines, line)
	}

	return lipgloss.NewStyle().
		Width(width).
		Padding(1, 1).
		Render(strings.Join(lines, "\n"))
}

func (m Model) viewEditor(width int) string {
	doc := m.ActiveDoc()

	if doc == nil {
		return lipgloss.NewStyle().
			Width(width).
			Padding(2, 3).
			Render(styleVeryDim.Render("select a file from the sidebar"))
	}

	var sections []string

	filename := styleBase.Bold(true).Render(doc.LinkedFile)
	meta := styleVeryDim.Render(fmt.Sprintf("author %s  ·  %s", doc.Author, doc.UpdatedAt))
	sections = append(sections, filename+"\n"+meta)
	sections = append(sections, styleVeryDim.Render(strings.Repeat("─", width-6)))

	descSection := m.viewField(
		"description",
		doc.Description,
		"no description yet...",
		m.Pane == PaneEditor && m.ActiveField == 0,
		m.EditMode && m.EditTarget == "description",
		width-6,
	)
	sections = append(sections, descSection)

	if len(doc.Functions) > 0 {
		sections = append(sections, styleVeryDim.Render(strings.Repeat("─", width-6)))
		sections = append(sections, styleDim.Render("functions"))

		for fnIdx, fn := range doc.Functions {
			isPrivate := len(fn.Name) > 0 && fn.Name[0] >= 'a' && fn.Name[0] <= 'z'
			fnColor := colBlue
			if isPrivate {
				fnColor = colTextMuted
			}

			maxSig := width - 6 - len(fn.Name) - 3
			sig := fn.Signature
			if len([]rune(sig)) > maxSig {
				sig = string([]rune(sig)[:maxSig]) + "…"
			}

			fnName := lipgloss.NewStyle().Foreground(fnColor).Render(fn.Name)
			fnSig := styleVeryDim.Render("  " + sig)
			sections = append(sections, fnName+fnSig)

			fieldIndex := fnIdx + 1
			notesSection := m.viewField(
				"notes",
				fn.Notes,
				"no notes yet...",
				m.Pane == PaneEditor && m.ActiveField == fieldIndex,
				m.EditMode && m.EditTarget == EditTarget("fn:"+fn.Name),
				width-6,
			)
			sections = append(sections, notesSection)
		}
	}

	return lipgloss.NewStyle().
		Width(width).
		Padding(1, 3).
		Render(strings.Join(sections, "\n\n"))
}

func (m Model) viewField(label, content, placeholder string, active, editing bool, width int) string {
	var boxStyle lipgloss.Style
	if editing {
		boxStyle = styleBorderActive
	} else if active {
		boxStyle = styleBorder
	} else {
		boxStyle = styleBorderDim
	}

	var inner string
	if editing {
		inner = styleBase.Render(m.EditBuffer + "█")
	} else if content != "" {
		inner = styleBase.Render(wrapString(content, width-4))
	} else {
		inner = styleVeryDim.Render(placeholder)
	}

	labelStr := styleDim.Render(label)
	if active || editing {
		labelStr = styleAmber.Render(label)
	}

	return labelStr + "\n" + boxStyle.Width(width-2).Padding(0, 1).Render(inner)
}

func statusColor(status string) lipgloss.Color {
	switch status {
	case "stable":
		return colGreen
	case "deprecated":
		return colRed
	default:
		return colAmberBase
	}
}

func (m Model) viewFooter() string {
	type bind struct{ key, desc string }

	binds := []bind{
		{"tab", "sidebar"},
		{"← →", "switch pane"},
		{"↑↓", "navigate"},
		{"enter", "open / edit"},
		{"ctrl+s", "save"},
		{"q", "quit"},
	}
	if m.ReadOnly {
		binds = []bind{
			{"← →", "switch pane"},
			{"↑↓", "navigate"},
			{"enter", "open"},
			{"q", "quit"},
		}
	}

	var parts []string
	for _, b := range binds {
		key := lipgloss.NewStyle().Foreground(colAmberBase).Render(b.key)
		desc := styleVeryDim.Render(" " + b.desc)
		parts = append(parts, key+desc)
	}

	left := strings.Join(parts, "  ")

	right := ""
	if m.Dirty {
		right = styleAmber.Render("● unsaved")
	}
	if m.ReadOnly {
		right = styleVeryDim.Render("read-only")
	}

	gap := m.Width - lipgloss.Width(left) - lipgloss.Width(right) - 6
	if gap < 0 {
		gap = 0
	}

	return left + strings.Repeat(" ", gap) + right
}

func wrapString(text string, width int) string {
	if width <= 0 {
		return text
	}
	words := strings.Fields(text)
	var lines []string
	current := ""
	for _, word := range words {
		if lipgloss.Width(current)+lipgloss.Width(word)+1 > width {
			if current != "" {
				lines = append(lines, current)
			}
			current = word
		} else {
			if current == "" {
				current = word
			} else {
				current += " " + word
			}
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return strings.Join(lines, "\n")
}
