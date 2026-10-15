package pages

import (
	"cmp"
	"slices"
	"strings"

	"workshop/internal/domain/sighting"
	"workshop/internal/domain/sighting/country"
)

type rank struct {
	Name, Value string
}

// hasPosition reports whether a sighting can go on a map.
func hasPosition(s sighting.Sighting) bool {
	return s.Location.Latitude != 0 || s.Location.Longitude != 0
}

// countryName is a country's name as the sightings give it, or its code.
func countryName(code country.Code, sightings []sighting.Sighting) string {
	for _, s := range sightings {
		if s.Location.CountryCode == code && s.Location.Country != "" {
			return s.Location.Country
		}
	}
	return code.String()
}

// scientificNames lists the distinct species names an animal's sightings were recorded under.
func scientificNames(sightings []sighting.Sighting) string {
	var names []string
	for _, s := range sightings {
		n := cmp.Or(s.Species.CanonicalName, s.Species.ScientificName)
		if n != "" && !slices.Contains(names, n) {
			names = append(names, n)
		}
	}
	slices.Sort(names)
	return strings.Join(names, " · ")
}

// taxonomy is a sighting's classification, broadest first, skipping ranks GBIF left blank.
func taxonomy(s sighting.Sighting) []rank {
	sp := s.Species
	all := []rank{
		{"Kingdom", sp.Kingdom}, {"Phylum", sp.Phylum}, {"Class", sp.Class},
		{"Order", sp.Order}, {"Family", sp.Family}, {"Genus", sp.Genus}, {"Species", sp.Species},
	}
	return slices.DeleteFunc(all, func(r rank) bool { return r.Value == "" })
}

// recorder is who recorded a sighting.
func recorder(s sighting.Sighting) string {
	return s.RecordedBy
}
