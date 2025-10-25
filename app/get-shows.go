package main

import (
	"encoding/json"
	"os"
	"slices"
	"time"
)

const showJSON = "../data/shows.json"
const speciesJSON = "../data/species.json"
const eventJSON = "../data/event-types.json"

type show struct {
	FullName  string     `json:"full"`
	Abbr      string     `json:"abbr"`
	Seasons   int        `json:"seasons"`
	StartDate *time.Time `json:"start"`
	EndDate   *time.Time `json:"end"`
}

type event struct {
	FullName string `json:"fullName"`
	Short    string `json:"shortened"`
}

func getShows() ([]show, map[string]show) {
	data, err := os.ReadFile(showJSON)
	if err != nil {
		panic(err)
	}

	var shows []show
	json.Unmarshal(data, &shows)

	slices.SortFunc(shows, func(s1, s2 show) int {
		return s1.StartDate.Compare(*s2.StartDate)
	})

	showMap := make(map[string]show)

	for _, s := range shows {
		showMap[s.Abbr] = s
	}

	return shows, showMap
}

func getSpecies() []string {
	data, err := os.ReadFile(speciesJSON)
	if err != nil {
		panic(err)
	}

	var species []string
	json.Unmarshal(data, &species)
	return species
}

func getEvents() []event {
	data, err := os.ReadFile(eventJSON)
	if err != nil {
		panic(err)
	}

	var events []event
	json.Unmarshal(data, &events)
	return events
}
