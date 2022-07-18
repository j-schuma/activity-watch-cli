package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	"math/rand"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hako/durafmt"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type activity struct {
	title    string
	duration time.Duration
}

func (a activity) Title() string       { return a.title }
func (a activity) Description() string { return durafmt.Parse(a.duration).String() }
func (a activity) FilterValue() string { return a.title }

type model struct {
	list  list.Model
	input textinput.Model
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		// when we press entter and the textbox has focus we want to emmit a message with the value
		// https://github.com/caarlos0/tasktimer/blob/main/internal/ui/main.go#L147
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)
	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return secondaryForeground.Render("project: ") +
		activeForegroundBold.Render("this could be a project") +
		separator + "I am a dummy timer" + "\n\n" +
		m.input.View() + "\n\n" +
		m.list.View() + "\n"
}

func main() {
	rand.Seed(time.Now().UnixNano())

	activities := []list.Item{
		activity{title: "do stuff", duration: 2 * time.Minute},
		activity{title: "do more stuff", duration: 5 * time.Hour},
		activity{title: "even more stuff", duration: 45 * time.Minute},
	}

	input := textinput.New()
	input.Prompt = "❯ "
	input.Placeholder = "New task description..."
	input.Focus()
	input.CharLimit = 250
	input.Width = 50

	m := model{list: list.New(activities, list.NewDefaultDelegate(), 0, 0), input: input}
	m.list.Title = "What are you doing?"

	p := tea.NewProgram(m, tea.WithAltScreen())

	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
