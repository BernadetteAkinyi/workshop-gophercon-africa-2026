package handlers

import (
	"log/slog"
	"net/http"
	"workshop/internal/domain/sighting"
	"workshop/web/views/pages"
)

func Stats(logger *slog.Logger, sightings SightingsLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, ok := list(w, r, logger, sightings, sighting.ListFilter{Animals: africanWildlife})
		if !ok {
			return
		}

		render(w, r, logger, pages.Stats(items))
	}
}

func Gallery(logger *slog.Logger, sightings SightingsLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, ok := list(w, r, logger, sightings, sighting.ListFilter{Animals: africanWildlife})
		if !ok {
			return
		}

		render(w, r, logger, pages.Gallery(items))
	}
}
