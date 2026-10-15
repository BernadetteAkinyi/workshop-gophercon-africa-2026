package components

import (
	"cmp"
	"fmt"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/a-h/templ"

	"workshop/internal/domain/sighting"
	"workshop/internal/domain/sighting/country"
	"workshop/internal/domain/species"
)

// AnimalTile is one animal on the "meet the animals" wall.
type AnimalTile struct {
	Animal species.Animal
	Name   string
	Count  int
	Photo  string
}

// AnimalOf finds which of the site's animals a species belongs to.
func AnimalOf(sp species.Species) (species.Animal, bool) {
	for _, a := range species.Animals() {
		if slices.Contains(a.TaxonKeys(), sp.GBIFKey) {
			return a, true
		}
	}
	return "", false
}

// AnimalName turns an animal such as african_wild_dog into "African wild dog".
func AnimalName(a species.Animal) string {
	name := strings.ReplaceAll(a.String(), "_", " ")
	r, size := utf8.DecodeRuneInString(name)
	return string(unicode.ToUpper(r)) + name[size:]
}

// AnimalURL is the path of an animal's page.
func AnimalURL(a species.Animal) string {
	return "/animals/" + a.String()
}

// CountryURL is the path of a country's page.
func CountryURL(c country.Code) string {
	return "/countries/" + strings.ToLower(c.String())
}

// Photos keeps the media that are pictures, with a link to show.
func Photos(s sighting.Sighting) []sighting.Media {
	var out []sighting.Media
	for _, m := range s.Media {
		if m.Type == "StillImage" && m.URL != "" {
			out = append(out, m)
		}
	}
	return out
}

// firstPhoto is the picture a card leads with, if there is one.
func firstPhoto(s sighting.Sighting) (sighting.Media, bool) {
	p := Photos(s)
	if len(p) == 0 {
		return sighting.Media{}, false
	}
	return p[0], true
}

// AnimalTiles groups sightings by animal, most seen first, each with a photo when one exists.
func AnimalTiles(sightings []sighting.Sighting) []AnimalTile {
	byAnimal := make(map[species.Animal]*AnimalTile)
	for _, s := range sightings {
		a, ok := AnimalOf(s.Species)
		if !ok {
			continue
		}
		t, seen := byAnimal[a]
		if !seen {
			t = &AnimalTile{Animal: a, Name: AnimalName(a)}
			byAnimal[a] = t
		}
		t.Count++
		if p, ok := firstPhoto(s); ok && t.Photo == "" {
			t.Photo = p.URL
		}
	}

	tiles := make([]AnimalTile, 0, len(byAnimal))
	for _, t := range byAnimal {
		tiles = append(tiles, *t)
	}
	slices.SortFunc(tiles, func(a, b AnimalTile) int {
		return cmp.Or(cmp.Compare(b.Count, a.Count), cmp.Compare(a.Name, b.Name))
	})
	return tiles
}

// GlobePoints is each sighting's position as a [lat, lng] pair, for the 3D globe.
func GlobePoints(sightings []sighting.Sighting) [][2]float64 {
	points := make([][2]float64, 0, len(sightings))
	for _, s := range sightings {
		if s.Location.Latitude == 0 && s.Location.Longitude == 0 {
			continue
		}
		points = append(points, [2]float64{round2(s.Location.Latitude), round2(s.Location.Longitude)})
	}
	return points
}

// round2 trims a coordinate to two decimal places, about a kilometre.
func round2(v float64) float64 {
	return float64(int(v*100)) / 100
}

// Basis puts GBIF's kind of record, such as HUMAN_OBSERVATION, into words.
func Basis(b string) string {
	switch b {
	case "":
		return ""
	case "HUMAN_OBSERVATION":
		return "Seen by a person"
	case "MACHINE_OBSERVATION":
		return "Camera trap or sensor"
	case "PRESERVED_SPECIMEN":
		return "Museum specimen"
	case "LIVING_SPECIMEN":
		return "Living specimen, such as in a zoo"
	case "MATERIAL_SAMPLE":
		return "Material sample"
	case "OCCURRENCE":
		return "Occurrence"
	}
	name := strings.ToLower(strings.ReplaceAll(b, "_", " "))
	r, size := utf8.DecodeRuneInString(name)
	return string(unicode.ToUpper(r)) + name[size:]
}

// Metres renders a distance to the nearest whole metre.
func Metres(v float64) string {
	return fmt.Sprintf("%.0f m", v)
}

// Distance renders a distance in metres, or in kilometres once it's a kilometre or more.
func Distance(v float64) string {
	if v >= 1000 {
		return fmt.Sprintf("%.1f km", v/1000)
	}
	return Metres(v)
}

// Coords renders a position as degrees north/south and east/west.
func Coords(c sighting.Coordinates) string {
	ns, ew := "N", "E"
	lat, lng := c.Latitude, c.Longitude
	if lat < 0 {
		ns, lat = "S", -lat
	}
	if lng < 0 {
		ew, lng = "W", -lng
	}
	return fmt.Sprintf("%.3f° %s, %.3f° %s", lat, ns, lng, ew)
}

// DisplayName is the name a reader would recognise for a sighting's species.
func DisplayName(s sighting.Sighting) string {
	return displayName(s.Species)
}

// Place joins whichever parts of the location GBIF supplied.
func Place(l sighting.Location) string {
	return place(l)
}

// Individuals renders how many animals were seen, or nothing when GBIF didn't say.
func Individuals(n *int) string {
	return count(n)
}

// PhotoItem is one photo on a wall, with the sighting it belongs to.
type PhotoItem struct {
	Photo    sighting.Media
	Sighting sighting.Sighting
}

// PhotoItems lists every photo across the sightings, keeping at most limit; 0 keeps all.
func PhotoItems(sightings []sighting.Sighting, limit int) []PhotoItem {
	var items []PhotoItem
	for _, s := range sightings {
		for _, p := range Photos(s) {
			if limit > 0 && len(items) == limit {
				return items
			}
			items = append(items, PhotoItem{Photo: p, Sighting: s})
		}
	}
	return items
}

// photoAttrs marks a photo with its sighting's filter attributes, when the wall is filtered.
func photoAttrs(p PhotoItem, filterable bool) templ.Attributes {
	if !filterable {
		return templ.Attributes{}
	}
	return FilterAttrs(p.Sighting)
}

// photoAlt describes a photo by the animal in it, unless the photo has its own title.
func photoAlt(p PhotoItem) string {
	return cmp.Or(p.Photo.Title, displayName(p.Sighting.Species)+" in "+place(p.Sighting.Location))
}

// photoCaption is the lightbox caption: the animal, where, and who took it.
func photoCaption(p PhotoItem) string {
	parts := []string{displayName(p.Sighting.Species), place(p.Sighting.Location)}
	if c := credit(p.Photo); c != "" {
		parts = append(parts, c)
	}
	return strings.Join(parts, " · ")
}

// plural renders a count with its noun, adding an s for anything but one.
func plural(n int, noun string) string {
	if n == 1 {
		return "1 " + noun
	}
	return fmt.Sprintf("%d %ss", n, noun)
}
