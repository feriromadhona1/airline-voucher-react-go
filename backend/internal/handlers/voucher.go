package handlers

import (
	"net/http"
	"strings"
	"time"

	"airline-voucher/internal/database"
	"airline-voucher/internal/models"
	"airline-voucher/internal/seats"

	"github.com/labstack/echo/v4"
)

type VoucherHandler struct {
	store *database.Store
}

func NewVoucherHandler(store *database.Store) *VoucherHandler {
	return &VoucherHandler{store: store}
}

func (h *VoucherHandler) Check(c echo.Context) error {
	var req models.CheckRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Invalid JSON request body.",
		})
	}

	if fieldErrors := validateCheck(req); len(fieldErrors) > 0 {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Validation failed.",
			Errors:  fieldErrors,
		})
	}

	exists, err := h.store.ExistsForFlight(strings.TrimSpace(req.FlightNumber), req.Date)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Failed to check voucher assignment.",
		})
	}

	return c.JSON(http.StatusOK, models.CheckResponse{Exists: exists})
}

func (h *VoucherHandler) Generate(c echo.Context) error {
	var req models.GenerateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Invalid JSON request body.",
		})
	}

	if fieldErrors := validateGenerate(req); len(fieldErrors) > 0 {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Validation failed.",
			Errors:  fieldErrors,
		})
	}

	name := strings.TrimSpace(req.Name)
	crewID := strings.TrimSpace(req.ID)
	flightNumber := strings.TrimSpace(req.FlightNumber)

	exists, err := h.store.ExistsForFlight(flightNumber, req.Date)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Failed to check existing vouchers.",
		})
	}
	if exists {
		return c.JSON(http.StatusConflict, models.ErrorResponse{
			Message: "Vouchers have already been generated for this flight and date.",
		})
	}

	selected, err := seats.Generate(req.Aircraft, 3)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: err.Error(),
		})
	}

	voucher := &models.Voucher{
		CrewName:     name,
		CrewID:       crewID,
		FlightNumber: flightNumber,
		FlightDate:   req.Date,
		AircraftType: req.Aircraft,
		Seat1:        selected[0],
		Seat2:        selected[1],
		Seat3:        selected[2],
		CreatedAt:    time.Now().UTC().Format(time.RFC3339),
	}

	if err := h.store.Create(voucher); err != nil {
		if database.IsUniqueViolation(err) {
			return c.JSON(http.StatusConflict, models.ErrorResponse{
				Message: "Vouchers have already been generated for this flight and date.",
			})
		}
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Failed to save voucher assignment.",
		})
	}

	return c.JSON(http.StatusOK, models.GenerateResponse{
		Success: true,
		Seats:   selected,
	})
}

func validateCheck(req models.CheckRequest) map[string][]string {
	errors := map[string][]string{}

	if strings.TrimSpace(req.FlightNumber) == "" {
		errors["flightNumber"] = []string{"The flight number is required."}
	}
	if req.Date == "" {
		errors["date"] = []string{"The flight date is required."}
	} else if !isValidDate(req.Date) {
		errors["date"] = []string{"The flight date must be in YYYY-MM-DD format."}
	}

	return errors
}

func validateGenerate(req models.GenerateRequest) map[string][]string {
	errors := map[string][]string{}

	if strings.TrimSpace(req.Name) == "" {
		errors["name"] = []string{"The crew name is required."}
	}
	if strings.TrimSpace(req.ID) == "" {
		errors["id"] = []string{"The crew ID is required."}
	}
	if strings.TrimSpace(req.FlightNumber) == "" {
		errors["flightNumber"] = []string{"The flight number is required."}
	}
	if req.Date == "" {
		errors["date"] = []string{"The flight date is required."}
	} else if !isValidDate(req.Date) {
		errors["date"] = []string{"The flight date must be in YYYY-MM-DD format."}
	}
	if strings.TrimSpace(req.Aircraft) == "" {
		errors["aircraft"] = []string{"The aircraft type is required."}
	} else if !seats.IsValidAircraft(req.Aircraft) {
		errors["aircraft"] = []string{"The selected aircraft type is invalid. Allowed: ATR, Airbus 320, Boeing 737 Max."}
	}

	return errors
}

func isValidDate(value string) bool {
	_, err := time.Parse("2006-01-02", value)
	return err == nil
}
