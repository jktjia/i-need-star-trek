package data

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"time"
)

type Show struct {
	FullName  string     `json:"full"`
	Abbr      string     `json:"abbr"`
	Seasons   int        `json:"seasons"`
	StartDate *time.Time `json:"start"`
	EndDate   *time.Time `json:"end"`
}

func importShows(writeJSON bool) *map[string]Show {
	f, err := os.Open(showCSV)
	if err != nil {
		log.Fatal("Unable to read input file "+showCSV, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+showCSV, err)
	}

	showMap := make(map[string]Show)

	var longIdx int
	var shortIdx int
	var seasonIdx int
	var startIdx int
	var endIdx int
	for i, s := range records[0] {
		switch s {
		case "Full Name":
			longIdx = i
		case "Abbreviation":
			shortIdx = i
		case "Seasons":
			seasonIdx = i
		case "Start Date":
			startIdx = i
		case "End Date":
			endIdx = i
		}
	}

	showList := []Show{}
	for idx, r := range records {
		if idx != 0 {
			show := Show{}
			show.Abbr = cleanString(r[shortIdx])
			n, err := strconv.Atoi(r[seasonIdx])
			if err != nil {
				log.Fatal(err.Error())
			} else {
				show.Seasons = n
			}
			start, err := time.Parse(time.DateOnly, r[startIdx])
			if err != nil {
				log.Fatal(err.Error())
			} else {
				show.StartDate = &start
			}
			end, err := time.Parse(time.DateOnly, r[endIdx])
			if err == nil {
				show.EndDate = &end
			}
			show.FullName = cleanString(r[longIdx])
			showMap[show.Abbr] = show
			showList = append(showList, show)
		}
	}

	if writeJSON {
		writeJSONFile(showJSON, showList)
	}

	return &showMap
}

var (
	ShowMap = importShows(false)
)
