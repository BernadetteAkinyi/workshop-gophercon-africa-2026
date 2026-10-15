package components

import (
	"cmp"
	"fmt"
	"slices"
	"strings"
	"time"

	"workshop/internal/domain/sighting"
	"workshop/internal/domain/sighting/country"
)

// BarLink is one ranked row that links to its own page.
type BarLink struct {
	Name  string
	Count int
	Href  string
}

// Column is one bar of a column chart, in the order it's drawn.
type Column struct {
	Label string
	Count int
}

// Record is a stand-out sighting, such as the hottest day one was seen on.
type Record struct {
	Label    string
	Value    string
	Sighting sighting.Sighting
}

// Fact is a single headline figure and what it counts.
type Fact struct {
	Value string
	Label string
}

// TopAnimals ranks the animals seen most, each linking to its page.
func TopAnimals(sightings []sighting.Sighting, n int) []BarLink {
	var rows []BarLink
	for _, t := range AnimalTiles(sightings) {
		rows = append(rows, BarLink{Name: t.Name, Count: t.Count, Href: AnimalURL(t.Animal)})
	}
	if n > 0 && len(rows) > n {
		rows = rows[:n]
	}
	return rows
}

// TopCountries ranks the countries with the most sightings, each linking to its page.
func TopCountries(sightings []sighting.Sighting, n int) []BarLink {
	names := make(map[string]string)
	counts := sighting.CountBy(sightings, func(s sighting.Sighting) string {
		code := s.Location.CountryCode.String()
		if code != "" && names[code] == "" {
			names[code] = cmp.Or(s.Location.Country, code)
		}
		return code
	})
	var rows []BarLink
	for _, t := range sighting.Top(counts, n) {
		rows = append(rows, BarLink{Name: names[t.Name], Count: t.Count, Href: CountryURL(country.Code(t.Name))})
	}
	return rows
}

// Tallies ranks sightings by any key, such as who recorded them.
func Tallies(sightings []sighting.Sighting, n int, key func(sighting.Sighting) string) []sighting.Tally {
	return sighting.Top(sighting.CountBy(sightings, key), n)
}

// Sky names the weather a sighting was seen in, or nothing when none was recorded.
func Sky(s sighting.Sighting) string {
	if s.Weather.ObservedAt.IsZero() {
		return ""
	}
	name := sky(s.Weather.Condition.Code)
	return strings.ToUpper(name[:1]) + name[1:]
}

// Months counts sightings in each month of the year, January first.
func Months(sightings []sighting.Sighting) []Column {
	cols := make([]Column, 12)
	for i := range cols {
		cols[i].Label = time.Month(i + 1).String()[:3]
	}
	for _, s := range sightings {
		cols[s.HappenedAt.Month()-1].Count++
	}
	return cols
}

// Temperatures counts sightings in bands of ten degrees.
func Temperatures(sightings []sighting.Sighting) []Column {
	cols := []Column{{Label: "Under 10 °C"}, {Label: "10–20 °C"}, {Label: "20–30 °C"}, {Label: "30 °C +"}}
	for _, s := range sightings {
		if s.Weather.ObservedAt.IsZero() {
			continue
		}
		switch t := s.Weather.Temperature.Actual; {
		case t < 10:
			cols[0].Count++
		case t < 20:
			cols[1].Count++
		case t < 30:
			cols[2].Count++
		default:
			cols[3].Count++
		}
	}
	return cols
}

// Decades counts sightings in each decade, from the first to the last.
func Decades(sightings []sighting.Sighting, sum sighting.Summary) []Column {
	if sum.Sightings == 0 {
		return nil
	}
	first, last := sum.Earliest.Year()/10*10, sum.Latest.Year()/10*10
	cols := make([]Column, 0, (last-first)/10+1)
	for d := first; d <= last; d += 10 {
		cols = append(cols, Column{Label: fmt.Sprintf("%ds", d)})
	}
	for _, s := range sightings {
		cols[(s.HappenedAt.Year()/10*10-first)/10].Count++
	}
	return cols
}

// ColumnHeight sets --bar to a column's share of the tallest one.
func ColumnHeight(c Column, cols []Column) map[string]string {
	high := 0
	for _, col := range cols {
		high = max(high, col.Count)
	}
	return bar(c.Count, high)
}

// BarWidth sets --bar to a row's share of the first, longest row.
func BarWidth(count int, rows []BarLink) map[string]string {
	if len(rows) == 0 {
		return bar(0, 0)
	}
	return bar(count, rows[0].Count)
}

