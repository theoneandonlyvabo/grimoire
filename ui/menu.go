package ui

import (
	"github.com/gdamore/tcell/v2"
)

type menuItem struct {
	command     string
	description string
}

var menuItems = []menuItem{
	{"forge", "initialize grimoire in this project"},
	{"carve", "write technical notes"},
	{"cast", "read grimoire"},
}

func StartMenu() error {
	screen, err := tcell.NewScreen()
	if err != nil {
		return err
	}
	if err := screen.Init(); err != nil {
		return err
	}
	defer screen.Fini()

	selected := 0

	for {
		screen.Clear()
		renderMenu(screen, selected)
		screen.Show()

		event := screen.PollEvent()
		switch ev := event.(type) {
		case *tcell.EventResize:
			screen.Sync()
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyUp:
				if selected > 0 {
					selected--
				}
			case tcell.KeyDown:
				if selected < len(menuItems)-1 {
					selected++
				}
			case tcell.KeyEnter:
				screen.Fini()
				return runFromMenu(menuItems[selected].command)
			case tcell.KeyEscape:
				return nil
			case tcell.KeyRune:
				if ev.Rune() == 'q' || ev.Rune() == 'Q' {
					return nil
				}
			}
		}
	}
}

func runFromMenu(command string) error {
	switch command {
	case "forge":
		return RunForge()
	case "carve":
		return RunCarve()
	case "cast":
		return RunCast()
	}
	return nil
}
