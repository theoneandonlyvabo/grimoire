package ui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/theoneandonlyvabo/grimoire/core"
)

var (
	styleDefault = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.NewRGBColor(201, 199, 190))
	styleDim     = tcell.StyleDefault.Foreground(tcell.NewRGBColor(85, 85, 96))
	styleAccent  = tcell.StyleDefault.Foreground(tcell.NewRGBColor(201, 169, 110))
	styleGreen   = tcell.StyleDefault.Foreground(tcell.NewRGBColor(58, 138, 90))
	styleRed     = tcell.StyleDefault.Foreground(tcell.NewRGBColor(138, 58, 58))
	styleMuted   = tcell.StyleDefault.Foreground(tcell.NewRGBColor(106, 106, 117))
	styleHeader  = tcell.StyleDefault.Background(tcell.NewRGBColor(16, 16, 18)).Foreground(tcell.NewRGBColor(201, 199, 190))
	styleSidebar = tcell.StyleDefault.Background(tcell.NewRGBColor(12, 12, 14)).Foreground(tcell.NewRGBColor(85, 85, 96))
	styleActive  = tcell.StyleDefault.Background(tcell.NewRGBColor(24, 24, 32)).Foreground(tcell.NewRGBColor(201, 169, 110))
	styleFooter  = tcell.StyleDefault.Background(tcell.NewRGBColor(16, 16, 18)).Foreground(tcell.NewRGBColor(53, 53, 56))
)

func drawText(screen tcell.Screen, x, y int, style tcell.Style, text string) {
	for i, ch := range text {
		screen.SetContent(x+i, y, ch, nil, style)
	}
}

func drawFill(screen tcell.Screen, x1, y1, x2, y2 int, style tcell.Style) {
	for y := y1; y <= y2; y++ {
		for x := x1; x <= x2; x++ {
			screen.SetContent(x, y, ' ', nil, style)
		}
	}
}

func renderHeader(screen tcell.Screen, grim *core.Grimoire, width int) {
	drawFill(screen, 0, 0, width, 0, styleHeader)
	drawFill(screen, 0, 1, width, 1, styleHeader)

	x := 1
	drawText(screen, x, 0, styleAccent.Bold(true), "GRIMOIRE")
	x += 9
	drawText(screen, x, 0, styleDim, "|")
	x += 2
	drawText(screen, x, 0, styleGreen, "⎇ "+grim.Meta.Branch)
	x += len("⎇ "+grim.Meta.Branch) + 2
	drawText(screen, x, 0, styleDim, "|")
	x += 2

	repo := grim.Meta.Repository
	repo = strings.TrimPrefix(repo, "https://github.com/")
	repo = strings.TrimPrefix(repo, "git@github.com:")
	drawText(screen, x, 0, styleMuted, repo)
	x += len(repo) + 2
	drawText(screen, x, 0, styleDim, "|")
	x += 2
	drawText(screen, x, 0, styleAccent, "↑ "+grim.Meta.LastCommit)
	x += len("↑ "+grim.Meta.LastCommit) + 1
	drawText(screen, x, 0, styleDim, "· "+grim.Meta.LastCommitDate[:10])

	drawText(screen, 1, 1, styleDim, "» ")
	drawText(screen, 3, 1, styleMuted, grim.Meta.LastCommitMessage)
}

func renderSidebar(screen tcell.Screen, state *AppState, startY int, height int) {
	sidebarWidth := 22
	drawFill(screen, 0, startY, sidebarWidth, startY+height, styleSidebar)

	y := startY
	drawText(screen, 1, y, styleDim, "files")
	y++

	for i, node := range state.Tree {
		if y >= startY+height {
			break
		}
		indent := node.Depth * 2
		prefix := strings.Repeat(" ", indent)

		if node.IsFolder {
			arrow := "▸ "
			if node.Expanded {
				arrow = "▾ "
			}
			folderStyle := styleDim
			if node.Expanded {
				folderStyle = styleGreen
			}
			drawText(screen, 1+indent, y, folderStyle, arrow+node.Name+"/")
		} else {
			dot := "·"
			dotStyle := styleDim
			if node.Doc != nil {
				switch node.Doc.Status {
				case "stable":
					dotStyle = styleGreen
				case "deprecated":
					dotStyle = styleRed
				default:
					dotStyle = styleAccent
				}
			}
			_ = prefix
			rowStyle := styleSidebar.Foreground(tcell.NewRGBColor(85, 85, 96))
			if i == state.ActiveIndex {
				rowStyle = styleActive
				drawFill(screen, 0, y, sidebarWidth, y, styleActive)
				drawText(screen, 0, y, styleAccent, "▌")
			}
			drawText(screen, 1+indent, y, dotStyle, dot)
			drawText(screen, 3+indent, y, rowStyle, node.Name)
		}
		y++
	}
}

