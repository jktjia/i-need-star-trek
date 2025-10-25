package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	accent = lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}

	statusMessageStyle = lipgloss.NewStyle().Foreground(accent).Render

	selectedStyle = lipgloss.NewStyle().Foreground(accent).BorderLeftForeground(accent)
	filterStyle   = lipgloss.NewStyle().Padding(1, 1)

	descStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: "#B2B2B2",
		Dark:  "#4A4A4A",
	})

	showList, showMap = getShows()
	speciesList       = getSpecies()
	eventList         = getEvents()
)

type keyMap struct {
	goBack         key.Binding
	toggleHelpMenu key.Binding
	showFilters    key.Binding
}

func newKeyMap() *keyMap {
	return &keyMap{
		goBack: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "go back"),
		),
		toggleHelpMenu: key.NewBinding(
			key.WithKeys("H"),
			key.WithHelp("H", "toggle help"),
		),
		showFilters: key.NewBinding(
			key.WithKeys("f"),
			key.WithHelp("f", "show filters"),
		),
	}
}

type selectInfo struct{ ep *episode }

type model struct {
	episodes    []episode
	keys        *keyMap
	list        episodelist
	selected    *selectInfo
	filters     filterModel
	showFilters bool
	width       int
}

func newModel() model {
	var (
		keys        = newKeyMap()
		episodes    = getEpisodes()
		selected    = selectInfo{}
		filtersBase = newFilters()
	)

	episodeList := newEpisodeList(episodes, keys, &selected)

	return model{
		list:        episodeList,
		keys:        keys,
		episodes:    episodes,
		selected:    &selected,
		filters:     filtersBase,
		showFilters: false,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var (
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width - appStyle.GetPaddingLeft() - appStyle.GetPaddingRight()
		newFilters, cmd := m.filters.Update(msg)
		m.filters = newFilters
		cmds = append(cmds, cmd)
		cmd = m.list.SetEpisodes(m.filters.selectedFilters.filterEpisodes(m.episodes))
		cmds = append(cmds, cmd)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.goBack):
			m.selected.ep = nil
			m.showFilters = false
			return m, nil

		case key.Matches(msg, m.keys.showFilters):
			m.showFilters = true
			return m, nil
		}
	}

	if m.showFilters {
		newFilters, cmd := m.filters.Update(msg)
		m.filters = newFilters
		cmds = append(cmds, cmd)
		cmd = m.list.SetEpisodes(m.filters.selectedFilters.filterEpisodes(m.episodes))
		cmds = append(cmds, cmd)
	} else {
		newList, cmd := m.list.Update(msg)
		m.list = newList
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var content string
	if m.selected.ep != nil {
		content = (*m.selected.ep).View(m.width)
	} else if m.showFilters {
		content = m.filters.View()
	} else {
		content = m.list.View()
	}
	return appStyle.Render(content)
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	if _, err := tea.NewProgram(newModel(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
