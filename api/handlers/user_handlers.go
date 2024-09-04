package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/joseph-gunnarsson/scheduling/api/errors"
	db "github.com/joseph-gunnarsson/scheduling/db/models"
	"github.com/joseph-gunnarsson/scheduling/internals/auth"
	"golang.org/x/crypto/bcrypt"
)

func (h *BaseHandler) CreateUserHandler(rw http.ResponseWriter, r *http.Request) {
	var newUser db.CreateUserParams
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		errors.HandleError(rw, errors.ValidationError{Message: "Invalid request body"})
		return
	}

	if newUser.Username == "" || newUser.Email == "" || newUser.PasswordHash == "" {
		errors.HandleError(rw, errors.ValidationError{Message: "Missing required fields"})
		return
	}

	newUser.PasswordHash, err = auth.HashPassword(newUser.PasswordHash)
	if err != nil {
		errors.HandleError(rw, err)
		return
	}

	query := db.New(h.db)
	user, err := query.CreateUser(r.Context(), newUser)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			errors.HandleError(rw, errors.ValidationError{Message: "Username or email already exists"})
		} else {
			errors.HandleError(rw, err)
		}
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	json.NewEncoder(rw).Encode(user)
}

func (h *BaseHandler) LoginHandler(rw http.ResponseWriter, r *http.Request) {
	var loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		errors.HandleError(rw, errors.ValidationError{Message: "Invalid request body: " + err.Error()})
		return
	}

	query := db.New(h.db)
	userInformation, err := query.LoginUser(r.Context(), loginRequest.Username)
	if err != nil {
		if err == pgx.ErrNoRows {
			errors.HandleError(rw, errors.UnauthorizedError{Message: "Invalid username or password"})
		} else {
			errors.HandleError(rw, err)
		}
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(userInformation.PasswordHash), []byte(loginRequest.Password))
	if err != nil {
		errors.HandleError(rw, errors.UnauthorizedError{Message: "Invalid username or password"})
		return
	}

	token, err := auth.GenerateJWTToken(userInformation.ID, userInformation.Username)
	if err != nil {
		errors.HandleError(rw, err)
		return
	}

	response := map[string]string{
		"token": token,
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(response)
}

func (h *BaseHandler) UpdatePassword(rw http.ResponseWriter, r *http.Request) {
	var updatePasswordRequest struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	err := json.NewDecoder(r.Body).Decode(&updatePasswordRequest)
	if err != nil {
		errors.HandleError(rw, errors.ValidationError{Message: "Invalid request body"})
		return
	}

	userID, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		errors.HandleError(rw, errors.ValidationError{Message: "Invalid user ID"})
		return
	}

	query := db.New(h.db)
	user, err := query.GetUserByID(r.Context(), int32(userID))
	if err != nil {
		if err == pgx.ErrNoRows {
			errors.HandleError(rw, errors.NotFoundError{Message: "User not found"})
		} else {
			errors.HandleError(rw, err)
		}
		return
	}
	log.Println(updatePasswordRequest.OldPassword)
	log.Println(user.PasswordHash)
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(updatePasswordRequest.OldPassword))
	if err != nil {
		errors.HandleError(rw, errors.UnauthorizedError{Message: "Invalid old password"})
		return
	}

	newPasswordHash, err := auth.HashPassword(updatePasswordRequest.NewPassword)
	if err != nil {
		errors.HandleError(rw, err)
		return
	}

	err = query.UpdateUserPassword(r.Context(), db.UpdateUserPasswordParams{
		ID:           int32(userID),
		PasswordHash: newPasswordHash,
	})
	if err != nil {
		errors.HandleError(rw, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(map[string]string{"message": "Password updated successfully"})
}
