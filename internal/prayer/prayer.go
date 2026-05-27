// Package prayer fetches and caches Islamic prayer times via the AlAdhan API.
package prayer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/andhikapraa/goccline/internal/cache"
)

// Times is the canonical set of daily prayer times.
type Times struct {
	Fajr    string `json:"Fajr"`
	Dhuhr   string `json:"Dhuhr"`
	Asr     string `json:"Asr"`
	Maghrib string `json:"Maghrib"`
	Isha    string `json:"Isha"`
}

type aladhanResponse struct {
	Code int `json:"code"`
	Data struct {
		Timings Times `json:"timings"`
	} `json:"data"`
}

// Fetch returns today's prayer times for the given coordinates. The result
// is cached on disk for 24h so consecutive statusline renders don't hit the
// network. Method 2 = Islamic Society of North America (sensible default).
// School 1 = Hanafi (matches the bash statusline's default for Indonesia).
func Fetch(latitude, longitude float64, method, school int) (Times, error) {
	key := fmt.Sprintf("prayer-%d-%.4f-%.4f-m%d-s%d",
		time.Now().YearDay(), latitude, longitude, method, school)

	var t Times
	if cache.Get(key, &t) {
		return t, nil
	}

	url := fmt.Sprintf("https://api.aladhan.com/v1/timings/%s?latitude=%f&longitude=%f&method=%d&school=%d",
		time.Now().Format("02-01-2006"), latitude, longitude, method, school)
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return Times{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Times{}, err
	}
	var ar aladhanResponse
	if err := json.Unmarshal(body, &ar); err != nil {
		return Times{}, err
	}
	if ar.Code != 200 {
		return Times{}, fmt.Errorf("aladhan: status %d", ar.Code)
	}
	// AlAdhan returns "HH:MM (TZ)" — trim parenthetical for display.
	t = trimTimes(ar.Data.Timings)
	_ = cache.Set(key, 24*time.Hour, t)
	return t, nil
}

func trimTimes(t Times) Times {
	t.Fajr = trimParen(t.Fajr)
	t.Dhuhr = trimParen(t.Dhuhr)
	t.Asr = trimParen(t.Asr)
	t.Maghrib = trimParen(t.Maghrib)
	t.Isha = trimParen(t.Isha)
	return t
}

func trimParen(s string) string {
	for i, r := range s {
		if r == ' ' {
			return s[:i]
		}
	}
	return s
}
