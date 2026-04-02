package ui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/theoneandonlyvabo/grimoire/core"
)

var (
	colBackground = tcell.NewRGBColor(10, 10, 10)
	colSurface    = tcell.NewRGBColor(16, 16, 16)
	colBorder     = tcell.NewRGBColor(201, 169, 110)
	colBorderDim  = tcell.NewRGBColor(50, 48, 40)
	colTextPri    = tcell.NewRGBColor(232, 228, 216)
	colTextMuted  = tcell.NewRGBColor(102, 102, 96)
	colTextDim    = tcell.NewRGBColor(51, 51, 48)
	colAmber      = tcell.NewRGBColor(240, 192, 112)
	colAmberBase  = tcell.NewRGBColor(201, 169, 110)
	colRed        = tcell.NewRGBColor(201, 110, 110)
	colYellow     = tcell.NewRGBColor(201, 194, 110)
	colGreen      = tcell.NewRGBColor(110, 201, 138)
	colBlue       = tcell.NewRGBColor(110, 158, 201)
	colIndigo     = tcell.NewRGBColor(126, 110, 201)
	colPurple     = tcell.NewRGBColor(169, 110, 201)

	stDefault = tcell.StyleDefault.Background(colBackground).Foreground(colTextPri)
	stMuted   = tcell.StyleDefault.Background(colBackground).Foreground(colTextMuted)
	stDim     = tcell.StyleDefault.Background(colBackground).Foreground(colTextDim)
	stAmber   = tcell.StyleDefault.Background(colBackground).Foreground(colAmber)
	stBorder  = tcell.StyleDefault.Background(colBackground).Foreground(colBorder)
	stSurface = tcell.StyleDefault.Background(colSurface).Foreground(colTextPri)
	stActive  = tcell.StyleDefault.Background(colSurface).Foreground(colAmber)
)

func fileIcon(name string) (string, tcell.Color) {
	ext := ""
	if idx := strings.LastIndex(name, "."); idx != -1 {
		ext = strings.ToLower(name[idx:])
	}
	switch ext {
	case ".go":
		return "go", colBlue
	case ".java":
		return "jv", colYellow
	case ".py":
		return "py", colGreen
	case ".js":
		return "js", colYellow
	case ".ts":
		return "ts", colBlue
	case ".rs":
		return "rs", colRed
	case ".rb":
		return "rb", colRed
	case ".c", ".cpp":
		return "c+", colIndigo
	case ".md":
		return "md", colTextMuted
	case ".json":
		return "{}", colTextMuted
	case ".yaml", ".yml":
		return "ym", colTextMuted
	case ".sh":
		return "sh", colGreen
	case ".html":
		return "ht", colRed
	case ".css":
		return "cs", colPurple
	default:
		return " ·", colTextDim
	}
}

func statusColor(status string) tcell.Color {
	switch status {
	case "stable":
		return colGreen
	case "deprecated":
		return colRed
	default:
		return colAmberBase
	}
}

func draw(screen tcell.Screen, x, y int, st tcell.Style, text string) {
	for i, ch := range text {
		screen.SetContent(x+i, y, ch, nil, st)
	}
}

func drawColor(screen tcell.Screen, x, y int, bg tcell.Color, fg tcell.Color, text string) {
	st := tcell.StyleDefault.Background(bg).Foreground(fg)
	for i, ch := range text {
		screen.SetContent(x+i, y, ch, nil, st)
	}
}

func fill(screen tcell.Screen, x1, y1, x2, y2 int, st tcell.Style) {
	for y := y1; y <= y2; y++ {
		for x := x1; x <= x2; x++ {
			screen.SetContent(x, y, ' ', nil, st)
		}
	}
}

func drawBox(screen tcell.Screen, x1, y1, x2, y2 int, st tcell.Style) {
	for x := x1 + 1; x < x2; x++ {
		screen.SetContent(x, y1, '─', nil, st)
		screen.SetContent(x, y2, '─', nil, st)
	}
	for y := y1 + 1; y < y2; y++ {
		screen.SetContent(x1, y, '│', nil, st)
		screen.SetContent(x2, y, '│', nil, st)
	}
	screen.SetContent(x1, y1, '┌', nil, st)
	screen.SetContent(x2, y1, '┐', nil, st)
	screen.SetContent(x1, y2, '└', nil, st)
	screen.SetContent(x2, y2, '┘', nil, st)
}

func drawHLine(screen tcell.Screen, x1, x2, y int, st tcell.Style) {
	screen.SetContent(x1, y, '├', nil, st)
	for x := x1 + 1; x < x2; x++ {
		screen.SetContent(x, y, '─', nil, st)
	}
	screen.SetContent(x2, y, '┤', nil, st)
}