func renderEditor(screen tcell.Screen, state *AppState, startX, startY, width, height int) {
	drawFill(screen, startX, startY, startX+width, startY+height, styleDefault)

	if state.ActiveDoc == nil {
		drawText(screen, startX+2, startY+2, styleDim, "select a file from the sidebar")
		return
	}

	doc := state.ActiveDoc
	y := startY + 1

	drawText(screen, startX+2, y, styleAccent.Bold(true), doc.LinkedFile)
	y++
	drawText(screen, startX+2, y, styleDim, fmt.Sprintf("author %s · %s", doc.Author, doc.UpdatedAt))
	y += 2

	drawText(screen, startX+2, y, styleDim, "DESCRIPTION")
	y++
	desc := doc.Description
	if desc == "" {
		desc = "no description yet..."
		drawText(screen, startX+2, y, styleDim, desc)
	} else {
		wrapped := wrapText(desc, width-6)
		for _, line := range wrapped {
			drawText(screen, startX+2, y, styleDefault, line)
			y++
		}
		y--
	}
	y += 2

	if len(doc.Functions) > 0 {
		drawText(screen, startX+2, y, styleDim, "FUNCTIONS")
		y++
		for _, fn := range doc.Functions {
			if y >= startY+height {
				break
			}
			drawText(screen, startX+2, y, styleGreen, fn.Name)
			drawText(screen, startX+2+len(fn.Name)+1, y, styleDim, fn.Signature[len(fn.Name):])
			y++
			if fn.Notes != "" {
				drawText(screen, startX+4, y, styleMuted, "» "+fn.Notes)
				y++
			}
		}
	}
}

func renderMenu(screen tcell.Screen, selected int) {
	w, h := screen.Size()
	cx := w / 2
	cy := h / 2

	drawText(screen, cx-4, cy-3, styleAccent.Bold(true), "GRIMOIRE")

	for i, item := range menuItems {
		y := cy - 1 + i
		style := styleDim
		prefix := "  "
		if i == selected {
			style = styleAccent
			prefix = "› "
		}
		drawText(screen, cx-4, y, style, prefix+item.command)
		drawText(screen, cx+8, y, styleDim, item.description)
	}

	drawText(screen, cx-4, cy+4, styleDim, "↑↓ navigate   Enter select   Q quit")
}

func renderFooter(screen tcell.Screen, state *AppState, y, width int) {
	drawFill(screen, 0, y, width, y, styleFooter)
	x := 1
	binds := []string{"Tab switch pane", "↑↓ navigate", "Enter expand", "Ctrl+S save", "Q quit"}
	if state.ReadOnly {
		binds = []string{"↑↓ navigate", "Enter expand", "Q quit"}
	}
	for _, b := range binds {
		drawText(screen, x, y, styleFooter, b)
		x += len(b) + 3
	}
	if state.Dirty {
		dirty := "● unsaved"
		drawText(screen, width-len(dirty)-1, y, styleAccent, dirty)
	}
	if state.ReadOnly {
		drawText(screen, width-10, y, styleDim, "read-only")
	}
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

func render(screen tcell.Screen, state *AppState) {
	width, height := screen.Size()
	sidebarWidth := 22

	renderHeader(screen, state.Grimoire, width)
	renderSidebar(screen, state, 2, height-4)
	renderEditor(screen, state, sidebarWidth+1, 2, width-sidebarWidth-1, height-4)
	renderFooter(screen, state, height-1, width)

	screen.SetContent(sidebarWidth, 2, '│', nil, styleDim)
	for y := 3; y < height-1; y++ {
		screen.SetContent(sidebarWidth, y, '│', nil, styleDim)
	}
}
