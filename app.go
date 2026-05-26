package main

import (
    tea "github.com/charmbracelet/bubbletea"
)

func main() {
    p := tea.NewProgram(
        initialModel(),
    )
    if err := p.Start(); err != nil {
        panic(err)
    }
}
