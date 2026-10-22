package shared

import (
	"slices"
	"time"
)

var (
	DateFormats = []string{"2006", "Jan 2006", "January 2006", "1 Jan 2006", "Jan 1 2006", "1 January 2006", "January 1 2006", time.DateOnly}

	Episodes = *importEpisodes(importSpecies(false), importEventTypes(false), false)
	EventMap = *importEventTypes(false)
	Species  = *importSpecies(false)
	ShowMap  = *importShows(false)
)

func FilterEpisodes(eps []EpisodeInfo, shows []string, species []string, other []string, minDate *time.Time, maxDate *time.Time, minLen *int, maxLen *int) []EpisodeInfo {
	filtered := []EpisodeInfo{}
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
