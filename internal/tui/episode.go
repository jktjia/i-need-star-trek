package tui

import (
	"fmt"
	"i-need-star-trek/internal/data"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

type episode struct {
	Name          string          `json:"name"`
	Show          string          `json:"show"`
	Season        int             `json:"season"`
	EpNumber      int             `json:"episode"`
	EpNumber2     *int            `json:"episode2"`
	ProdNumber    string          `json:"production"`
	ReleaseDate   *time.Time      `json:"release"`
	StardateStart *string         `json:"start"`
	StardateEnd   *string         `json:"end"`
	Directors     *[]string       `json:"directors"`
	Writers       *[]string       `json:"writers"`
	Length        int             `json:"length"`
	Crossover     *[]string       `json:"crossover"`
	MainChars     *[]string       `json:"characters"`
	Species       map[string]bool `json:"species"`
	EpCategory    map[string]bool `json:"category"`
	Summary       string          `json:"summary"`
}

func (i episode) Title() string { return i.Name }
func (i episode) Description() string {
	return fmt.Sprintf("%s, %s", i.fullShowName(), i.seasonEpisode())
}
func (i episode) fullShowName() string {
	return (*data.ShowMap)[i.Show].FullName
}
func (i episode) seasonEpisode() string {
	if i.Season == 0 && (i.EpNumber == 0 || i.EpNumber == 1) {
		return "Pilot"
	} else if i.EpNumber2 != nil {
		return fmt.Sprintf("Season %d, Episodes %d and %d", i.Season, i.EpNumber, *i.EpNumber2)
	} else {
		return fmt.Sprintf("Season %d, Episode %d", i.Season, i.EpNumber)
	}
}
func (e episode) FilterValue() string { return e.Name }
func (e episode) View(width int) string {
	base := fmt.Sprintf("%s\n%s\n%s", e.headerView(), e.mainView(width), e.footerView())
	return base
}

func (e episode) mainView(width int) string {
	showEp := fmt.Sprintf("%s\n%s", e.fullShowName(), e.seasonEpisode())
	directors := strings.Join(*e.Directors, ", ")
	if len(*e.Directors) == 1 {
		directors = "Director:  " + directors
	} else {
		directors = "Directors: " + directors
	}
	writers := strings.Join(*e.Writers, ", ")
	if len(*e.Writers) == 1 {
		writers = "Writer:    " + writers
	} else {
		writers = "Writers:   " + writers
	}
	summary := []string{""}
	for _, s := range strings.Split(e.Summary, " ") {
		last := len(summary) - 1
		if len(summary[last])+len(s) > width {
			summary = append(summary, s+" ")
		} else {
			summary[last] = summary[last] + s + " "
		}
	}
	summaryStr := strings.Join(summary, "\n")
	return fmt.Sprintf("%s\n\n%s\n\n%s\n%s\n", showEp, summaryStr, directors, writers)
}

func (e episode) headerView() string {
	title := titleStyle.Padding(0, 1).Render(e.Name)
	length := descStyle.Padding(1, 0).Render(fmt.Sprintf("%d minutes", e.Length))
	return fmt.Sprintf("%s\n%s", title, length)
}

func (e episode) footerView() string {
	return help.New().View(e)
}

func (e episode) ShortHelp() []key.Binding {
	return []key.Binding{key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "go back"),
	)}
}

func (e episode) FullHelp() [][]key.Binding {
	return [][]key.Binding{e.ShortHelp()}
}

func newEpisode(ep data.EpisodeInfo) episode {
	return episode(ep)
}

func ViewEpisode(ep data.EpisodeInfo) string {
	e := newEpisode(ep)
	return fmt.Sprintf("\n%s\n%s\n", e.headerView(), e.mainView(120))
}
