package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type menuModel struct {
	selected  int
	confirmed bool
	items     []menuItem
	width     int
	height    int
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

	brand := lipgloss.NewStyle().
		Foreground(colAmber).
		Bold(true).
		Render("GRIMOIRE")

	divider := styleVeryDim.Render("────────────────────────────────────")

	var itemLines []string
	for i, item := range m.items {
		var line string
		if i == m.selected {
			cursor := styleAmber.Render("›")
			cmd := styleAmber.Render(item.command)
			desc := styleVeryDim.Render("  " + item.description)
			line = cursor + " " + cmd + desc
		} else {
			cmd := styleDim.Render("  " + item.command)
			desc := styleVeryDim.Render("  " + item.description)
			line = cmd + desc
		}
		itemLines = append(itemLines, line)
	}

	hints := styleVeryDim.Render("↑↓ navigate   enter select   q quit")

	content := fmt.Sprintf("%s\n%s\n\n%s\n\n%s",
		brand,
		divider,
		joinLines(itemLines),
		hints,
	)

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colBorder).
		Padding(1, 3)

	boxed := box.Render(content)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		boxed,
	)
}

type selectMsg struct {
	command string
}

func joinLines(lines []string) string {
	result := ""
	for i, l := range lines {
		if i > 0 {
			result += "\n"
		}
		result += l
	}
	return result
}

func StartMenu() error {
	m := menuModel{items: menuItems, selected: -1}
	p := tea.NewProgram(m, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		return err
	}

	result, ok := finalModel.(menuModel)
	if !ok || result.selected < 0 {
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
