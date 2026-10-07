package data

// import (
// 	"encoding/json"
// 	"fmt"
// 	"i-need-star-trek/internal/consts"
// 	"os"
// 	"slices"
// 	"strings"
// 	"time"

// 	"github.com/charmbracelet/bubbles/help"
// 	"github.com/charmbracelet/bubbles/key"
// )

// type show struct {
// 	FullName  string     `json:"full"`
// 	Abbr      string     `json:"abbr"`
// 	Seasons   int        `json:"seasons"`
// 	StartDate *time.Time `json:"start"`
// 	EndDate   *time.Time `json:"end"`
// }

// func getShows() ([]show, map[string]show) {
// 	data, err := os.ReadFile(showJSON)
// 	if err != nil {
// 		panic(err)
// 	}

// 	var shows []show
// 	json.Unmarshal(data, &shows)

// 	slices.SortFunc(shows, func(s1, s2 show) int {
// 		return s1.StartDate.Compare(*s2.StartDate)
// 	})

// 	showMap := make(map[string]show)

// 	for _, s := range shows {
// 		showMap[s.Abbr] = s
// 	}

// 	return shows, showMap
// }

// func getSpecies() []string {
// 	data, err := os.ReadFile(speciesJSON)
// 	if err != nil {
// 		panic(err)
// 	}

// 	var species []string
// 	json.Unmarshal(data, &species)
// 	return species
// }

// type event struct {
// 	FullName string `json:"fullName"`
// 	Short    string `json:"shortened"`
// }

// func getEvents() []event {
// 	data, err := os.ReadFile(eventTypeJSON)
// 	if err != nil {
// 		panic(err)
// 	}

// 	var events []event
// 	json.Unmarshal(data, &events)
// 	return events
// }

// var (
// 	Shows, showMap = getShows()
// 	Species        = getSpecies()
// 	Events         = getEvents()
// )

// type episode struct {
// 	Name          string          `json:"name"`
// 	Show          string          `json:"show"`
// 	Season        int             `json:"season"`
// 	EpNumber      int             `json:"episode"`
// 	EpNumber2     *int            `json:"episode2"`
// 	ProdNumber    string          `json:"production"`
// 	ReleaseDate   *time.Time      `json:"release"`
// 	StardateStart *string         `json:"start"`
// 	StardateEnd   *string         `json:"end"`
// 	Directors     *[]string       `json:"directors"`
// 	Writers       *[]string       `json:"writers"`
// 	Length        int             `json:"length"`
// 	Crossover     *[]string       `json:"crossover"`
// 	MainChars     *[]string       `json:"characters"`
// 	Species       map[string]bool `json:"species"`
// 	EpCategory    map[string]bool `json:"category"`
// 	Summary       string          `json:"summary"`
// }

// func (i episode) Title() string { return i.Name }
// func (i episode) Description() string {
// 	return fmt.Sprintf("%s, %s", i.fullShowName(), i.seasonEpisode())
// }
// func (i episode) fullShowName() string {
// 	return showMap[i.Show].FullName
// }
// func (i episode) seasonEpisode() string {
// 	if i.Season == 0 && (i.EpNumber == 0 || i.EpNumber == 1) {
// 		return "Pilot"
// 	} else if i.EpNumber2 != nil {
// 		return fmt.Sprintf("Season %d, Episodes %d and %d", i.Season, i.EpNumber, *i.EpNumber2)
// 	} else {
// 		return fmt.Sprintf("Season %d, Episode %d", i.Season, i.EpNumber)
// 	}
// }
// func (e episode) FilterValue() string { return e.Name }
// func (e episode) View(width int) string {
// 	showEp := fmt.Sprintf("%s\n%s", e.fullShowName(), e.seasonEpisode())
// 	directors := strings.Join(*e.Directors, ", ")
// 	if len(*e.Directors) == 1 {
// 		directors = "Director:  " + directors
// 	} else {
// 		directors = "Directors: " + directors
// 	}
// 	writers := strings.Join(*e.Writers, ", ")
// 	if len(*e.Writers) == 1 {
// 		writers = "Writer:    " + writers
// 	} else {
// 		writers = "Writers:   " + writers
// 	}
// 	summary := []string{""}
// 	for _, s := range strings.Split(e.Summary, " ") {
// 		last := len(summary) - 1
// 		if len(summary[last])+len(s) > width {
// 			summary = append(summary, s+" ")
// 		} else {
// 			summary[last] = summary[last] + s + " "
// 		}
// 	}
// 	summaryStr := strings.Join(summary, "\n")
// 	base := fmt.Sprintf("%s\n%s\n\n%s\n\n%s\n%s\n\n%s", e.headerView(), showEp, summaryStr, directors, writers, e.footerView())
// 	return base
// }

// func (e episode) headerView() string {
// 	title := consts.TitleStyle.Padding(0, 1).Render(e.Name)
// 	length := consts.DescStyle.Padding(1, 0).Render(fmt.Sprintf("%d minutes", e.Length))
// 	return fmt.Sprintf("%s\n%s", title, length)
// }

// func (e episode) footerView() string {
// 	return help.New().View(e)
// }

// func (e episode) ShortHelp() []key.Binding {
// 	return []key.Binding{key.NewBinding(
// 		key.WithKeys("esc"),
// 		key.WithHelp("esc", "go back"),
// 	)}
// }

// func (e episode) FullHelp() [][]key.Binding {
// 	return [][]key.Binding{e.ShortHelp()}
// }

// func getEpisodes() []episode {
// 	data, err := os.ReadFile(episodeJSON)
// 	if err != nil {
// 		panic(err)
// 	}

// 	var episodes []episode
// 	json.Unmarshal(data, &episodes)

// 	return episodes
// }

// var Episodes = getEpisodes()
