package ui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/theoneandonlyvabo/grimoire/core"
)

var (
	colBackground = tcell.NewRGBColor(10, 10, 10)
	colSurface    = tcell.NewRGBColor(18, 18, 18)
	colBorder     = tcell.NewRGBColor(201, 169, 110)
	colTextPri    = tcell.NewRGBColor(232, 228, 216)
	colTextMuted  = tcell.NewRGBColor(100, 100, 94)
	colTextDim    = tcell.NewRGBColor(50, 50, 46)
	colAmber      = tcell.NewRGBColor(240, 192, 112)
	colAmberBase  = tcell.NewRGBColor(201, 169, 110)
	colRed        = tcell.NewRGBColor(201, 110, 110)
	colGreen      = tcell.NewRGBColor(110, 201, 138)
	colBlue       = tcell.NewRGBColor(110, 158, 201)

	stDefault = tcell.StyleDefault.Background(colBackground).Foreground(colTextPri)
	stMuted   = tcell.StyleDefault.Background(colBackground).Foreground(colTextMuted)
	stDim     = tcell.StyleDefault.Background(colBackground).Foreground(colTextDim)
	stAmber   = tcell.StyleDefault.Background(colBackground).Foreground(colAmber)
	stBorder  = tcell.StyleDefault.Background(colBackground).Foreground(colBorder)
	stSurface = tcell.StyleDefault.Background(colSurface).Foreground(colTextPri)
	stActive  = tcell.StyleDefault.Background(colBackground).Foreground(colAmber)
)

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

func drawCol(screen tcell.Screen, x, y int, bg, fg tcell.Color, text string) {
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

func hline(screen tcell.Screen, x1, x2, y int, leftCap, rightCap rune, st tcell.Style) {
	screen.SetContent(x1, y, leftCap, nil, st)
	for x := x1 + 1; x < x2; x++ {
		screen.SetContent(x, y, '─', nil, st)
	}
	screen.SetContent(x2, y, rightCap, nil, st)
}

func vline(screen tcell.Screen, x, y1, y2 int, st tcell.Style) {
	for y := y1; y <= y2; y++ {
		screen.SetContent(x, y, '│', nil, st)
	}
}

func sep(screen tcell.Screen, x1, x2, y int, st tcell.Style) {
	for x := x1; x <= x2; x++ {
		screen.SetContent(x, y, '─', nil, st)
	}
}

func renderHeader(screen tcell.Screen, grim *core.Grimoire, w int) {
	fill(screen, 1, 1, w-2, 3, stSurface)

	repo := grim.Meta.Repository
	repo = strings.TrimPrefix(repo, "https://github.com/")
	repo = strings.TrimPrefix(repo, "git@github.com:")
	repo = strings.TrimSuffix(repo, ".git")

	date := grim.Meta.LastCommitDate
	if len(date) > 10 {
		date = date[:10]
	}

	x := 2
	drawCol(screen, x, 1, colSurface, colAmber, "GRIMOIRE")
	x += 9
	drawCol(screen, x, 1, colSurface, colTextDim, "─")
	x += 2
	drawCol(screen, x, 1, colSurface, colBlue, "⎇ "+grim.Meta.Branch)
	x += len("⎇ "+grim.Meta.Branch) + 2
	drawCol(screen, x, 1, colSurface, colTextDim, "│")
	x += 2
	drawCol(screen, x, 1, colSurface, colTextMuted, repo)
	x += len(repo) + 2
	drawCol(screen, x, 1, colSurface, colTextDim, "│")
	x += 2
	commitStr := fmt.Sprintf("Commit ↑ %s · %s", grim.Meta.LastCommit, date)
	drawCol(screen, x, 1, colSurface, colAmberBase, commitStr)

	hline(screen, 0, w-1, 2, '├', '┤', stBorder)

	msgLine := fmt.Sprintf("  »  %s", grim.Meta.LastCommitMessage)
	fill(screen, 1, 3, w-2, 3, stSurface)
	drawCol(screen, 1, 3, colSurface, colTextDim, strings.Repeat(" ", w-2))
	drawCol(screen, 2, 3, colSurface, colTextMuted, msgLine)
}

func renderSidebar(screen tcell.Screen, state *AppState, x1, y1, x2, y2 int) {
	fill(screen, x1, y1, x2, y2, stDefault)

	draw(screen, x1+2, y1+1, stDim, "files")
	sep(screen, x1+2, x2-1, y1+2, stDim)

	y := y1 + 3
	for i, node := range state.Tree {
		if y > y2-1 {
			break
		}
		indent := node.Depth * 2
		isFocused := i == state.ActiveIndex && state.ActivePane == 0

		if node.IsFolder {
			arrow := "▸ "
			if node.Expanded {
				arrow = "▾ "
			}
			fg := colTextMuted
			if isFocused {
				fg = colAmber
			} else if node.Expanded {
				fg = colAmberBase
			}
			if isFocused {
				drawCol(screen, x1, y, colBackground, colAmber, "▌")
			}
			drawCol(screen, x1+2+indent, y, colBackground, fg, arrow+node.Name+"/")
		} else {
			fg := colTextMuted
			if isFocused {
				fg = colAmber
				drawCol(screen, x1, y, colBackground, colAmber, "▌")
			}

			dotColor := colTextDim
			if node.Doc != nil {
				dotColor = statusColor(node.Doc.Status)
			}

			drawCol(screen, x1+2+indent, y, colBackground, dotColor, "●")

			maxLen := x2 - x1 - 5 - indent
			name := node.Name
			if len(name) > maxLen {
				name = name[:maxLen-1] + "…"
			}
			drawCol(screen, x1+4+indent, y, colBackground, fg, name)
		}
		y++
	}
}

func renderEditor(screen tcell.Screen, state *AppState, x1, y1, x2, y2 int) {
	fill(screen, x1, y1, x2, y2, stDefault)

	if state.ActiveDoc == nil {
		draw(screen, x1+3, y1+3, stDim, "select a file from the sidebar")
		return
	}

	doc := state.ActiveDoc
	y := y1 + 2
	maxW := x2 - x1 - 5

	drawCol(screen, x1+3, y, colBackground, colAmber, doc.LinkedFile)
	y++
	drawCol(screen, x1+3, y, colBackground, colTextDim,
		fmt.Sprintf("author %s  ·  %s", doc.Author, doc.UpdatedAt))
	y++
	sep(screen, x1+3, x2-2, y, stDim)
	y++

	draw(screen, x1+3, y, stDim, "DESCRIPTION")
	y++

	desc := doc.Description
	if desc == "" {
		draw(screen, x1+3, y, stDim, "no description yet...")
		y++
	} else {
		for _, line := range wrapText(desc, maxW) {
			if y > y2-2 {
				break
			}
			draw(screen, x1+3, y, stDefault, line)
			y++
		}
	}

	if len(doc.Functions) > 0 {
		y++
		sep(screen, x1+3, x2-2, y, stDim)
		y++
		draw(screen, x1+3, y, stDim, "FUNCTIONS")
		y++

		for _, fn := range doc.Functions {
			if y > y2-2 {
				break
			}
			isPrivate := len(fn.Name) > 0 && fn.Name[0] >= 'a' && fn.Name[0] <= 'z'
			fnColor := colBlue
			if isPrivate {
				fnColor = colTextMuted
			}

			maxSig := maxW - len(fn.Name) - 3
			sig := fn.Signature
			if len(sig) > maxSig {
				sig = sig[:maxSig] + "…"
			}

			drawCol(screen, x1+3, y, colBackground, fnColor, fn.Name)
			drawCol(screen, x1+3+len(fn.Name)+2, y, colBackground, colTextDim, sig)
			y++

			if fn.Notes != "" {
				draw(screen, x1+5, y, stMuted, "»  "+fn.Notes)
				y++
			}
		}
	}
}

func renderFooter(screen tcell.Screen, state *AppState, w, h int) {
	fill(screen, 1, h-2, w-2, h-2, stSurface)

	type bind struct{ key, desc string }
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
			{"Enter", "open"},
			{"Q", "quit"},
		}
	}

	x := 3
	for _, b := range binds {
		drawCol(screen, x, h-2, colSurface, colAmber, b.key)
		x += len(b.key) + 1
		drawCol(screen, x, h-2, colSurface, colTextDim, b.desc)
		x += len(b.desc) + 3
	}

	if state.Dirty {
		msg := "● unsaved"
		drawCol(screen, w-len(msg)-2, h-2, colSurface, colAmber, msg)
	}
	if state.ReadOnly {
		msg := "read-only"
		drawCol(screen, w-len(msg)-2, h-2, colSurface, colTextDim, msg)
	}
}

