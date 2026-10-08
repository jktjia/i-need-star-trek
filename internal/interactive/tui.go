package interactive

import (
	"errors"
	"fmt"
	"i-need-star-trek/internal/shared"
	"maps"
	"math/rand"
	"os"
	"slices"
	"strconv"
	"time"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

type styles struct {
	app           lipgloss.Style
	title         lipgloss.Style
	statusMessage lipgloss.Style
}

func newStyles(darkBG bool) styles {
	lightDark := lipgloss.LightDark(darkBG)

	return styles{
		app: lipgloss.NewStyle().Padding(1, 2),
		title: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1),
		statusMessage: lipgloss.NewStyle().
			Foreground(lightDark(lipgloss.Color("#04B575"), lipgloss.Color("#04B575"))),
	}
}

type model struct {
	width, height int

	styles styles
	darkBG bool

	tabs tabs

	epList   list.Model
	selected *episode

	filterForm *huh.Form
	minDate    *string
	maxDate    *string
	minLen     *string
	maxLen     *string
	shows      *[]string
	species    *[]string
	tags       *[]string

	random bool
	filter bool
}

var tabOptions = []string{"Episodes", "Filters", "Random"}

func newModel() model {
	zone.NewGlobal()

	var (
		m = model{}
	)

	m.styles = newStyles(false)

	m.tabs = tabs{id: "filters", items: tabOptions}
	m.tabs.active = "Episodes"

	items := []list.Item{}
	for _, e := range shared.Episodes {
		items = append(items, newEpisode(e))
	}

	delegate := newItemDelegate()
	episodeList := list.New(items, delegate, 0, 0)
	episodeList.SetShowTitle(false)
	episodeList.SetFilteringEnabled(false)
	episodeList.SetShowStatusBar(true)
	episodeList.KeyMap.Quit = key.NewBinding(key.WithKeys("q"),
		key.WithHelp("q", "quit"))
	m.epList = episodeList

	// showOptions := []huh.Option[string]{}
	// showList := slices.SortedFunc(maps.Values(*data.ShowMap), func(s1, s2 data.Show) int { return s1.StartDate.Compare(*s2.StartDate) })
	// for _, show := range showList {
	// 	showOptions = append(showOptions, huh.NewOption(show.FullName, show.Abbr))
	// }

	// tagOptions := []huh.Option[string]{}
	// for k, v := range *data.EventMap {
	// 	tagOptions = append(tagOptions, huh.NewOption(k, v))
	// }

	emptyStr := ""
	m.minDate = &emptyStr
	emptyStr2 := ""
	m.maxDate = &emptyStr2
	emptyStr3 := ""
	m.minLen = &emptyStr3
	emptyStr4 := ""
	m.maxLen = &emptyStr4
	emptySlice := []string{}
	m.shows = &emptySlice
	emptySlice2 := []string{}
	m.species = &emptySlice2
	emptySlice3 := []string{}
	m.tags = &emptySlice3

	m.newForm()

	return m
}

func (m *model) newForm() {
	showOptions := []huh.Option[string]{}
	showList := slices.SortedFunc(maps.Values(shared.ShowMap), func(s1, s2 shared.Show) int { return s1.StartDate.Compare(*s2.StartDate) })
	for _, show := range showList {
		showOptions = append(showOptions, huh.NewOption(show.FullName, show.Abbr))
	}

	tagOptions := []huh.Option[string]{}
	for k, v := range shared.EventMap {
		tagOptions = append(tagOptions, huh.NewOption(k, v))
	}

	m.filterForm = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("minDate").
				Placeholder("1966-09-08").
				Title("Start date").
				Validate(shared.CanParseDate).
				Value(m.minDate),

			huh.NewInput().
				Key("maxDate").
				Placeholder("2026-09-25").
				Title("End date").
				Validate(shared.CanParseDate).
				Value(m.maxDate),

			huh.NewInput().
				Key("minLen").
				Title("Min Length (minutes)").
				Validate(func(s string) error {
					if s != "" {
						i, err := strconv.Atoi(s)
						if err != nil {
							return errors.New("min length must be an integer")
						} else if i < 0 {
							return errors.New("min length must be positive")
						}
					}
					return nil
				}).
				Value(m.minLen),

			huh.NewInput().
				Key("maxLen").
				Title("Max Length (minutes)").
				Validate(func(s string) error {
					if s != "" {
						i, err := strconv.Atoi(s)
						if err != nil {
							return errors.New("max length must be an integer")
						} else {
							min, err := strconv.Atoi(*m.minLen)
							if err == nil && i < min {
								return errors.New("max length must not be less than min length")
							}
						}
					}
					return nil
				}).
				Value(m.maxLen),

			huh.NewMultiSelect[string]().
				Key("shows").
				Options(showOptions...).
				Title("Shows to include:").
				Height(8).
				Value(m.shows),
		),

		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Key("species").
				Options(huh.NewOptions(shared.Species...)...).
				Title("Species in the episodes:").
				Height(8).
				Value(m.species),

			huh.NewMultiSelect[string]().
				Key("tags").
				Options(tagOptions...).
				Title("Tags for the episodes:").
				Height(8).
				Value(m.tags),
		),
	).WithAccessible(true).WithLayout(huh.LayoutColumns(2)).WithShowErrors(true).WithShowHelp(true)

	m.filterForm.Init()
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) isInitialized() bool {
	return m.height != 0 && m.width != 0
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+e":
			zone.SetEnabled(!zone.Enabled())
			return m, nil
		case "ctrl+c":
			return m, tea.Quit
		case "shift+left":
			idx := slices.Index(tabOptions, m.tabs.active)
			if idx > 0 {
				m.tabs.active = tabOptions[idx-1]
			}
		case "shift+right":
			idx := slices.Index(tabOptions, m.tabs.active)
			if idx < len(tabOptions)-1 {
				m.tabs.active = tabOptions[idx+1]
			}
		}

	case tea.BackgroundColorMsg:
		m.darkBG = msg.IsDark()
		m.updateListProperties()

	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width

	}

	cmds = append(cmds, m.tabs.Update(msg))

	switch m.tabs.active {
	case "Episodes":
		m.random = false
		m.filter = false
		m.updateListProperties()

		list, cmd := m.epList.Update(msg)
		m.epList = list
		cmds = append(cmds, cmd)

	case "Filters":
		m.random = false
		m.updateFormProperties()

		if !m.filter {
			m.resetForm()
			m.filter = true
		}

		form, cmd := m.filterForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.filterForm = f
			cmds = append(cmds, cmd)
			if f.State == huh.StateCompleted {
				m.tabs.active = "Episodes"
			}
			m.applyFilters()
		}

	case "Random":
		m.filter = false

		if !m.random {
			m.selectRandom()
			m.random = true
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *model) updateListProperties() {
	// Update list size.
	h, v := m.styles.app.GetFrameSize()
	m.epList.SetSize(m.width-h, m.height-v-m.tabs.GetHeight())

	// Update the model and list styles.
	m.styles = newStyles(m.darkBG)
	m.epList.Styles.Title = m.styles.title
}

