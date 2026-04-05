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

	styleBase  = lipgloss.NewStyle().Foreground(colTextPri)
	styleDim   = lipgloss.NewStyle().Foreground(colTextMuted)
	styleFaint = lipgloss.NewStyle().Foreground(colTextDim)
	styleAmber = lipgloss.NewStyle().Foreground(colAmber)

	styleBorderNormal = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(colTextDim)

	styleBorderActive = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(colAmber)
)

func (m Model) View() string {
	if m.Width == 0 {
		return ""
	}

	if m.Width < 80 || m.Height < 24 {
		return fmt.Sprintf(
			"\n  Terminal too small.\n  Minimum 80×24 required.\n  Current: %d×%d\n",
			m.Width, m.Height,
		)
	}

	header := m.viewHeader()
	body := m.viewBody()
	footer := m.viewFooter()

	outerBorder := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(colBorder).
		Width(m.Width - 2)

	inner := lipgloss.JoinVertical(lipgloss.Left,
		header,
		lipgloss.NewStyle().Foreground(colBorder).Render(strings.Repeat("─", m.Width-4)),
		body,
		lipgloss.NewStyle().Foreground(colBorder).Render(strings.Repeat("─", m.Width-4)),
		footer,
	)

	return outerBorder.Render(inner)
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
	branch := lipgloss.NewStyle().Foreground(colTextMuted).Render("  ⎇ " + meta.Branch)
	repoStr := styleFaint.Render("  " + repo)

	row1 := brand + branch + repoStr

	commitHash := lipgloss.NewStyle().Foreground(colAmberBase).Render("↑ " + meta.LastCommit)
	commitDate := styleFaint.Render(" · " + date)
	commitMsg := styleFaint.Render("  »  " + meta.LastCommitMessage)

	row2 := commitHash + commitDate + commitMsg

	return lipgloss.NewStyle().Padding(0, 1).Render(
		lipgloss.JoinVertical(lipgloss.Left, row1, row2),
	)
}

func (m Model) viewBody() string {
	sidebarWidth := 26
	editorWidth := m.Width - sidebarWidth - 5

	bodyHeight := m.Height - 8

	sidebar := m.viewSidebar(sidebarWidth, bodyHeight)
	divider := lipgloss.NewStyle().
		Foreground(colBorder).
		Render(strings.Repeat("│\n", bodyHeight))
	editor := m.viewEditor(editorWidth, bodyHeight)

	return lipgloss.JoinHorizontal(lipgloss.Top, sidebar, divider, editor)
}

func (m Model) viewSidebar(width, height int) string {
	var lines []string

	title := styleFaint.Render("files")
	lines = append(lines, title)
	lines = append(lines,
		lipgloss.NewStyle().Foreground(colTextDim).Render(strings.Repeat("─", width-3)),
	)

	visibleStart := 0
	visibleEnd := height - 4

	if m.ActiveIndex >= visibleEnd+visibleStart {
		visibleStart = m.ActiveIndex - visibleEnd + 1
	}

	visibleNodes := m.Tree
	if visibleStart > 0 {
		visibleNodes = m.Tree[visibleStart:]
	}

	shown := 0
	for i, node := range visibleNodes {
		if shown >= height-4 {
			break
		}

		actualIndex := i + visibleStart
		indent := strings.Repeat("  ", node.Depth)
		isFocused := actualIndex == m.ActiveIndex && m.Pane == PaneSidebar

		var line string
		if node.IsFolder {
			arrow := "▸ "
			if node.Expanded {
				arrow = "▾ "
			}
			text := indent + arrow + node.Name + "/"
			if isFocused {
				line = styleAmber.Render("▌ ") + styleAmber.Render(text)
			} else if node.Expanded {
				line = "  " + lipgloss.NewStyle().Foreground(colAmberBase).Render(text)
			} else {
				line = "  " + styleDim.Render(text)
			}
		} else {
			dotColor := colTextDim
			if node.Doc != nil {
				dotColor = statusColor(node.Doc.Status)
			}
			dot := lipgloss.NewStyle().Foreground(dotColor).Render("●")

			maxLen := width - 6 - node.Depth*2
			name := node.Name
			if len([]rune(name)) > maxLen {
				name = string([]rune(name)[:maxLen-1]) + "…"
			}

			if isFocused {
				line = styleAmber.Render("▌ ") + indent + dot + " " + styleAmber.Render(name)
			} else {
				line = "  " + indent + dot + " " + styleDim.Render(name)
			}
		}

		lines = append(lines, line)
		shown++
	}

	hasMore := visibleStart+shown < len(m.Tree)
	hasAbove := visibleStart > 0

	scrollIndicator := ""
	if hasAbove && hasMore {
		scrollIndicator = styleFaint.Render("↑↓")
	} else if hasAbove {
		scrollIndicator = styleFaint.Render("↑")
	} else if hasMore {
		scrollIndicator = styleFaint.Render("↓")
	}

	if scrollIndicator != "" {
		lines = append(lines, scrollIndicator)
	}

	return lipgloss.NewStyle().
		Width(width).
		Padding(0, 1).
		Render(strings.Join(lines, "\n"))
}

