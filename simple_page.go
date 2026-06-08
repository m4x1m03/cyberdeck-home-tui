package main

import (
	"os"
	"os/exec"
	"time"
	"strings"
	"strconv"

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

	leftPanelStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#874BFD"))

	rightPanelStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#555555")).
		Padding(1, 2)

	statusBarStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("#7D56F4")).
		Foreground(lipgloss.Color("#FAFAFA")).
		Padding(0, 1)

	previewTitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		MarginBottom(1)

	errorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF5555"))
)


// App structure
type App struct{
	Name string
	Icon string
	Binary string
	Args []string
	Preview string
}

var apps = []App{
	{
		Name: "Music",
		Icon: "",
		Binary: "",
		Args: []string{},
		Preview: "This where the good tunes are",
	},
	{
		Name: "Anime",
		Icon: "",
		Binary: "",
		Args: []string{},
		Preview: "Good things to watch",
	},
	{
		Name: "Files",
		Icon: "",
		Binary: "",
		Args: []string{},
		Preview: "Don't look at my secret files",
	},

}

// Extra msg structs
type tickMsg time.Time
type appFinishedMsg struct{ err error }
type batteryMsg struct{
	level string
	status string
}


// Battery funcs
func getBattery() tea.Msg {
	capacity, err := os.ReadFile("/sys/class/power_supply/BAT0/capacity")
	if err != nil {
		capacity, err = os.ReadFile("/sys/class/power_supply/BAT1/capacity")
		if err != nil {
			return batteryMsg{level: "--", status: ""}
		}
	}

	status, err := os.ReadFile("/sys/class/power_supply/BAT0/status")
	if err != nil {
		status, _ = os.ReadFile("/sys/class/power_supply/BAT1/status")
	}

	level := strings.TrimSpace(string(capacity))
	state := strings.TrimSpace(string(status))
	return batteryMsg{level: level, status: state}
}

func pollBattery() tea.Cmd {
	return tea.Tick(30*time.Second, func(t time.Time) tea.Msg {
		return getBattery()
	})
}

func batteryIcon(level string, status string) string {
	if status == "Charging" {
		return "󰂄"
	}
	n, err := strconv.Atoi(level)
	if err != nil {
		return "󰂑"
	}
	switch {
	case n >= 90: return "󰁹"
	case n >= 70: return "󰂀"
	case n >= 50: return "󰁾"
	case n >= 30: return "󰁼"
	case n >= 10: return "󰁺"
	default:      return "󰂃" // critical
	}
}


// TUI Model structure

type model struct {
	cursor int
	width    int
	height   int
	time time.Time
	lastErr error
	battery string
	batteryStatus string
}

func initialModel() model {
	return model{
		time: time.Now(),
		battery: "--",
	}
}



// Init Model

func (m model) Init() tea.Cmd {
	return tea.Batch(
		tea.Tick(time.Second, func(t time.Time) tea.Msg { return tickMsg(t) }),
		getBattery, 
		)
}

// View Model
func (m model) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}

	statusBar := m.renderStatusBar()
	statusHeight := lipgloss.Height(statusBar)

	panelHeight := m.height - statusHeight

	leftPanel  := m.renderLeftPanel(panelHeight)
	rightPanel := m.renderRightPanel(panelHeight, lipgloss.Width(leftPanel))

	main := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)
	return lipgloss.JoinVertical(lipgloss.Left, main, statusBar)
}

func (m model) renderLeftPanel(panelHeight int) string {
	leftWidth := m.width / 3

	var list string
	for i, app := range apps {
		label := app.Icon + "  " + app.Name
		if i == m.cursor {
			list += selectedStyle.Render("> "+label) + "\n"
		} else {
			list += itemStyle.Render("  "+label) + "\n"
		}
	}

	content := titleStyle.Render("Apps") + "\n\n" + list

	return leftPanelStyle.
		Width(leftWidth - 2).
		Height(panelHeight - 2).
		Render(content)
}

func (m model) renderRightPanel(panelHeight, leftWidth int) string {
	rightWidth := m.width - leftWidth

	app := apps[m.cursor]

	content := previewTitleStyle.Render(app.Icon+"  "+app.Name) +
	"\n" + app.Preview +
	"\n\n" + itemStyle.Render("Press Enter to launch")

	return rightPanelStyle.
		Width(rightWidth - 2).
		Height(panelHeight - 2).
		Render(content)
}

func (m model) renderStatusBar() string {
	clock := m.time.Format("15:04:05")
	batIcon := batteryIcon(m.battery, m.batteryStatus)
	battery := batIcon + " " + m.battery + "%"

	left := " 󰒋 cyberdeck"
	right := clock + " | " + battery + " | weather | volume % | wifi status"

	totalWidth := m.width
	middle := totalWidth - lipgloss.Width(left) - lipgloss.Width(right) - 2
	if middle < 0 {
		middle = 0
	}

	var errMsg string
	if m.lastErr != nil {
		errMsg = errorStyle.Render("  error: " + m.lastErr.Error())
	}

	bar := statusBarStyle.Width(m.width).Render(
		left + errMsg + lipgloss.NewStyle().Width(middle).Render("") + right,
		)
	return bar
}


// Update Model

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tickMsg:
		m.time = time.Time(msg)
		return m, tea.Tick(time.Second, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})

	case appFinishedMsg:
		m.lastErr = msg.err
		return m, nil

	case batteryMsg:
		m.battery = msg.level
		m.batteryStatus = msg.status
		return m, pollBattery()

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(apps)-1 {
				m.cursor++
			}
		case "enter", " ":
			return m, launchApp(apps[m.cursor])
		}
	}
	return m, nil
}


// launching selected application given arguments in app list
func launchApp(app App) tea.Cmd {
	binary, err := exec.LookPath(app.Binary)
	if err != nil {
		return func() tea.Msg {
			return appFinishedMsg{err}
		}
	}

	c := exec.Command(binary, app.Args...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	return tea.ExecProcess(c, func(err error) tea.Msg {
		return appFinishedMsg{err}
	})
}
