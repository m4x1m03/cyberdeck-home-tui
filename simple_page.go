package main

import (

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

//styles
var (
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1)

	itemStyle = lipgloss.NewStyle().
		PaddingLeft(2)

	selectedStyle = lipgloss.NewStyle().
		PaddingLeft(2).
		Foreground(lipgloss.Color("#7D56F4")).
		Bold(true)

	boxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#874BFD")).
		Padding(1, 2)

	statusBarStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		Width(0)

)

// MODEL DATA

type model struct {
	apps []string
	cursor int
	selected map[int]bool
	width    int
	height   int
}

func initialModel() model {
	return model{
		apps: []string{"Music", "Anime", "Pictures"},
		selected: make(map[int]bool),
	}
}

func (m model) Init() tea.Cmd { return nil }

// VIEW

// Main box apps
func (m model) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}
	title := titleStyle.Render("Apps")

	var apps string

	for i, app := range m.apps{
		if i == m.cursor {
			apps += selectedStyle.Render(app) + "\n"
		} else {
			apps += itemStyle.Render(app) + "\n"
		}
	}

	content := title + "\n\n" + apps

	sized := boxStyle.
		Width(m.width - 2).
		Height(m.height - 4).
		Align(lipgloss.Center, lipgloss.Center).
		Render(content)


	//Status Bar
	status := statusBarStyle.
        Width(m.width).
        Render("cyberdeck  |  date time  |  battery  |  weather  |  volume %  |  wifi status")


	return lipgloss.JoinVertical(lipgloss.Left, sized, status)
}

// UPDATE

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
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
