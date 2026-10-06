package data

import (
	"encoding/csv"
	"log"
	"os"
)

type EventType struct {
	FullName  string `json:"fullName"`
	Shortened string `json:"shortened"`
}

func importEventTypes(writeJSON bool) *map[string]string {
	f, err := os.Open(eventTypeCSV)
	if err != nil {
		log.Fatal("Unable to read input file "+eventTypeCSV, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+eventTypeCSV, err)
	}

	eventTypes := make(map[string]string)

	var longIdx int
	var shortIdx int
	for i, s := range records[0] {
		if s == "Event Type" {
			longIdx = i
		} else if s == "Shortened" {
			shortIdx = i
		}
	}

	for idx, r := range records {
		if idx != 0 {
			eventTypes[cleanString(r[longIdx])] = cleanString(r[shortIdx])
		}
	}

	if writeJSON {
		eventList := []EventType{}
		for k, v := range eventTypes {
			eventList = append(eventList, EventType{k, v})
		}
		writeJSONFile(eventTypeJSON, eventList)
	}

	return &eventTypes
}

var EventMap = importEventTypes(false)
