package seats

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

type layout struct {
	rows  int
	seats []string
}

var layouts = map[string]layout{
	"ATR": {
		rows:  18,
		seats: []string{"A", "C", "D", "F"},
	},
	"Airbus 320": {
		rows:  32,
		seats: []string{"A", "B", "C", "D", "E", "F"},
	},
	"Boeing 737 Max": {
		rows:  32,
		seats: []string{"A", "B", "C", "D", "E", "F"},
	},
}

func SupportedAircraftTypes() []string {
	types := make([]string, 0, len(layouts))
	for name := range layouts {
		types = append(types, name)
	}
	sort.Strings(types)
	return types
}

func IsValidAircraft(aircraft string) bool {
	_, ok := layouts[aircraft]
	return ok
}

func AvailableSeats(aircraftType string) ([]string, error) {
	layout, ok := layouts[aircraftType]
	if !ok {
		return nil, fmt.Errorf("unsupported aircraft type [%s]", aircraftType)
	}

	seats := make([]string, 0, layout.rows*len(layout.seats))
	for row := 1; row <= layout.rows; row++ {
		for _, letter := range layout.seats {
			seats = append(seats, fmt.Sprintf("%d%s", row, letter))
		}
	}
	return seats, nil
}

func Generate(aircraftType string, count int) ([]string, error) {
	available, err := AvailableSeats(aircraftType)
	if err != nil {
		return nil, err
	}
	if count > len(available) {
		return nil, fmt.Errorf("cannot generate %d unique seats for aircraft type [%s]", count, aircraftType)
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	rng.Shuffle(len(available), func(i, j int) {
		available[i], available[j] = available[j], available[i]
	})

	selected := append([]string(nil), available[:count]...)
	sort.Slice(selected, func(i, j int) bool {
		return naturalLess(selected[i], selected[j])
	})
	return selected, nil
}

func naturalLess(a, b string) bool {
	ra, la := splitSeat(a)
	rb, lb := splitSeat(b)
	if ra != rb {
		return ra < rb
	}
	return la < lb
}

func splitSeat(seat string) (int, string) {
	i := 0
	for i < len(seat) && seat[i] >= '0' && seat[i] <= '9' {
		i++
	}
	row := 0
	for _, c := range seat[:i] {
		row = row*10 + int(c-'0')
	}
	return row, seat[i:]
}
