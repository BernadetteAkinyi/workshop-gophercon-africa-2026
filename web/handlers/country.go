package handlers

import (
	"log/slog"
	"net/http"
	"workshop/internal/domain/sighting"
	"workshop/internal/domain/sighting/country"
	"workshop/web/views/pages"
)

func Country(logger *slog.Logger, sightings SightingsLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code, err := country.Parse(r.PathValue("code"))
		if err != nil || code == "" {
			http.NotFound(w, r)
			return
		}

		items, ok := list(w, r, logger, sightings, sighting.ListFilter{Animals: africanWildlife, Country: code})
		if !ok {
			return
		}

		render(w, r, logger, pages.Country(code, items))
	}
}
