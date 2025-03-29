package utils

import (
	"encoding/json"
	"errors"
	"frappuccino/internal/models"
	"net/http"
	"strconv"
	"math"
)

// sending responses in the json format
//
//	{
//		"error": "Internal Server Error"
//	}
type Response map[string]interface{}

func SendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func MapErrorToResponse(err error, validationMap any) (int, any) {
	switch {
	// General errors
	case errors.Is(err, models.ErrInvalidID):
		return http.StatusBadRequest, Response{"error": err.Error()}
	case errors.Is(err, models.ErrNoRecord):
		return http.StatusNotFound, Response{"error": err.Error()}
	case errors.Is(err, models.ErrMissingFields):
		return http.StatusBadRequest, validationMap

	// Inventory errors
	case errors.Is(err, models.ErrDuplicateInventory),
		errors.Is(err, models.ErrNegativeQuantity),
		errors.Is(err, models.ErrInvalidEnumTypeInventory):
		return http.StatusBadRequest, Response{"error": err.Error()}

	// Menu errors
	case errors.Is(err, models.ErrDuplicateMenuItem),
		errors.Is(err, models.ErrNegativePrice),
		errors.Is(err, models.ErrForeignKeyConstraintMenuInventory):
		return http.StatusBadRequest, Response{"error": err.Error()}

	// Order errors
	case errors.Is(err, models.ErrDuplicateOrder),
		errors.Is(err, models.ErrInvalidFilterOption),
		errors.Is(err, models.ErrForeignKeyConstraintOrderMenu):
		return http.StatusBadRequest, Response{"error": err.Error()}
	
	// Report errors
	case errors.Is(err, models.ErrInvalidPrice):
		return http.StatusBadRequest, Response{"error": err.Error()}

	// Default catch-all
	default:
		return http.StatusInternalServerError, Response{"error": "Internal Server Error"}
	}
}

func ValidatePrices(minPriceStr string, maxPriceStr string) (float64, float64, error) {
	var minPrice, maxPrice float64
	var err error

	if minPriceStr == "absent" {
		minPrice = -1
	} else {
		minPrice, err = strconv.ParseFloat(minPriceStr, 64)
		if err != nil || math.IsNaN(minPrice) || math.IsInf(minPrice, 0) {
			return 0, 0, models.ErrInvalidPrice
		}
	}
	if maxPriceStr == "absent" {
		maxPrice = -1
	} else {
		maxPrice, err = strconv.ParseFloat(maxPriceStr, 64)
		if err != nil || math.IsNaN(maxPrice) || math.IsInf(maxPrice, 0) {
			return 0, 0, models.ErrInvalidPrice
		}
	}
	if minPrice != -1 && maxPrice != -1 && minPrice > maxPrice {
		return 0, 0, models.ErrInvalidPrice
	}
	return minPrice, maxPrice, nil
}
