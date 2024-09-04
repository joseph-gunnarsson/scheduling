package middleware

import (
	"context"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/joseph-gunnarsson/scheduling/api/errors"
	db "github.com/joseph-gunnarsson/scheduling/db/models"
	"github.com/joseph-gunnarsson/scheduling/internals/auth"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

type MiddlewareManager struct {
	db *pgx.Conn
}
type ContextKey string

const UserKey ContextKey = "user"

func NewMiddlewareManager(db *pgx.Conn) *MiddlewareManager {
	return &MiddlewareManager{
		db: db,
	}
}

func (m *MiddlewareManager) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			errors.HandleError(rw, errors.UnauthorizedError{Message: "Missing Authorization header"})
			return
		}
		tokenParts := strings.SplitN(authHeader, " ", 2)
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			errors.HandleError(rw, errors.UnauthorizedError{Message: "Invalid Authorization header format"})
			return
		}
		token := tokenParts[1]
		err := auth.VerifyToken(token)
		if err != nil {
			errors.HandleError(rw, errors.UnauthorizedError{Message: "Invalid token"})
			return
		}
		userID, err := auth.ExtractSubFromToken(token)
		if err != nil {
			errors.HandleError(rw, errors.UnauthorizedError{Message: "Failed to extract user ID from token"})
			return
		}
		query := db.New(m.db)
		user, err := query.GetUserByID(r.Context(), userID)

		if err != nil {
			errors.HandleError(rw, err)
			return
		}
		ctx := context.WithValue(r.Context(), UserKey, user)

		next.ServeHTTP(rw, r.WithContext(ctx))
	}
}

func (m *MiddlewareManager) GroupPermissionMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		groupIDStr := r.PathValue("id")
		groupID, err := strconv.ParseInt(groupIDStr, 10, 32)
		log.Println()
		if err != nil {
			errors.HandleError(rw, err)
			return
		}

		user := r.Context().Value(UserKey).(db.User)

		query := db.New(m.db)
		group, err := query.GetGroupByID(r.Context(), int32(groupID))
		if err != nil {
			errors.HandleError(rw, err)
			return
		}

		if group.OwnerID.Int32 != user.ID {
			errors.HandleError(rw, errors.UnauthorizedError{Message: "User is not the group owner"})
			return
		}

		next.ServeHTTP(rw, r)
	}
}

func (m *MiddlewareManager) ErrorHandlerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				debug.PrintStack()
				log.Printf("Panic: %v", err)
				errors.SendErrorResponse(rw, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(rw, r)
	}
}

func MultipleMiddleware(h http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	if len(middlewares) < 1 {
		return h
	}
	wrapped := h
	for i := len(middlewares) - 1; i >= 0; i-- {
		wrapped = middlewares[i](wrapped)

	}
	return wrapped
}
