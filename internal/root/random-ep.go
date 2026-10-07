package root

import (
	"errors"
	"fmt"
	"i-need-star-trek/internal/data"
	"i-need-star-trek/internal/interactive"
	"i-need-star-trek/internal/shared"
	"math/rand"
	"time"
)

func RandomEp(shows []string, species []string, other []string, minDate time.Time, maxDate time.Time, minLen int, maxLen int) error {
	eps := *data.Episodes

	showFilters, err := shared.ParseAllShows(shows)
	if err != nil {
		return err
	}

	speciesFilters, err := shared.ParseAllSpecies(species)
	if err != nil {
		return err
	}

	otherFilters, err := shared.ParseAllOther(other)
	if err != nil {
		return err
	}

	eps = shared.FilterEpisodes(eps, showFilters, speciesFilters, otherFilters, &minDate, &maxDate, &minLen, &maxLen)

	if len(eps) == 0 {
		return errors.New("No episode matches filters")
	}

	idx := rand.Intn(len(eps))
	episode := eps[idx]

	fmt.Print(interactive.ViewEpisode(episode))

	return nil
}
