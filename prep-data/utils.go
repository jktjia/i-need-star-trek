package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"unicode"
)

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

func cleanString(s string) string {
	cleaned := strings.Map(func(r rune) rune {
		if unicode.IsGraphic(r) {
			return r
		}
		return -1
	}, s)
	return strings.TrimSpace(cleaned)
}
