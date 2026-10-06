package internal

import (
	"errors"
	"i-need-star-trek/internal/data"
	"slices"
	"strings"
)

func filterEpisodes(eps []data.EpisodeInfo, shows []string, species []string, other []string) []data.EpisodeInfo {
	filtered := []data.EpisodeInfo{}
	for _, ep := range eps {
		showViable := true
		if len(shows) > 0 {
			showViable = slices.Contains(shows, ep.Show)
		}
		speciesViable := !(len(species) > 0)
		for _, s := range species {
			speciesViable = speciesViable || ep.Species[s]
		}
		eventViable := true
		for _, e := range other {
			eventViable = eventViable && ep.EpCategory[e]
		}
		if showViable && speciesViable && eventViable {
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
	uppercase := strings.ToUpper(show)
	showMap := *data.ShowMap

	_, ok := showMap[uppercase]
	if ok {
		return uppercase, nil
	} else {
		for k, v := range showMap {
			if strings.EqualFold(show, v) || strings.EqualFold("Star Trek: "+show, v) {
				return k, nil
			}
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
		if strings.EqualFold(species, sp) {
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