func drawVLine(screen tcell.Screen, x, y1, y2 int, st tcell.Style) {
	for y := y1; y <= y2; y++ {
		screen.SetContent(x, y, '│', nil, st)
	}
}

func renderHeader(screen tcell.Screen, grim *core.Grimoire, w, h int) {
	fill(screen, 1, 1, w-2, 3, stSurface)

	repo := grim.Meta.Repository
	repo = strings.TrimPrefix(repo, "https://github.com/")
	repo = strings.TrimPrefix(repo, "git@github.com:")
	repo = strings.TrimSuffix(repo, ".git")

	row1 := fmt.Sprintf(" ⎇ %s  │  %s  │  Commit ↑ %s · %s ",
		grim.Meta.Branch,
		repo,
		grim.Meta.LastCommit,
		grim.Meta.LastCommitDate[:10],
	)

	x := 2
	drawColor(screen, x, 1, colSurface, colAmber, " GRIMOIRE ")
	x += 10
	drawColor(screen, x, 1, colSurface, colBorderDim, "─")
	x++
	drawColor(screen, x, 1, colSurface, colTextMuted, row1)

	drawHLine(screen, 0, w-1, 2, stBorder)

	msg := fmt.Sprintf("  »  %s", grim.Meta.LastCommitMessage)
	drawColor(screen, 1, 3, colSurface, colTextDim, strings.Repeat(" ", w-2))
	drawColor(screen, 1, 3, colSurface, colTextMuted, msg)
}

func renderSidebar(screen tcell.Screen, state *AppState, x1, y1, x2, y2 int) {
	fill(screen, x1, y1, x2, y2, stDefault)

	draw(screen, x1+2, y1+1, stDim, "files")
	draw(screen, x1+2, y1+2, stDim, strings.Repeat("─", x2-x1-3))

	y := y1 + 3
	for i, node := range state.Tree {
		if y > y2-1 {
			break
		}
		indent := node.Depth * 2

		if node.IsFolder {
			arrow := "▸ "
			if node.Expanded {
				arrow = "▾ "
			}
			fg := colTextMuted
			if node.Expanded {
				fg = colAmberBase
			}
			fill(screen, x1, y, x2-1, y, stDefault)
			drawColor(screen, x1+2+indent, y, colBackground, fg, arrow+node.Name+"/")
		} else {
			isFocused := i == state.ActiveIndex
			bg := colBackground
			if isFocused {
				bg = colSurface
				fill(screen, x1, y, x2-1, y, stActive)
				drawColor(screen, x1, y, bg, colAmber, "▌")
			}

			icon, iconColor := fileIcon(node.Name)
			dotColor := colTextDim
			if node.Doc != nil {
				dotColor = statusColor(node.Doc.Status)
			}

			drawColor(screen, x1+2+indent, y, bg, iconColor, icon)
			drawColor(screen, x1+4+indent, y, bg, colTextDim, " ")
			drawColor(screen, x1+5+indent, y, bg, dotColor, "● ")

			nameStyle := colTextMuted
			if isFocused {
				nameStyle = colAmber
			}
			maxLen := x2 - x1 - 8 - indent
			name := node.Name
			if len(name) > maxLen {
				name = name[:maxLen-1] + "…"
			}
			drawColor(screen, x1+7+indent, y, bg, nameStyle, name)
		}
		y++
	}
}

func renderEditor(screen tcell.Screen, state *AppState, x1, y1, x2, y2 int) {
	fill(screen, x1, y1, x2, y2, stDefault)

	if state.ActiveDoc == nil {
		draw(screen, x1+3, y1+2, stDim, "select a file from the sidebar")
		return
	}

	doc := state.ActiveDoc
	y := y1 + 1
	maxW := x2 - x1 - 4

	draw(screen, x1+3, y, stAmber.Bold(true), doc.LinkedFile)
	y++
	draw(screen, x1+3, y, stDim, fmt.Sprintf("author %s · %s", doc.Author, doc.UpdatedAt))
	y++
	draw(screen, x1+3, y, stDim, strings.Repeat("─", maxW))
	y++

	draw(screen, x1+3, y, stDim, "DESCRIPTION")
	y++
	if doc.Description == "" {
		draw(screen, x1+3, y, stDim, "no description yet...")
		y++
	} else {
		for _, line := range wrapText(doc.Description, maxW) {
			if y > y2-2 {
				break
			}
			draw(screen, x1+3, y, stDefault, line)
			y++
		}
	}
	y++

	if len(doc.Functions) > 0 {
		draw(screen, x1+3, y, stDim, strings.Repeat("─", maxW))
		y++
		draw(screen, x1+3, y, stDim, "FUNCTIONS")
		y++
		for _, fn := range doc.Functions {
			if y > y2-2 {
				break
			}
			_, iconColor := fileIcon(doc.LinkedFile)
			isPrivate := len(fn.Name) > 0 && fn.Name[0] >= 'a' && fn.Name[0] <= 'z'
			fnColor := iconColor
			if isPrivate {
				fnColor = colTextMuted
			}
			sigDisplay := fn.Signature
			if len(sigDisplay) > maxW-len(fn.Name)-2 {
				sigDisplay = sigDisplay[:maxW-len(fn.Name)-3] + "…"
			}
			drawColor(screen, x1+3, y, colBackground, fnColor, fn.Name)
			draw(screen, x1+3+len(fn.Name), y, stDim, "  "+sigDisplay)
			y++
			if fn.Notes != "" {
				draw(screen, x1+5, y, stMuted, "» "+fn.Notes)
				y++
			}
			if fn.Author != "" {
				draw(screen, x1+5, y, stDim, fn.Author+" · "+fn.UpdatedAt)
				y++
			}
		}
	}
}

