package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	spinner "github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

func waitForRooms(r chan roomChanMsg) tea.Cmd {
	return func() tea.Msg {
		return roomChanMsg(<-r)
	}
}

func listRooms(r chan roomChanMsg) tea.Cmd {
	return func() tea.Msg {
		for {
			c := &http.Client{
				Timeout: 10 * time.Second,
			}
			res, err := c.Get(url)
			if err != nil {
				return errMsg{err}
			}

			y := make([]string, 0)
			json.NewDecoder(res.Body).Decode(&y)

			r <- roomChanMsg(y)
		}
	}
}

// Control.
func controlRoom(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "up", "k":
				if m.cursor["rooms"] > 0 {
						m.cursor["rooms"]--
				}

			case "down", "j":
				if m.cursor["rooms"] < len(m.rooms)-1 {
						m.cursor["rooms"]++
				}

			case "enter", " ":
				m.selected = m.cursor["rooms"] + 1
				return m, listUsers(m.u)

			default:
				return m, nil
			}

		case spinner.TickMsg:
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd

		case roomChanMsg:
			m.rooms = []string(msg)
			return m, waitForRooms(m.r)
	}
	
	return m, nil
}

// View.
func roomsView(m model) string {
	p := termenv.EnvColorProfile()
	
	var s string
	if (m.rooms == nil) {
		s += fmt.Sprintf("\n %s Fetching URL: %v\n\n", m.spinner.View(), url)
	} else {
		s += fmt.Sprintf("%s", termenv.String("\nResults\n\n").Bold())
	}

	for i, choice := range m.rooms {
		cursor := " " 
		if m.cursor["rooms"] == i {
				cursor = "→" 
		}

		s += fmt.Sprintf(
			"%s %s\n", 
			termenv.String(cursor).Foreground(p.Color("#00FF00")), 
			choice)
	}

	return s
}