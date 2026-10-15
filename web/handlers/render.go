package handlers

import (
	"log/slog"
	"net/http"
	"workshop/internal/domain/sighting"

	"github.com/a-h/templ"
)

// render writes page as HTML, falling back to the error page if it fails.
func render(w http.ResponseWriter, r *http.Request, logger *slog.Logger, page templ.Component) {
	ctx := r.Context()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := page.Render(ctx, w); err != nil {
		logger.ErrorContext(ctx, "handlers: rendering page: "+err.Error())
		renderError(ctx, w, logger, http.StatusInternalServerError, "Unable to render page.")
	}
}

// list fetches sightings, answering a 500 and reporting false when it can't.
func list(w http.ResponseWriter, r *http.Request, logger *slog.Logger, sightings SightingsLister, filter sighting.ListFilter) ([]sighting.Sighting, bool) {
	ctx := r.Context()

	items, err := sightings.List(ctx, filter)
	if err != nil {
		logger.ErrorContext(ctx, "handlers: listing sightings: "+err.Error())
		renderError(ctx, w, logger, http.StatusInternalServerError, "Unable to load sightings.")
		return nil, false
	}

	return items, true
}