func renderFooter(screen tcell.Screen, state *AppState, w, h int) {
	fill(screen, 1, h-2, w-2, h-2, stSurface)

	x := 3
	type bind struct {
		key  string
		desc string
	}
	binds := []bind{
		{"Tab", "switch pane"},
		{"↑↓", "navigate"},
		{"Enter", "expand"},
		{"Ctrl+S", "save"},
		{"Q", "quit"},
	}
	if state.ReadOnly {
		binds = []bind{
			{"↑↓", "navigate"},
			{"Enter", "expand"},
			{"Q", "quit"},
		}
	}
	for _, b := range binds {
		drawColor(screen, x, h-2, colSurface, colAmberBase, b.key)
		x += len(b.key) + 1
		drawColor(screen, x, h-2, colSurface, colTextDim, b.desc)
		x += len(b.desc) + 3
	}

	if state.Dirty {
		msg := "● unsaved"
		drawColor(screen, w-len(msg)-2, h-2, colSurface, colAmber, msg)
	}
	if state.ReadOnly {
		msg := "read-only"
		drawColor(screen, w-len(msg)-2, h-2, colSurface, colTextDim, msg)
	}
}

func renderMenu(screen tcell.Screen, selected int) {
	w, h := screen.Size()
	fill(screen, 0, 0, w-1, h-1, stDefault)
	drawBox(screen, 0, 0, w-1, h-1, stBorder)

	cx := w/2 - 4
	cy := h/2 - 3

	draw(screen, cx, cy, stAmber.Bold(true), "GRIMOIRE")
	draw(screen, cx, cy+1, stDim, strings.Repeat("─", 36))

	for i, item := range menuItems {
		y := cy + 3 + i
		if i == selected {
			drawColor(screen, cx-2, y, colBackground, colAmber, "›")
			drawColor(screen, cx, y, colBackground, colAmber, item.command)
		} else {
			draw(screen, cx-2, y, stDim, " ")
			draw(screen, cx, y, stMuted, item.command)
		}
		draw(screen, cx+10, y, stDim, item.description)
	}

	draw(screen, cx, cy+8, stDim, "↑↓ navigate   Enter select   Q quit")
}

func render(screen tcell.Screen, state *AppState) {
	w, h := screen.Size()

	fill(screen, 0, 0, w-1, h-1, stDefault)
	drawBox(screen, 0, 0, w-1, h-1, stBorder)

	renderHeader(screen, state.Grimoire, w, h)

	sidebarX2 := 24
	bodyY1 := 4
	bodyY2 := h - 3

	drawHLine(screen, 0, w-1, bodyY1, stBorder)
	drawVLine(screen, sidebarX2, bodyY1, bodyY2, stBorder)
	screen.SetContent(sidebarX2, bodyY1, '┬', nil, stBorder)
	screen.SetContent(0, bodyY1, '├', nil, stBorder)
	screen.SetContent(w-1, bodyY1, '┤', nil, stBorder)

	drawHLine(screen, 0, w-1, bodyY2, stBorder)
	screen.SetContent(0, bodyY2, '├', nil, stBorder)
	screen.SetContent(w-1, bodyY2, '┤', nil, stBorder)
	screen.SetContent(sidebarX2, bodyY2, '┴', nil, stBorder)

	renderSidebar(screen, state, 1, bodyY1, sidebarX2-1, bodyY2)
	renderEditor(screen, state, sidebarX2+1, bodyY1, w-2, bodyY2)
	renderFooter(screen, state, w, h)
}

func wrapText(text string, width int) []string {
	if width <= 0 {
		return []string{text}
	}
	var lines []string
	words := strings.Fields(text)
	current := ""
	for _, word := range words {
		if len(current)+len(word)+1 > width {
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
	return lines
}
