package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"bot/internal/models"
	"bot/pkg/logger"
)


type VotingService interface {
	CreatePoll(req models.CreatePollRequest, creatorID string) (models.Poll, error)
	GetPoll(id string) (models.Poll, error)
	GetPollResults(id string) (models.PollResult, error)
	Vote(pollID string, req models.VoteRequest) error
	EndPoll(pollID string, userID string) error
	DeletePoll(pollID string, userID string) error
}


type Handler struct {
	service VotingService
	logger  logger.Logger
}


func NewHandler(service VotingService, logger logger.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}


func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}


func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, models.APIResponse{
		Success: false,
		Message: message,
	})
}


func (h *Handler) CreatePollHandler(w http.ResponseWriter, r *http.Request) {
	var req models.CreatePollRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		h.logger.Error("Error decoding request body", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()


	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		respondWithError(w, http.StatusBadRequest, "X-User-ID header is required")
		return
	}

	poll, err := h.service.CreatePoll(req, userID)
	if err != nil {
		h.logger.Error("Error creating poll", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.logger.Info("Poll created successfully with ID: " + poll.ID)
	respondWithJSON(w, http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Poll created successfully",
		Data:    poll,
	})
}


func (h *Handler) GetPollHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pollID := vars["id"]

	poll, err := h.service.GetPoll(pollID)
	if err != nil {
		h.logger.Error("Error getting poll", err)
		respondWithError(w, http.StatusNotFound, "Poll not found")
		return
	}

	respondWithJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    poll,
	})
}


func (h *Handler) GetPollResultsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pollID := vars["id"]

	results, err := h.service.GetPollResults(pollID)
	if err != nil {
		h.logger.Error("Error getting poll results", err)
		respondWithError(w, http.StatusNotFound, "Poll not found")
		return
	}

	respondWithJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    results,
	})
}


func (h *Handler) VoteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pollID := vars["id"]

	var req models.VoteRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		h.logger.Error("Error decoding request body", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()


	if req.UserID == "" {
		userID := r.Header.Get("X-User-ID")
		if userID == "" {
			respondWithError(w, http.StatusBadRequest, "User ID is required")
			return
		}
		req.UserID = userID
	}

	if err := h.service.Vote(pollID, req); err != nil {
		h.logger.Error("Error voting", err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Vote recorded successfully",
	})
}


func (h *Handler) EndPollHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pollID := vars["id"]


	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		respondWithError(w, http.StatusBadRequest, "X-User-ID header is required")
		return
	}

	if err := h.service.EndPoll(pollID, userID); err != nil {
		h.logger.Error("Error ending poll", err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Poll ended successfully",
	})
}


func (h *Handler) DeletePollHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pollID := vars["id"]


	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		respondWithError(w, http.StatusBadRequest, "X-User-ID header is required")
		return
	}

	if err := h.service.DeletePoll(pollID, userID); err != nil {
		h.logger.Error("Error deleting poll", err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Poll deleted successfully",
	})
}