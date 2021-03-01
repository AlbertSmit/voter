package main

// A simple program that makes a GET request and prints the response status.

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	spinner "github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

const url = "https://api.mocki.io/v1/d201c428"

type model struct {
	rooms			[]string

	cursor		int
	selected	int

	spinner   spinner.Model
	status 		int
	err    		error
}

type statusMsg int
type payloadMsg []string
type errMsg struct{ error }

func (e errMsg) Error() string { return e.Error() }

func main() {
	p := tea.NewProgram(model{
		spinner: spinner.NewModel(),
		selected: 0,
		cursor: 0,
	})

	termenv.ClearScreen()
	p.EnterAltScreen()

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}

func (m model) Init() tea.Cmd {
	return checkServer
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			fallthrough
		case "ctrl+c", "q":
			return m, tea.Quit
		case "r":
			return m, checkServer
		case "up", "k":
			if m.cursor > 0 {
					m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.rooms)-1 {
					m.cursor++
			}
		case "enter", " ":
			m.selected = m.cursor + 1
		default:
			return m, nil
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case payloadMsg:
		m.rooms = []string(msg)
		return m, nil

	case statusMsg:
		m.status = int(msg)
		return m, tea.Quit

	case errMsg:
		m.err = msg
		return m, nil

	default:
		return m, nil
	}

	return m, nil
}

func (m model) View() string {
	p := termenv.ColorProfile()
	var s string

	// The header
	if (len(m.rooms) == 0) {
		s += fmt.Sprintf("\n %s Fetching URL: %v\n\n", m.spinner.View(), url)
	} else {
		s += fmt.Sprintf("%s", termenv.String("\nResults\n\n").Bold())
	}

	for i, choice := range m.rooms {
		cursor := " " 
		if m.cursor == i {
				cursor = "→" 
		}

		s += fmt.Sprintf(
			"%s %s\n", 
			termenv.String(cursor).Foreground(p.Color("#00FF00")), 
			choice)
	}

	if (m.selected != 0) {
		s += fmt.Sprintf("\nSelected: %v\n", m.selected)
	}

	// The footer
	s += fmt.Sprintf(
		"\nPress %s to refresh, %s to quit.\n", 
		termenv.String("R").Foreground(p.Color("#E88388")), 
		termenv.String("Q").Foreground(p.Color("#E88388")))

	return s
}

func checkServer() tea.Msg {
	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	res, err := c.Get(url)
	if err != nil {
		return errMsg{err}
	}

	y := make([]string, 0)
	json.NewDecoder(res.Body).Decode(&y)

	return payloadMsg(y)
}