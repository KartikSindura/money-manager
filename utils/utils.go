package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type Envelope map[string]any

func WriteJSON(w http.ResponseWriter, status int, data Envelope) error {
	js, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}

	js = append(js, '\n')
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}

func ReadIDParam(r *http.Request) (int64, error) {
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		return 0, errors.New("invalid id parameter")
	}
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return 0, errors.New("invalid id parameter")
	}
	return id, nil
}

func GetLimitOffset(r *http.Request) (limit, offset int) {
	limit = 10
	offset = 0
	limitParam := r.URL.Query().Get("limit")
	if limitParam != "" {
		limit64, err := strconv.ParseInt(limitParam, 10, 64)
		if err == nil {
			limit = int(limit64)
		}
	}
	offsetParam := r.URL.Query().Get("offset")
	if offsetParam != "" {
		offset64, err := strconv.ParseInt(offsetParam, 10, 64)
		if err == nil {
			offset = int(offset64)
		}
	}
	return limit, offset
}

func GetTransactionQueryParams(r *http.Request) (*time.Time, *time.Time, *int, *int, *string, error) {
	// /transactions?from=2022-01-01&to=2022-01-31&month=1&year=2022&type=expense

	// this gets an empty string if key is not found
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")
	monthStr := r.URL.Query().Get("month")
	yearStr := r.URL.Query().Get("year")
	typeStr := r.URL.Query().Get("type")

	var from, to *time.Time
	var month, year *int
	var _type *string

	if fromStr != "" {
		t, err := time.Parse(time.RFC3339, fromStr)
		if err != nil {
			return nil, nil, nil, nil, nil, err
		}
		from = &t
	}
	if toStr != "" {
		t, err := time.Parse(time.RFC3339, toStr)
		if err != nil {
			return nil, nil, nil, nil, nil, err
		}
		to = &t
	}
	if monthStr != "" {
		m, err := strconv.Atoi(monthStr)
		if err != nil {
			return nil, nil, nil, nil, nil, err
		}
		month = &m
	}
	if yearStr != "" {
		y, err := strconv.Atoi(yearStr)
		if err != nil {
			return nil, nil, nil, nil, nil, err
		}
		year = &y
	}
	if typeStr != "" {
		_type = &typeStr
	}
	return from, to, month, year, _type, nil
}
