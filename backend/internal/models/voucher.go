package models

type Voucher struct {
	ID           int64  `json:"id"`
	CrewName     string `json:"crew_name"`
	CrewID       string `json:"crew_id"`
	FlightNumber string `json:"flight_number"`
	FlightDate   string `json:"flight_date"`
	AircraftType string `json:"aircraft_type"`
	Seat1        string `json:"seat1"`
	Seat2        string `json:"seat2"`
	Seat3        string `json:"seat3"`
	CreatedAt    string `json:"created_at"`
}

type CheckRequest struct {
	FlightNumber string `json:"flightNumber"`
	Date         string `json:"date"`
}

type GenerateRequest struct {
	Name         string `json:"name"`
	ID           string `json:"id"`
	FlightNumber string `json:"flightNumber"`
	Date         string `json:"date"`
	Aircraft     string `json:"aircraft"`
}

type CheckResponse struct {
	Exists bool `json:"exists"`
}

type GenerateResponse struct {
	Success bool     `json:"success"`
	Seats   []string `json:"seats"`
}

type ErrorResponse struct {
	Message string              `json:"message"`
	Errors  map[string][]string `json:"errors,omitempty"`
}
