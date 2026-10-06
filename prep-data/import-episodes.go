package main

import (
	"encoding/csv"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

type EpisodeInfo struct {
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

func importEpisodes(species *[]string, events *map[string]string, writeJSON bool) *[]EpisodeInfo {
	f, err := os.Open(episodeCSV)
	if err != nil {
		log.Fatal("Unable to read input file "+episodeCSV, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+episodeCSV, err)
	}

	episodes := []EpisodeInfo{}
	var header []string
	for idx, row := range records {
		if idx == 0 {
			header = row
		} else {
			e := parseRow(row, header, *species, *events)
			episodes = append(episodes, e)
		}
	}

	if writeJSON {
		writeJSONFile(episodeJSON, episodes)
	}

	return &episodes
}

func parseRow(row []string, header []string, species []string, events map[string]string) EpisodeInfo {
	var e EpisodeInfo
	e.Species = make(map[string]bool)
	e.EpCategory = make(map[string]bool)
	for idx, val := range row {
		col := cleanString(header[idx])
		val = cleanString(val)
		if col == "Show" {
			e.Show = val
		} else if col == "Production Number" {
			e.ProdNumber = val
		} else if col == "Stardate Start" && val != "" {
			e.StardateStart = &val
		} else if col == "Stardate End" && val != "" {
			e.StardateEnd = &val
		} else if col == "Name" {
			e.Name = strings.Trim(val, "\"")
		} else if col == "Summary" {
			e.Summary = val
		} else if (col == "Season" || col == "Episode Number" || col == "Part 2 Number" || col == "Length (minutes)") && val != "" {
			n, err := strconv.Atoi(val)
			if err != nil {
				log.Print(err.Error())
			} else if col == "Season" {
				e.Season = n
			} else if col == "Episode Number" {
				e.EpNumber = n
			} else if col == "Part 2 Number" {
				e.EpNumber2 = &n
			} else if col == "Length (minutes)" {
				e.Length = n
			}
		} else if col == "Release Date" {
			date, err := time.Parse(time.DateOnly, val)
			if err != nil {
				log.Print(err.Error())
			} else {
				e.ReleaseDate = &date
			}
		} else if (col == "Director(s)" || col == "Writer(s)" || col == "Crossover" || col == "Central Characters") && val != "" {
			l := strings.Split(val, ",")
			for idx, dir := range l {
				l[idx] = strings.TrimSpace(dir)
			}
			switch col {
			case "Director(s)":
				e.Directors = &l
			case "Writer(s)":
				e.Writers = &l
			case "Crossover":
				e.Crossover = &l
			case "Central Characters":
				e.MainChars = &l
			}
		} else if slices.Contains(species, col) {
			e.Species[col] = (val == "1")
		} else if shortName, containsKey := events[col]; containsKey {
			e.EpCategory[shortName] = (val == "1")
		} else {
			log.Print(col)
		}
	}
	return e
}