func (m *model) updateFormProperties() {
	// Update form size and theme.
	h, v := m.styles.app.GetFrameSize()
	m.filterForm = m.filterForm.WithHeight(m.height - v - m.tabs.GetHeight()).WithWidth(m.width - h)
}

func (m *model) resetForm() {
	m.newForm()
	m.updateFormProperties()
}

func (m *model) applyFilters() {
	minDateStr := *m.minDate
	var minDate *time.Time
	if minDateStr != "" {
		d, _ := shared.ParseDate(minDateStr)
		minDate = &d
	}

	maxDateStr := *m.maxDate
	var maxDate *time.Time
	if maxDateStr != "" {
		d, _ := shared.ParseDate(maxDateStr)
		maxDate = &d
	}

	minLenStr := *m.minLen
	var minLen *int
	if minLenStr != "" {
		i, _ := strconv.Atoi(minLenStr)
		minLen = &i
	}

	maxLenStr := *m.maxLen
	var maxLen *int
	if maxLenStr != "" {
		i, _ := strconv.Atoi(maxLenStr)
		maxLen = &i
	}

	eps := shared.FilterEpisodes(shared.Episodes, *m.shows, *m.species, *m.tags, minDate, maxDate, minLen, maxLen)

	items := []list.Item{}
	for _, e := range eps {
		items = append(items, newEpisode(e))
	}

	delegate := newItemDelegate()
	episodeList := list.New(items, delegate, 0, 0)
	episodeList.SetShowTitle(false)
	episodeList.SetFilteringEnabled(false)
	episodeList.SetShowStatusBar(true)
	episodeList.KeyMap.Quit = key.NewBinding(key.WithKeys("q"),
		key.WithHelp("q", "quit"))
	m.epList = episodeList

	m.updateListProperties()
}

func (m *model) selectRandom() {
	idx := rand.Intn(len(m.epList.Items()))
	m.epList.Select(idx)
}

func (m model) View() tea.View {
	var view tea.View
	view.AltScreen = true
	view.MouseMode = tea.MouseModeCellMotion

	if !m.isInitialized() {
		return view
	}

	s := lipgloss.NewStyle().MaxHeight(m.height).MaxWidth(m.width)

	epView := ""
	ep, ok := m.epList.SelectedItem().(episode)
	if ok {
		epView = ep.View(0)
	}

	content := ""
	switch m.tabs.active {
	case "Episodes":
		content = lipgloss.JoinHorizontal(
			lipgloss.Top,
			m.styles.app.Width(m.width/2-2).
				Render(m.epList.View()), m.styles.app.
				Width(m.width/2).
				Border(lipgloss.NormalBorder()).Render(epView),
		)
	case "Filters":
		content = m.styles.app.Width(m.width).Render(m.filterForm.View())
	case "Random":
		content = m.styles.app.Width(m.width).Render(epView)
	}

	// Wrap the main models view in [zone.Scan].
	view.SetContent(zone.Scan(s.Render(
		lipgloss.JoinVertical(lipgloss.Top,
			lipgloss.NewStyle().MarginBottom(1).Render(m.tabs.View()),
			lipgloss.PlaceHorizontal(
				m.width, lipgloss.Left,
				content,
				lipgloss.WithWhitespaceChars(" "),
			),
		),
	)))

	return view
}

func Start() {
	p := tea.NewProgram(newModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("there's been an error: %v", err)
		os.Exit(1)
	}
}
