package internal

import (
	"errors"
	"i-need-star-trek/internal/data"
	"slices"
	"strings"
	"time"
)

func filterEpisodes(eps []data.EpisodeInfo, shows []string, species []string, other []string, before *time.Time, after *time.Time) []data.EpisodeInfo {
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
		if before != nil && ep.ReleaseDate != nil {
			viable = viable && ep.ReleaseDate.Compare(*before) < 0
		}
		if after != nil && ep.ReleaseDate != nil {
			viable = viable && ep.ReleaseDate.Compare(*after) > 0
		}
		if viable {
			filtered = append(filtered, ep)
		}
	}
	return filtered
}

func parseAllStr(strs []string, parseFunc func(s string) (string, error)) ([]string, error) {
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

func parseAllShows(shows []string) ([]string, error) {
	return parseAllStr(shows, parseShow)
}

func parseShow(show string) (string, error) {
	showMap := *data.ShowMap

	for k, v := range showMap {
		if strings.EqualFold(show, v.Abbr) || strings.EqualFold(show, v.FullName) || strings.EqualFold("Star Trek: "+show, v.FullName) {
			return k, nil
		}
	}

	return "", errors.New("Unknown show: " + show)
}

func parseAllSpecies(species []string) ([]string, error) {
	return parseAllStr(species, parseSpecies)
}

func parseSpecies(species string) (string, error) {
	knownSpecies := *data.Species

	for _, sp := range knownSpecies {
		if strings.EqualFold(species, sp) || strings.EqualFold(species, sp+"s") {
			return sp, nil
		}
	}
	return "", errors.New("Unknown species: " + species)
}

func parseAllOther(other []string) ([]string, error) {
	return parseAllStr(other, parseOther)
}

func parseOther(other string) (string, error) {
	eventMap := *data.EventMap

	for k, v := range eventMap {
		if strings.EqualFold(other, k) {
			return v, nil
		}
	}
	return "", errors.New("Unknown event type: " + other)
}
