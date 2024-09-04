package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/joseph-gunnarsson/scheduling/api/errors"
	db "github.com/joseph-gunnarsson/scheduling/db/models"
)

func (h BaseHandler) CreateGroupHandler(rw http.ResponseWriter, r *http.Request) {
	var newGroup db.CreateGroupParams

	err := json.NewDecoder(r.Body).Decode(&newGroup)

	if err != nil {
		errors.HandleError(rw, errors.ValidationError{Message: "Invalid request body"})
		return
	}

	query := db.New(h.db)
	group, err := query.CreateGroup(r.Context(), newGroup)
	log.Printf("%d", newGroup.OwnerID.Int32)
	if err != nil {
		errors.HandleError(rw, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	json.NewEncoder(rw).Encode(group)
}

func (h BaseHandler) DeleteGroupHandler(rw http.ResponseWriter, r *http.Request) {
	groupIDstr := r.PathValue("id")

	groupID, err := strconv.ParseInt(groupIDstr, 10, 32)
	if err != nil {
		errors.HandleError(rw, errors.ValidationError{Message: "Invalid or missing group id."})
	}

	query := db.New(h.db)
	err = query.DeleteGroup(r.Context(), int32(groupID))

	if err != nil {
		errors.HandleError(rw, err)
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	json.NewEncoder(rw).Encode(map[string]string{"message": "Deleted group successfully"})
}

func (h *BaseHandler) UpdateGroupHandler(rw http.ResponseWriter, r *http.Request) {
	groupIDStr := r.PathValue("id")
	groupID, err := strconv.ParseInt(groupIDStr, 10, 32)
	if err != nil {
		errors.HandleError(rw, errors.ValidationError{Message: "Invalid group id"})
		return
	}

	var updateGroup db.UpdateGroupParams
	err = json.NewDecoder(r.Body).Decode(&updateGroup)
	if err != nil {
		errors.HandleError(rw, errors.ValidationError{Message: "Invalid request body"})
		return
	}
	updateGroup.ID = int32(groupID)

	query := db.New(h.db)
	group, err := query.UpdateGroup(r.Context(), updateGroup)
	if err != nil {
		errors.HandleError(rw, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(group)
}

func (h *BaseHandler) PatchGroupHandler(rw http.ResponseWriter, r *http.Request) {
	groupIDStr := r.PathValue("id")
	groupID, err := strconv.ParseInt(groupIDStr, 10, 32)
	if err != nil {
		errors.HandleError(rw, errors.ValidationError{Message: "Invalid group id"})
		return
	}

	var patchData map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&patchData)
	if err != nil {
		errors.HandleError(rw, errors.ValidationError{Message: "Invalid request body"})
		return
	}

	query := db.New(h.db)
	var patchGroup db.PatchGroupParams
	patchGroup.ID = int32(groupID)

	if name, ok := patchData["name"].(string); ok {
		patchGroup.Name = pgtype.Text{String: name, Valid: true}
	}

	if description, ok := patchData["description"].(string); ok {
		patchGroup.Description = pgtype.Text{String: description, Valid: true}
	}

	group, err := query.PatchGroup(r.Context(), patchGroup)
	if err != nil {
		errors.HandleError(rw, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(group)
}

func (h *BaseHandler) GetGroupsByOwnerHandler(rw http.ResponseWriter, r *http.Request) {
	ownerIDStr := r.PathValue("id")
	ownerID, err := strconv.ParseInt(ownerIDStr, 10, 32)
	if err != nil {
		errors.HandleError(rw, errors.ValidationError{Message: "Invalid owner id"})
		return
	}

	query := db.New(h.db)

	groups, err := query.GetGroupsByOwner(r.Context(), pgtype.Int4{Int32: int32(ownerID),
		Valid: true,
	})
	if err != nil {
		errors.HandleError(rw, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(groups)
}

func (h *BaseHandler) AddUserToGroupHandler(rw http.ResponseWriter, r *http.Request) {
	var addUserToGroup db.AddUserToGroupParams
	err := json.NewDecoder(r.Body).Decode(&addUserToGroup)
	if err != nil {
		errors.HandleError(rw, errors.ValidationError{Message: "Invalid request body"})
		return
	}

	query := db.New(h.db)
	userGroup, err := query.AddUserToGroup(r.Context(), addUserToGroup)
	if err != nil {
		errors.HandleError(rw, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	json.NewEncoder(rw).Encode(userGroup)
}

func (h *BaseHandler) DeleteUserFromGroupHandler(rw http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(r.PathValue("user_id"), 10, 32)
	if err != nil {
		errors.HandleError(rw, errors.ValidationError{Message: "Invalid user id"})
		return
	}

	groupID, err := strconv.ParseInt(r.PathValue("group_id"), 10, 32)
	if err != nil {
		errors.HandleError(rw, errors.ValidationError{Message: "Invalid group id"})
		return
	}

	query := db.New(h.db)
	err = query.DeleteUserFromGroup(r.Context(), db.DeleteUserFromGroupParams{
		UserID:  int32(userID),
		GroupID: int32(groupID),
	})
	if err != nil {
		errors.HandleError(rw, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(map[string]string{"message": "User removed from group successfully"})
}

func (h *BaseHandler) GetUserGroupsHandler(rw http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 32)
	if err != nil {
		errors.HandleError(rw, errors.ValidationError{Message: "Invalid user id"})
		return
	}

	query := db.New(h.db)
	userGroups, err := query.GetUserGroups(r.Context(), int32(userID))
	if err != nil {
		errors.HandleError(rw, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(userGroups)
}

func (h *BaseHandler) GetGroupMembersHandler(rw http.ResponseWriter, r *http.Request) {
	groupIDStr := r.PathValue("group_id")
	groupID, err := strconv.ParseInt(groupIDStr, 10, 32)
	if err != nil {
		errors.HandleError(rw, errors.ValidationError{Message: "Invalid group id"})
		return
	}

	query := db.New(h.db)
	groupMembers, err := query.GetGroupMembers(r.Context(), int32(groupID))
	if err != nil {
		errors.HandleError(rw, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(groupMembers)
}
