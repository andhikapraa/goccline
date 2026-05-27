// Package location resolves GPS coordinates for prayer-time API lookups.
//
// On macOS we shell out to CoreLocationCLI (brew install corelocationcli),
// caching the result for 1 hour so we don't hit the system framework on
// every statusline render. On other platforms or when CoreLocationCLI is
// missing, callers should fall back to configured coordinates.
package location

import (
	"errors"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/andhikapraa/goccline/internal/cache"
)

// Coords is a lat/lon pair.
type Coords struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// FromGPS asks the OS for current coordinates. Returns (zero, err) when
// CoreLocationCLI isn't installed, location services are off, or the call
// times out. The result is cached for 1h.
func FromGPS() (Coords, error) {
	const cacheKey = "gps-location"
	var cached Coords
	if cache.Get(cacheKey, &cached) && cached.Latitude != 0 {
		return cached, nil
	}

	bin, err := exec.LookPath("CoreLocationCLI")
	if err != nil {
		return Coords{}, errors.New("CoreLocationCLI not installed (brew install corelocationcli)")
	}
	cmd := exec.Command(bin, "--format", "%latitude %longitude")
	out, err := cmd.Output()
	if err != nil {
		return Coords{}, err
	}
	fields := strings.Fields(strings.TrimSpace(string(out)))
	if len(fields) < 2 {
		return Coords{}, errors.New("CoreLocationCLI returned unexpected output: " + string(out))
	}
	lat, err1 := strconv.ParseFloat(fields[0], 64)
	lon, err2 := strconv.ParseFloat(fields[1], 64)
	if err1 != nil || err2 != nil {
		return Coords{}, errors.New("CoreLocationCLI returned non-numeric output")
	}
	c := Coords{Latitude: lat, Longitude: lon}
	_ = cache.Set(cacheKey, time.Hour, c)
	return c, nil
}
