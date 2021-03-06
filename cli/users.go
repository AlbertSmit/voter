package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	spinner "github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

// User expected from room.
type User struct {
	Name		string `json:"name"`
	UUID		string `json:"uuid"`
	role		int 
}

func waitForUsers(u chan usersChanMsg) tea.Cmd {
	return func() tea.Msg {
		return usersChanMsg(<-u)
	}
}

func listUsers(u chan usersChanMsg) tea.Cmd {
	return func() tea.Msg {
		for {
			c := &http.Client{
				Timeout: 10 * time.Second,
			}
			res, err := c.Get(users)
			if err != nil {
				log.Fatal(err)
			}

			users := make(usersChanMsg, 0)
			json.NewDecoder(res.Body).Decode(&users)

			u <- usersChanMsg(users)
		}
	}
}

// View.
func roomView(m model) string {
	p := termenv.EnvColorProfile()
	room := m.rooms[m.selected - 1]

	var s string
	if (len(m.users) == 0) {
		s += fmt.Sprintf("\n %s Getting users\n\n", m.spinner.View())
	} else {
		s += fmt.Sprintf(
			"\n%s %s\n\n", 
			termenv.String("Users in").Bold(), 
			termenv.String(room).Bold())
	}

	for i, user := range m.users {
		cursor := " " 
		if m.cursor["users"] == i {
				cursor = "→" 
		}

		s += fmt.Sprintf(
			"%s %v\n", 
			termenv.String(cursor).Foreground(p.Color("#00FF00")), 
			strings.Title(user.Name))
	}

	s += termenv.String("\n\npress ESC to return\npress  D  to kick user").Faint().String()

	return s
}

// Control.
func controlUsers(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc":
				m.selected = 0

			case "up", "k":
				if m.cursor["users"] > 0 {
						m.cursor["users"]--
				}

			case "down", "j":
				if m.cursor["users"] < len(m.rooms)-1 {
						m.cursor["users"]++
				}

			case "D":
				log.Println("Trying to delete %i", m.users[m.cursor["users"]].Name)

			default:
				return m, nil
			}

		case spinner.TickMsg:
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd

		case usersChanMsg:
			m.users = []User(msg)
			return m, waitForUsers(m.u)
	}
	
	return m, nil
}