func renderMenu(screen tcell.Screen, selected int) {
	w, h := screen.Size()
	fill(screen, 0, 0, w-1, h-1, stDefault)
	drawBox(screen, 0, 0, w-1, h-1, stBorder)

	cx := w/2 - 4
	cy := h/2 - 3

	draw(screen, cx, cy, stAmber.Bold(true), "GRIMOIRE")
	sep(screen, cx, cx+35, cy+1, stDim)

	for i, item := range menuItems {
		y := cy + 3 + i
		if i == selected {
			drawCol(screen, cx-2, y, colBackground, colAmber, "›")
			drawCol(screen, cx, y, colBackground, colAmber, item.command)
		} else {
			draw(screen, cx, y, stMuted, item.command)
		}
		draw(screen, cx+12, y, stDim, item.description)
	}

	draw(screen, cx, cy+8, stDim, "↑↓ navigate   Enter select   Q quit")
}

func render(screen tcell.Screen, state *AppState) {
	w, h := screen.Size()

	fill(screen, 0, 0, w-1, h-1, stDefault)
	drawBox(screen, 0, 0, w-1, h-1, stBorder)

	renderHeader(screen, state.Grimoire, w)

	sidebarX := 26
	bodyY1 := 4
	bodyY2 := h - 3

	hline(screen, 0, w-1, bodyY1, '├', '┤', stBorder)
	vline(screen, sidebarX, bodyY1, bodyY2, stBorder)
	hline(screen, 0, w-1, bodyY2, '├', '┤', stBorder)

	screen.SetContent(sidebarX, bodyY1, '┬', nil, stBorder)
	screen.SetContent(sidebarX, bodyY2, '┴', nil, stBorder)

	renderSidebar(screen, state, 1, bodyY1, sidebarX-1, bodyY2)
	renderEditor(screen, state, sidebarX+1, bodyY1, w-2, bodyY2)
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
