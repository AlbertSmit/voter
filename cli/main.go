package main

// A simple program that makes a GET request and prints the response status.

import (
	"log"

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
		if k == "q" || k == "ctrl+c" {
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

func (m model) View() string {
	var s string

	if (m.selected != 0) {
		s = roomView(m)
	} else {
		s = roomsView(m)
	}

	return indent.String("\n"+s+"\n\n", 2)
}