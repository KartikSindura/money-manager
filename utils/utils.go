package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v4"
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

func GetTransactionQueryParams(r *http.Request) (*time.Time, *time.Time, *int, *int, *string, *string, error) {
	// /transactions?from=2022-01-01&to=2022-01-31&month=1&year=2022&type=expense

	// this gets an empty string if key is not found
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")
	monthStr := r.URL.Query().Get("month")
	yearStr := r.URL.Query().Get("year")
	typeStr := r.URL.Query().Get("type")
	categoryStr := r.URL.Query().Get("category")

	var from, to *time.Time
	var month, year *int
	var _type, categoryName *string

	if fromStr != "" {
		t, err := time.Parse(time.RFC3339, fromStr)
		if err != nil {
			return nil, nil, nil, nil, nil, nil, err
		}
		from = &t
	}
	if toStr != "" {
		t, err := time.Parse(time.RFC3339, toStr)
		if err != nil {
			return nil, nil, nil, nil, nil, nil, err
		}
		to = &t
	}
	if monthStr != "" {
		m, err := strconv.Atoi(monthStr)
		if err != nil {
			return nil, nil, nil, nil, nil, nil, err
		}
		month = &m
	}
	if yearStr != "" {
		y, err := strconv.Atoi(yearStr)
		if err != nil {
			return nil, nil, nil, nil, nil, nil, err
		}
		year = &y
	}
	if typeStr != "" {
		_type = &typeStr
	}

	if categoryStr != "" {
		lower := strings.ToLower(categoryStr)
		categoryName = &lower
	}

	return from, to, month, year, _type, categoryName, nil
}

const (
	ScopeAuth = "authentication"
)

var secretKey = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

func CreateToken(user_id int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id": user_id,
			"exp":     time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
