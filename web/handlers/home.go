package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"workshop/internal/domain/sighting"
	"workshop/internal/domain/species"
	"workshop/web/views/pages"
)

type SightingsLister interface {
	List(ctx context.Context, filter sighting.ListFilter) ([]sighting.Sighting, error)
}

// africanWildlife is every animal the site shows, on every page that lists sightings.
var africanWildlife = []species.Animal{
	species.AnimalLion,
	species.AnimalLeopard,
	species.AnimalElephant,
	species.AnimalRhinoceros,
	species.AnimalAfricanBuffalo,
	species.AnimalCheetah,
	species.AnimalGiraffe,
	species.AnimalZebra,
	species.AnimalHippopotamus,
	species.AnimalWildebeest,
	species.AnimalImpala,
	species.AnimalCommonWarthog,
	species.AnimalSpottedHyena,
	species.AnimalAfricanWildDog,
	species.AnimalNileCrocodile,
	species.AnimalGorilla,
	species.AnimalChimpanzee,
	species.AnimalRingTailedLemur,
	species.AnimalAfricanPenguin,
	species.AnimalGreaterFlamingo,
	species.AnimalShoebill,
}

// featuredWildlife narrows the home page to the Big Five and the classic
// safari animals; the animal, country and statistics pages keep the full set.
var featuredWildlife = []species.Animal{
	species.AnimalLion,
	species.AnimalLeopard,
	species.AnimalElephant,
	species.AnimalRhinoceros,
	species.AnimalAfricanBuffalo,
	species.AnimalGiraffe,
	species.AnimalZebra,
	species.AnimalHippopotamus,
	species.AnimalCheetah,
}

func Home(logger *slog.Logger, sightings SightingsLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		items, err := sightings.List(ctx, sighting.ListFilter{
			Animals: featuredWildlife,
		})
		if err != nil {
			logger.ErrorContext(ctx, "handlers: listing sightings: "+err.Error())
			renderError(ctx, w, logger, http.StatusInternalServerError, "Unable to load sightings.")
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err = pages.Home(items).Render(ctx, w); err != nil {
			logger.ErrorContext(ctx, "handlers: rendering sightings: "+err.Error())
			renderError(ctx, w, logger, http.StatusInternalServerError, "Unable to render page.")
		}
	}
}
