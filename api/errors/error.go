package errors

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

type NotFoundError struct {
	Message string
}

func (e NotFoundError) Error() string {
	return e.Message
}

type UnauthorizedError struct {
	Message string
}

func (e UnauthorizedError) Error() string {
	return e.Message
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func SendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

func HandleError(w http.ResponseWriter, err error) {
	var pgErr *pgconn.PgError
	switch {
	case errors.As(err, &ValidationError{}):
		SendErrorResponse(w, err.Error(), http.StatusBadRequest)
	case errors.As(err, &NotFoundError{}):
		SendErrorResponse(w, err.Error(), http.StatusNotFound)
	case errors.As(err, &UnauthorizedError{}):
		time.Sleep(time.Second)
		SendErrorResponse(w, err.Error(), http.StatusUnauthorized)
	case errors.Is(err, pgx.ErrNoRows):
		SendErrorResponse(w, "Resource not found", http.StatusNotFound)
	case errors.As(err, &pgErr):
		pgErr := err.(*pgconn.PgError)
		if pgErr.Code == "23505" {
			SendErrorResponse(w, "Resource already exists", http.StatusConflict)
		} else {
			log.Printf("(%v)Database error: %v", pgErr.Code, err)
			SendErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
	default:
		debug.PrintStack()
		log.Printf("Unexpected error: %v", err)
		SendErrorResponse(w, "Internal server error", http.StatusInternalServerError)
	}
}
