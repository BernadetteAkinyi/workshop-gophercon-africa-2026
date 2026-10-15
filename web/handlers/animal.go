package handlers

import (
	"log/slog"
	"net/http"
	"workshop/internal/domain/sighting"
	"workshop/internal/domain/species"
	"workshop/web/views/pages"
)

func Animal(logger *slog.Logger, sightings SightingsLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		animal, err := species.ParseAnimal(r.PathValue("animal"))
		if err != nil {
			http.NotFound(w, r)
			return
		}

		items, ok := list(w, r, logger, sightings, sighting.ListFilter{Animals: []species.Animal{animal}})
		if !ok {
			return
		}

		render(w, r, logger, pages.Animal(animal, items))
	}
}
