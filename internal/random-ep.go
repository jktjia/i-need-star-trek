package internal

import (
	"errors"
	"fmt"
	"i-need-star-trek/internal/data"
	"i-need-star-trek/internal/tui"
	"math/rand"
	"time"
)

func RandomEp(shows []string, species []string, other []string, minDate time.Time, maxDate time.Time, minLen int, maxLen int) error {
	eps := *data.Episodes

	showFilters, err := parseAllShows(shows)
	if err != nil {
		return err
	}

	speciesFilters, err := parseAllSpecies(species)
	if err != nil {
		return err
	}

	otherFilters, err := parseAllOther(other)
	if err != nil {
		return err
	}

	eps = filterEpisodes(eps, showFilters, speciesFilters, otherFilters, &minDate, &maxDate, &minLen, &maxLen)

	if len(eps) == 0 {
		return errors.New("No episode matches filters")
	}

	idx := rand.Intn(len(eps))
	episode := eps[idx]

	fmt.Print(tui.ViewEpisode(episode))

	return nil
}
