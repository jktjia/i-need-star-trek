package shared

import (
	"encoding/csv"
	"log"
	"os"
)

func importSpecies(writeJSON bool) *[]string {
	f, err := os.Open(speciesCSV)
	if err != nil {
		log.Fatal("Unable to read input file "+speciesCSV, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+speciesCSV, err)
	}

	species := []string{}

	for idx, r := range records {
		if idx > 0 {
			species = append(species, cleanString(r[0]))
		}
	}

	if writeJSON {
		writeJSONFile(speciesJSON, species)
	}

	return &species
}
