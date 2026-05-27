package location

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/andhikapraa/goccline/internal/cache"
)

// reverseResponse is the slice of Nominatim's payload we care about.
type reverseResponse struct {
	DisplayName string `json:"display_name"`
	Address     struct {
		Suburb       string `json:"suburb"`
		Village      string `json:"village"`
		Town         string `json:"town"`
		City         string `json:"city"`
		County       string `json:"county"`
		State        string `json:"state"`
		CountryCode  string `json:"country_code"`
		Country      string `json:"country"`
	} `json:"address"`
}

// Label returns a human-readable place name for the coordinates by reverse
// geocoding via Nominatim (OpenStreetMap). Cached for 24h per (lat,lon)
// rounded to 2 decimals so small drift doesn't bust the cache.
//
// Returns "" on any failure — caller should fall back to the configured
// label.
func Label(c Coords) string {
	key := fmt.Sprintf("geo-%.2f-%.2f", c.Latitude, c.Longitude)
	var cached string
	if cache.Get(key, &cached) && cached != "" {
		return cached
	}

	url := fmt.Sprintf(
		"https://nominatim.openstreetmap.org/reverse?format=json&lat=%f&lon=%f&zoom=10",
		c.Latitude, c.Longitude)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ""
	}
	// Nominatim's usage policy requires identifying yourself.
	req.Header.Set("User-Agent", "goccline/0.7 (https://github.com/andhikapraa/goccline)")

	client := http.Client{Timeout: 3 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return ""
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	var r reverseResponse
	if err := json.Unmarshal(body, &r); err != nil {
		return ""
	}

	place := firstNonEmpty(r.Address.City, r.Address.Town, r.Address.Village, r.Address.Suburb, r.Address.County, r.Address.State)
	region := firstNonEmpty(r.Address.State, r.Address.Country)
	label := ""
	switch {
	case place != "" && region != "" && place != region:
		label = place + ", " + region
	case place != "":
		label = place
	case region != "":
		label = region
	default:
		label = r.DisplayName
	}
	if label != "" {
		_ = cache.Set(key, 24*time.Hour, label)
	}
	return label
}

func firstNonEmpty(ss ...string) string {
	for _, s := range ss {
		if s != "" {
			return s
		}
	}
	return ""
}
