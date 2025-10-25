package main

import (
	"slices"
)

var filterTabs = []string{"Show", "Species", "Other"}

type filters struct {
	species *[]string
	events  *[]string
	shows   *[]string
}

func (f filters) filterEpisodes(eps []episode) []episode {
	filtered := []episode{}
	for _, ep := range eps {
		showViable := true
		if len(*f.shows) > 0 {
			showViable = slices.Contains(*f.shows, ep.Show)
		}
		speciesViable := !(len(*f.species) > 0)
		for _, s := range *f.species {
			speciesViable = speciesViable || ep.Species[s]
		}
		eventViable := true
		for _, e := range *f.events {
			eventViable = eventViable && ep.EpCategory[e]
		}
		if showViable && speciesViable && eventViable {
			filtered = append(filtered, ep)
		}
	}
	return filtered
}

func (f *filters) addShow(show show) {
	*f.shows = append(*f.shows, show.Abbr)
}

func (f *filters) removeShow(show show) {
	newShows := []string{}
	for _, s := range *f.shows {
		if s != show.Abbr {
			newShows = append(newShows, s)
		}
	}
	f.shows = &newShows
}

func (f *filters) toggleShow(show show) {
	if slices.Contains(*f.shows, show.Abbr) {
		f.removeShow(show)
	} else {
		f.addShow(show)
	}
}

func (f *filters) addSpecies(species string) {
	*f.species = append(*f.species, species)
}

func (f *filters) removeSpecies(species string) {
	newSpecies := []string{}
	for _, s := range *f.species {
		if s != species {
			newSpecies = append(newSpecies, s)
		}
	}
	f.species = &newSpecies
}

func (f *filters) toggleSpecies(species string) {
	if slices.Contains(*f.species, species) {
		f.removeSpecies(species)
	} else {
		f.addSpecies(species)
	}
}

func (f *filters) addEvent(event event) {
	*f.events = append(*f.events, event.Short)
}

func (f *filters) removeEvent(event event) {
	newEvents := []string{}
	for _, s := range *f.events {
		if s != event.Short {
			newEvents = append(newEvents, s)
		}
	}
	f.events = &newEvents
}

func (f *filters) toggleEvent(event event) {
	if slices.Contains(*f.events, event.Short) {
		f.removeEvent(event)
	} else {
		f.addEvent(event)
	}
}
