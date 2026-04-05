package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/theoneandonlyvabo/grimoire/core"
)

const asciiArt = `
 ██████╗ ██████╗ ██╗███╗   ███╗ ██████╗ ██╗██████╗ ███████╗
██╔════╝ ██╔══██╗██║████╗ ████║██╔═══██╗██║██╔══██╗██╔════╝
██║  ███╗██████╔╝██║██╔████╔██║██║   ██║██║██████╔╝█████╗  
██║   ██║██╔══██╗██║██║╚██╔╝██║██║   ██║██║██╔══██╗██╔══╝  
╚██████╔╝██║  ██║██║██║ ╚═╝ ██║╚██████╔╝██║██║  ██║███████╗
 ╚═════╝ ╚═╝  ╚═╝╚═╝╚═╝     ╚═╝ ╚═════╝ ╚═╝╚═╝  ╚═╝╚══════╝`

type menuModel struct {
	selected  int
	confirmed bool
	items     []menuItem
	width     int
	height    int
	version   core.VersionInfo
}

type menuItem struct {
	command     string
	description string
}

var menuItems = []menuItem{
	{"forge", "initialize grimoire in this project"},
	{"carve", "write technical notes"},
	{"cast", "read grimoire"},
}

func (m menuModel) Init() tea.Cmd {
	return nil
}

func (m menuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.selected > 0 {
				m.selected--
			}
		case "down":
			if m.selected < len(m.items)-1 {
				m.selected++
			}
		case "enter":
			m.confirmed = true
			return m, tea.Quit
		case "q", "Q", "ctrl+c", "esc":
			m.selected = -1
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m menuModel) View() string {
	if m.width == 0 {
		return ""
	}

	var brand string
	if m.width >= 75 {
		brand = lipgloss.NewStyle().
			Foreground(colAmber).
			Render(asciiArt)
	} else {
		brand = lipgloss.NewStyle().
			Foreground(colAmber).
			Bold(true).
			Render("GRIMOIRE")
	}

	release := lipgloss.NewStyle().
		Foreground(colTextDim).
		Render("release  " + m.version.Release)

	divider := lipgloss.NewStyle().
		Foreground(colTextDim).
		Render("────────────────────────────────────────")

	var itemLines []string
	for i, item := range m.items {
		if i == m.selected {
			cursor := lipgloss.NewStyle().Foreground(colAmber).Render("›")
			cmd := lipgloss.NewStyle().Foreground(colAmber).Render(item.command)
			desc := lipgloss.NewStyle().Foreground(colTextMuted).Render("  " + item.description)
			itemLines = append(itemLines, cursor+" "+cmd+desc)
		} else {
			cmd := lipgloss.NewStyle().Foreground(colTextMuted).Render("  " + item.command)
			desc := lipgloss.NewStyle().Foreground(colTextDim).Render("  " + item.description)
			itemLines = append(itemLines, cmd+desc)
		}
	}

	hints := lipgloss.NewStyle().
		Foreground(colTextDim).
		Render("↑↓ navigate   enter select   q quit")

	padding := lipgloss.NewStyle().Padding(1, 3)

	content := lipgloss.JoinVertical(lipgloss.Left,
		brand,
		"",
		release,
		divider,
		"",
		itemLines[0],
		itemLines[1],
		itemLines[2],
		"",
		hints,
	)

	return padding.Render(content)
}

func StartMenu() error {
	v := core.GetVersion()
	m := menuModel{
		items:    menuItems,
		selected: 0,
		version:  v,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return err
	}

	result, ok := finalModel.(menuModel)
	if !ok || !result.confirmed || result.selected < 0 {
		return nil
	}

	switch result.items[result.selected].command {
	case "forge":
		return RunForge()
	case "carve":
		return RunCarve()
	case "cast":
		return RunCast()
	}

	return nil
}
