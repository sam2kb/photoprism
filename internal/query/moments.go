package query

import (
	"fmt"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"github.com/photoprism/photoprism/internal/maps"
)

// Moment contains photo counts per month and year
type Moment struct {
	Label      string `json:"Label"`
	Country    string `json:"Country"`
	State      string `json:"State"`
	Year       int    `json:"Year"`
	Month      int    `json:"Month"`
	PhotoCount int    `json:"PhotoCount"`
}

var MomentLabels = map[string]string{
	"botanical-garden": "Botanical Gardens",
	"nature-reserve":   "Nature & Landscape",
	"landscape":        "Nature & Landscape",
	"nature":           "Nature & Landscape",
	"bay":              "Bays, Capes & Beaches",
	"beach":            "Bays, Capes & Beaches",
	"cape":             "Bays, Capes & Beaches",
}

// Slug returns an identifier string for a moment.
func (m Moment) Slug() string {
	return slug.Make(m.Title())
}

// Title returns an english title for the moment.
func (m Moment) Title() string {
	if m.Year == 0 && m.Month == 0 {
		if m.Label != "" {
			return MomentLabels[m.Label]
		}

		country := maps.CountryName(m.Country)

		if strings.Contains(m.State, country) {
			return m.State
		}

		if m.State == "" {
			return m.Country
		}

		return fmt.Sprintf("%s / %s", m.State, country)
	}

	if m.Country != "" && m.Year > 1900 && m.Month == 0 {
		if m.State != "" {
			return fmt.Sprintf("%s / %s / %d", m.State, maps.CountryName(m.Country), m.Year)
		}

		return fmt.Sprintf("%s %d", maps.CountryName(m.Country), m.Year)
	}

	if m.Year > 1900 && m.Month > 0 && m.Month <= 12 {
		date := time.Date(m.Year, time.Month(m.Month), 1, 0, 0, 0, 0, time.UTC)

		if m.Country == "" {
			return date.Format("January 2006")
		}

		return fmt.Sprintf("%s / %s", maps.CountryName(m.Country), date.Format("January 2006"))
	}

	if m.Month > 0 && m.Month <= 12 {
		return time.Month(m.Month).String()
	}

	return maps.CountryName(m.Country)
}

// A list of moments.
type Moments []Moment

// MomentsTime counts photos by month and year.
func MomentsTime(threshold int) (results Moments, err error) {
	db := UnscopedDb().Table("photos").
		Select("photos.photo_year AS year, photos.photo_month AS month, COUNT(*) AS photo_count").
		Where("photos.photo_quality >= 3 AND deleted_at IS NULL AND photos.photo_year > 0 AND photos.photo_month > 0").
		Group("photos.photo_year, photos.photo_month").
		Order("photos.photo_year DESC, photos.photo_month DESC").
		Having("photo_count >= ?", threshold)

	if err := db.Scan(&results).Error; err != nil {
		return results, err
	}

	return results, nil
}

// MomentsCountries returns the most popular countries by year.
func MomentsCountries(threshold int) (results Moments, err error) {
	db := UnscopedDb().Table("photos").
		Select("photo_country AS country, photo_year AS year, COUNT(*) AS photo_count ").
		Where("photos.photo_quality >= 3 AND deleted_at IS NULL AND photo_country <> 'zz' AND photo_year > 0").
		Group("photo_country, photo_year").
		Having("photo_count >= ?", threshold)

	if err := db.Scan(&results).Error; err != nil {
		return results, err
	}

	return results, nil
}

// MomentsStates returns the most popular states and countries by year.
func MomentsStates(threshold int) (results Moments, err error) {
	db := UnscopedDb().Table("photos").
		Select("p.loc_country AS country, p.loc_state AS state, COUNT(*) AS photo_count").
		Joins("JOIN places p ON p.id = photos.place_id").
		Where("photos.photo_quality >= 3 AND photos.deleted_at IS NULL AND p.loc_state <> '' AND p.loc_country <> 'zz'").
		Group("p.loc_country, p.loc_state").
		Having("photo_count >= ?", threshold)

	if err := db.Scan(&results).Error; err != nil {
		return results, err
	}

	return results, nil
}

// MomentsLabels returns the most popular photo labels.
func MomentsLabels(threshold int) (results Moments, err error) {
	var cats []string

	for cat, _ := range MomentLabels {
		cats = append(cats, cat)
	}

	db := UnscopedDb().Table("photos").
		Select("l.label_slug AS label, COUNT(*) AS photo_count").
		Joins("JOIN photos_labels pl ON pl.photo_id = photos.id AND pl.uncertainty < 100").
		Joins("JOIN labels l ON l.id = pl.label_id").
		Where("photos.photo_quality >= 3 AND photos.deleted_at IS NULL AND l.label_slug IN (?)", cats).
		Group("l.label_slug").
		Having("photo_count >= ?", threshold)

	if err := db.Scan(&results).Error; err != nil {
		return results, err
	}

	return results, nil
}