func (m Model) viewEditor(width, height int) string {
	doc := m.ActiveDoc()

	if doc == nil {
		return lipgloss.NewStyle().
			Width(width).
			Padding(1, 2).
			Render(styleFaint.Render("select a file from the sidebar"))
	}

	var sections []string

	filename := styleBase.Bold(true).Render(doc.LinkedFile)
	meta := styleFaint.Render(fmt.Sprintf("author %s  ·  %s", doc.Author, doc.UpdatedAt))
	sections = append(sections, filename)
	sections = append(sections, meta)
	sections = append(sections,
		lipgloss.NewStyle().Foreground(colTextDim).Render(strings.Repeat("─", width-4)),
	)

	isDescActive := m.Pane == PaneEditor && m.ActiveField == 0
	isDescEditing := m.EditMode && m.EditTarget == "description"
	sections = append(sections, m.viewField(
		"description",
		doc.Description,
		"no description yet...",
		isDescActive,
		isDescEditing,
		width-4,
	))

	if len(doc.Functions) > 0 {
		sections = append(sections, "")
		sections = append(sections,
			lipgloss.NewStyle().Foreground(colTextDim).Render(strings.Repeat("─", width-4)),
		)
		sections = append(sections, styleFaint.Render("functions"))
		sections = append(sections, "")

		for fnIdx, fn := range doc.Functions {
			isPrivate := len(fn.Name) > 0 && fn.Name[0] >= 'a' && fn.Name[0] <= 'z'
			fnColor := colBlue
			if isPrivate {
				fnColor = colTextMuted
			}

			maxSig := width - 4 - len(fn.Name) - 3
			sig := fn.Signature
			if len([]rune(sig)) > maxSig {
				sig = string([]rune(sig)[:maxSig]) + "…"
			}

			fnName := lipgloss.NewStyle().Foreground(fnColor).Render(fn.Name)
			fnSig := styleFaint.Render("  " + sig)
			sections = append(sections, fnName+fnSig)

			fieldIndex := fnIdx + 1
			isFnActive := m.Pane == PaneEditor && m.ActiveField == fieldIndex
			isFnEditing := m.EditMode && m.EditTarget == EditTarget("fn:"+fn.Name)

			sections = append(sections, m.viewField(
				"notes",
				fn.Notes,
				"no notes yet...",
				isFnActive,
				isFnEditing,
				width-4,
			))
			sections = append(sections, "")
		}
	}

	content := strings.Join(sections, "\n")
	lines := strings.Split(content, "\n")

	visibleStart := 0
	visibleEnd := height - 2

	if m.ActiveField > 0 {
		linesPerField := 4
		targetLine := m.ActiveField * linesPerField
		if targetLine > visibleEnd {
			visibleStart = targetLine - visibleEnd + 2
		}
	}

	hasAbove := visibleStart > 0
	hasMore := visibleStart+visibleEnd < len(lines)

	if visibleStart+visibleEnd > len(lines) {
		visibleEnd = len(lines) - visibleStart
	}

	visibleLines := lines
	if visibleStart < len(lines) {
		visibleLines = lines[visibleStart:]
		if len(visibleLines) > visibleEnd {
			visibleLines = visibleLines[:visibleEnd]
		}
	}

	scrollHint := ""
	if hasAbove && hasMore {
		scrollHint = styleFaint.Render("  ↑↓ more")
	} else if hasAbove {
		scrollHint = styleFaint.Render("  ↑ more above")
	} else if hasMore {
		scrollHint = styleFaint.Render("  ↓ more below")
	}

	result := strings.Join(visibleLines, "\n")
	if scrollHint != "" {
		result += "\n" + scrollHint
	}

	return lipgloss.NewStyle().
		Width(width).
		Padding(0, 2).
		Render(result)
}

func (m Model) viewField(label, content, placeholder string, active, editing bool, width int) string {
	var boxStyle lipgloss.Style
	if active || editing {
		boxStyle = styleBorderActive.Width(width-2).Padding(0, 1)
	} else {
		boxStyle = styleBorderNormal.Width(width-2).Padding(0, 1)
	}

	var inner string
	if editing {
		inner = styleBase.Render(m.EditBuffer + "█")
	} else if content != "" {
		inner = styleBase.Render(wrapString(content, width-4))
	} else {
		inner = styleFaint.Render(placeholder)
	}

	labelStyle := styleFaint
	if active || editing {
		labelStyle = styleAmber
	}

	return labelStyle.Render(label) + "\n" + boxStyle.Render(inner)
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
		{"← →", "pane"},
		{"↑↓", "navigate"},
		{"enter", "open/edit"},
		{"ctrl+s", "save"},
		{"q", "quit"},
	}
	if m.ReadOnly {
		binds = []bind{
			{"← →", "pane"},
			{"↑↓", "navigate"},
			{"enter", "open"},
			{"q", "quit"},
		}
	}

	var parts []string
	for _, b := range binds {
		key := lipgloss.NewStyle().Foreground(colAmberBase).Render(b.key)
		desc := styleFaint.Render(" " + b.desc)
		parts = append(parts, key+desc)
	}

	left := strings.Join(parts, "   ")

	right := ""
	if m.Dirty {
		right = styleAmber.Render("● unsaved")
	}
	if m.ReadOnly {
		right = styleFaint.Render("read-only")
	}

	gap := m.Width - lipgloss.Width(left) - lipgloss.Width(right) - 6
	if gap < 0 {
		gap = 0
	}

	return lipgloss.NewStyle().Padding(0, 1).Render(
		left + strings.Repeat(" ", gap) + right,
	)
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
