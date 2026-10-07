package shared

import (
	"errors"
	"i-need-star-trek/internal/data"
	"slices"
	"strings"
	"time"
)

func FilterEpisodes(eps []data.EpisodeInfo, shows []string, species []string, other []string, minDate *time.Time, maxDate *time.Time, minLen *int, maxLen *int) []data.EpisodeInfo {
	filtered := []data.EpisodeInfo{}
	for _, ep := range eps {
		showViable := true
		if len(shows) > 0 {
			showViable = slices.Contains(shows, ep.Show)
		}
		speciesViable := true
		for _, s := range species {
			speciesViable = speciesViable && ep.Species[s]
		}
		eventViable := true
		for _, e := range other {
			eventViable = eventViable && ep.EpCategory[e]
		}
		viable := showViable && speciesViable && eventViable
		if minDate != nil && ep.ReleaseDate != nil {
			viable = viable && ep.ReleaseDate.Compare(*minDate) > 0
		}
		if maxDate != nil && ep.ReleaseDate != nil {
			viable = viable && ep.ReleaseDate.Compare(*maxDate) < 0
		}
		if minLen != nil {
			viable = viable && ep.Length >= *minLen
		}
		if maxLen != nil {
			viable = viable && ep.Length <= *maxLen
		}
		if viable {
			filtered = append(filtered, ep)
		}
	}
	return filtered
}

func ParseAllStr(strs []string, parseFunc func(s string) (string, error)) ([]string, error) {
	parsed := []string{}
	for i := range strs {
		str, err := parseFunc(strs[i])
		if err != nil {
			return nil, err
		} else {
			parsed = append(parsed, str)
		}
	}
	return parsed, nil
}

func ParseAllShows(shows []string) ([]string, error) {
	return ParseAllStr(shows, ParseShow)
}

func ParseShow(show string) (string, error) {
	showMap := *data.ShowMap

	for k, v := range showMap {
		if strings.EqualFold(show, v.Abbr) || strings.EqualFold(show, v.FullName) || strings.EqualFold("Star Trek: "+show, v.FullName) {
			return k, nil
		}
	}

	return "", errors.New("Unknown show: " + show)
}

func ParseAllSpecies(species []string) ([]string, error) {
	return ParseAllStr(species, ParseSpecies)
}

func ParseSpecies(species string) (string, error) {
	knownSpecies := *data.Species

	for _, sp := range knownSpecies {
		if strings.EqualFold(species, sp) || strings.EqualFold(species, sp+"s") {
			return sp, nil
		}
	}
	return "", errors.New("Unknown species: " + species)
}

func ParseAllOther(other []string) ([]string, error) {
	return ParseAllStr(other, ParseOther)
}

func ParseOther(other string) (string, error) {
	eventMap := *data.EventMap

	for k, v := range eventMap {
		if strings.EqualFold(other, k) {
			return v, nil
		}
	}
	return "", errors.New("Unknown event type: " + other)
}

func ParseDate(dateStr string) (time.Time, error) {
	for _, format := range DateFormats {
		date, err := time.Parse(format, dateStr)
		if err == nil {
			return date, err
		}
	}
	return time.Now(), errors.New("Invalid time format")
}

func CanParseDate(dateStr string) error {
	if dateStr == "" {
		return nil
	}
	_, err := ParseDate(dateStr)
	return err
}

var (
	DateFormats = []string{"2006", "Jan 2006", "January 2006", "1 Jan 2006", "Jan 1 2006", "1 January 2006", "January 1 2006", time.DateOnly}
)
