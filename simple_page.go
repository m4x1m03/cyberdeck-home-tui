package main

import (
    "fmt"

    tea "github.com/charmbracelet/bubbletea"
)

// MODEL DATA

type model struct {
	apps []string
	cursor int
	selected map[int]bool
}

func initialModel() model {
    return model{
			apps: []string{"Music", "Anime", "Pictures"},
			selected: make(map[int]bool),
	}
}

func (m model) Init() tea.Cmd { return nil }

// VIEW

func (m model) View() string {
	s := "Apps:\n\n"

	for i, apps := range m.apps{
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if m.selected[i]{
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, apps)
	}

	s += "\nPress q to quit.\n"
	return s
}

// UPDATE

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg.(type) {
    case tea.KeyMsg:
        switch msg.(tea.KeyMsg).String() {
        case "ctrl+c":
            return m, tea.Quit
      	case "up", "k":
						if m.cursor > 0 {
							m.cursor--
						}
				case "down", "j":
						if m.cursor < len(m.apps)-1 {
							m.cursor++
						}
				case "enter", " ":
						m.selected[m.cursor] = !m.selected[m.cursor]
			}
    }
    return m, nil
}
