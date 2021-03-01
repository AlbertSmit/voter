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

const url = "https://api.mocki.io/v1/d201c428"

type model struct {
	rooms			[]string
	
	cursor		int
	selected	int
	
	sub 			chan payloadMsg
	spinner   spinner.Model
	status 		int
	err    		error
}

type statusMsg int
type payloadMsg []string
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
		sub:     	make(chan payloadMsg),
		spinner: 	s,
		selected: 0,
		cursor: 	0,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		spinner.Tick,
		listRooms(m.sub),
		waitForFetch(m.sub),
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

	switch msg := msg.(type) {
	// Move this; see https://github.com/charmbracelet/bubbletea/blob/df0da429545895356259fecf23ea35fb9c938b61/examples/views/main.go#L265
	case tea.KeyMsg:
		switch msg.String() {
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
		return m, waitForFetch(m.sub)

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
	var s string

	if len(m.rooms) == 0 {
		s = fetchingView(m)
	} else if (m.selected != 0) {
		s = roomView(m)
	} else {
		s = roomsView(m)
	}

	return indent.String("\n"+s+"\n\n", 2)
}

func fetchingView(m model) string {
	var s string
	if (len(m.rooms) == 0) {
		s += fmt.Sprintf("\n %s Fetching URL: %v\n\n", m.spinner.View(), url)
	} else {
		s += fmt.Sprintf("%s", termenv.String("\nResults\n\n").Bold())
	}

	return s
}

// User expected from room.
type User struct {
	Name		string `json:"name"`
	UUID		string `json:"uuid"`
	role		int 
}

var users = "https://api.mocki.io/v1/aeeb90f6"

func roomView(m model) string {
	p := termenv.EnvColorProfile()
	room := m.rooms[m.selected - 1]

	var s string
	e := fmt.Sprint(room)
	s += fmt.Sprintf(
		"Users in %s.\n\n", 
		termenv.String(e).Bold())
		
	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	res, err := c.Get(users)
	if err != nil {
		log.Fatal(err)
	}

	users := make([]User, 0)
	json.NewDecoder(res.Body).Decode(&users)

	for i, user := range users {
		cursor := " " 
		if m.cursor == i {
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

	return s
}

func waitForFetch(sub chan payloadMsg) tea.Cmd {
	return func() tea.Msg {
		return payloadMsg(<-sub)
	}
}

func listRooms(sub chan payloadMsg) tea.Cmd {
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

			sub <- payloadMsg(y)
		}
	}
}