// Facts works out the figures the statistics page leads with.
func Facts(sightings []sighting.Sighting, sum sighting.Summary) []Fact {
	withPhoto, withWeather, groups, grouped := 0, 0, 0, 0
	var temps float64
	for _, s := range sightings {
		if len(Photos(s)) > 0 {
			withPhoto++
		}
		if !s.Weather.ObservedAt.IsZero() {
			withWeather++
			temps += s.Weather.Temperature.Actual
		}
		if s.IndividualCount != nil && *s.IndividualCount > 0 {
			groups += *s.IndividualCount
			grouped++
		}
	}

	facts := []Fact{
		{fmt.Sprint(sum.Sightings), "sightings"},
		{fmt.Sprint(len(AnimalTiles(sightings))), "kinds of animal"},
		{fmt.Sprint(sum.Countries), "countries"},
		{fmt.Sprint(sum.Recorders), "people recording"},
		{fmt.Sprint(sum.Media), "photos and media"},
		{percent(withPhoto, sum.Sightings), "of sightings have a photo"},
	}
	if withWeather > 0 {
		facts = append(facts, Fact{fmt.Sprintf("%.0f °C", temps/float64(withWeather)), "average temperature when seen"})
	}
	if grouped > 0 {
		facts = append(facts, Fact{fmt.Sprintf("%.1f", float64(groups)/float64(grouped)), "animals per sighting, where counted"})
	}
	if sum.Sightings > 0 {
		facts = append(facts, Fact{fmt.Sprint(sum.Latest.Year() - sum.Earliest.Year() + 1), "years of records"})
	}
	return facts
}

// percent renders part of a whole to the nearest whole percent.
func percent(part, whole int) string {
	if whole == 0 {
		return "0%"
	}
	return fmt.Sprintf("%d%%", part*100/whole)
}

// Records picks out the stand-out sightings: hottest, coldest, highest, biggest group, windiest.
func Records(sightings []sighting.Sighting) []Record {
	var weathered, high, counted []sighting.Sighting
	for _, s := range sightings {
		if !s.Weather.ObservedAt.IsZero() {
			weathered = append(weathered, s)
		}
		if s.Location.ElevationInMetres != nil {
			high = append(high, s)
		}
		if s.IndividualCount != nil {
			counted = append(counted, s)
		}
	}

	var records []Record
	if len(weathered) > 0 {
		hot := slices.MaxFunc(weathered, func(a, b sighting.Sighting) int {
			return cmp.Compare(a.Weather.Temperature.Actual, b.Weather.Temperature.Actual)
		})
		cold := slices.MinFunc(weathered, func(a, b sighting.Sighting) int {
			return cmp.Compare(a.Weather.Temperature.Actual, b.Weather.Temperature.Actual)
		})
		windy := slices.MaxFunc(weathered, func(a, b sighting.Sighting) int {
			return cmp.Compare(a.Weather.Wind.Speed, b.Weather.Wind.Speed)
		})
		wet := slices.MaxFunc(weathered, func(a, b sighting.Sighting) int {
			return cmp.Compare(a.Weather.Precipitation.Total, b.Weather.Precipitation.Total)
		})
		records = append(records,
			Record{"Hottest day", degrees(hot.Weather.Temperature.Actual), hot},
			Record{"Coldest day", degrees(cold.Weather.Temperature.Actual), cold},
			Record{"Windiest day", kmh(windy.Weather.Wind.Speed), windy},
		)
		if wet.Weather.Precipitation.Total > 0 {
			records = append(records, Record{"Wettest day", fmt.Sprintf("%.1f mm of rain", wet.Weather.Precipitation.Total), wet})
		}
	}
	if len(high) > 0 {
		top := slices.MaxFunc(high, func(a, b sighting.Sighting) int {
			return cmp.Compare(*a.Location.ElevationInMetres, *b.Location.ElevationInMetres)
		})
		records = append(records, Record{"Highest up", Metres(*top.Location.ElevationInMetres) + " above sea level", top})
	}
	if len(counted) > 0 {
		big := slices.MaxFunc(counted, func(a, b sighting.Sighting) int {
			return cmp.Compare(*a.IndividualCount, *b.IndividualCount)
		})
		records = append(records, Record{"Biggest group", individuals(*big.IndividualCount), big})
	}
	if len(sightings) > 0 {
		oldest := slices.MinFunc(sightings, func(a, b sighting.Sighting) int { return a.HappenedAt.Compare(b.HappenedAt) })
		records = append(records, Record{"Oldest record", oldest.HappenedAt.Format("2 January 2006"), oldest})
	}
	return records
}
