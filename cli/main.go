package main

// A simple program that makes a GET request and prints the response status.

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	spinner "github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/indent"
	"github.com/muesli/termenv"
)

const (
	url = 		"https://api.mocki.io/v1/d201c428"
	users = 	"https://api.mocki.io/v1/aeeb90f6"
)

type model struct {
	rooms			[]string
	users			[]User
	
	cursor		map[string]int
	selected	int
	
	r 				chan roomChanMsg
	u 				chan usersChanMsg

	spinner   spinner.Model
	status 		int
	err    		error
}

type statusMsg int
type roomChanMsg []string
type usersChanMsg []User
type errMsg struct{ error }

func (e errMsg) Error() string { return e.Error() }

func main() {
	p := tea.NewProgram(initialModel())

	termenv.ClearScreen()
	p.EnterAltScreen()

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}

func initialModel() model {
	s := spinner.NewModel()
	s.Spinner = spinner.Dot
	return model{
		r:     		make(chan roomChanMsg),
		u:     		make(chan usersChanMsg),
		spinner: 	s,
		selected: 0,
		cursor: 	make(map[string]int, 0),
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		spinner.Tick,
		listRooms(m.r),
		waitForRooms(m.r),
		waitForUsers(m.u),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Make sure these keys always quit
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			return m, tea.Quit
		}
	}

	// Room controls
	if m.selected == 0 {
		return controlRoom(msg, m)
	}

	// User controls
	if m.selected != 0 {
		return controlUsers(msg, m)
	} 

	return m, nil
}

// Update loop for the first view where you're choosing a room.
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

// Update loop for the second view where you're controlling users.
func controlUsers(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "up", "k":
				if m.cursor["users"] > 0 {
						m.cursor["users"]--
				}

			case "down", "j":
				if m.cursor["users"] < len(m.rooms)-1 {
						m.cursor["users"]++
				}

			case "enter", " ":
				m.selected = m.cursor["users"] + 1

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

func (m model) View() string {
	var s string

	if (m.selected != 0) {
		s = roomView(m)
	} else {
		s = roomsView(m)
	}

	return indent.String("\n"+s+"\n\n", 2)
}

// User expected from room.
type User struct {
	Name		string `json:"name"`
	UUID		string `json:"uuid"`
	role		int 
}

func roomView(m model) string {
	p := termenv.EnvColorProfile()
	room := m.rooms[m.selected - 1]

	var s string
	if (len(m.users) == 0) {
		s += fmt.Sprintf("\n %s Getting users\n\n", m.spinner.View())
	} else {
		s += fmt.Sprintf("\n%s %s\n\n", termenv.String("Users in").Bold(), termenv.String(room).Bold())
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

	return s
}

func roomsView(m model) string {
	p := termenv.EnvColorProfile()
	
	var s string
	if (len(m.rooms) == 0) {
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

func waitForRooms(r chan roomChanMsg) tea.Cmd {
	return func() tea.Msg {
		return roomChanMsg(<-r)
	}
}

func waitForUsers(u chan usersChanMsg) tea.Cmd {
	return func() tea.Msg {
		return usersChanMsg(<-u)
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