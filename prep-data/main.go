package main

func main() {
	species := importSpecies(true)
	eventTypes := importEventTypes(true)

	importEpisodes(species, eventTypes, true)
	importShows(true)
}
