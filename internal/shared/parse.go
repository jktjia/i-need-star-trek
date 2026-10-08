package shared

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"strings"
	"time"
	"unicode"
)

func cleanString(s string) string {
	cleaned := strings.Map(func(r rune) rune {
		if unicode.IsGraphic(r) {
			return r
		}
		return -1
	}, s)
	return strings.TrimSpace(cleaned)
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

func ParseAllShows(shows []string) ([]string, error) {
	return parseAllStr(shows, parseShow)
}

func parseShow(show string) (string, error) {
	showMap := ShowMap

	for k, v := range showMap {
		if strings.EqualFold(show, v.Abbr) || strings.EqualFold(show, v.FullName) || strings.EqualFold("Star Trek: "+show, v.FullName) {
			return k, nil
		}
	}

	return "", errors.New("Unknown show: " + show)
}

func ParseAllSpecies(species []string) ([]string, error) {
	return parseAllStr(species, parseSpecies)
}

func parseSpecies(species string) (string, error) {
	knownSpecies := Species

	for _, sp := range knownSpecies {
		if strings.EqualFold(species, sp) || strings.EqualFold(species, sp+"s") {
			return sp, nil
		}
	}
	return "", errors.New("Unknown species: " + species)
}

func ParseAllOther(other []string) ([]string, error) {
	return parseAllStr(other, parseOther)
}

func parseOther(other string) (string, error) {
	eventMap := EventMap

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

func writeJSONFile[T any](path string, data []T) {
	file, err := os.Create(path)

	if err != nil {
		log.Fatal("Unable to write output file for "+path, err)
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		log.Fatal("Unable to marshal data for "+path, err)
	}

	file.Write(bytes)
	log.Print("Successfully created file " + path)
}
