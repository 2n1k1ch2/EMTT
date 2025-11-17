package handlers

import (
	"EMTT/internal/dto"
	"EMTT/internal/repository"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
)

type Handler struct {
	SubRepo repository.SubscriptionRepository
}

var log = slog.Default()

// ========================== LIST ==========================

// ListSubscriptions
// @Summary List subscriptions
// @Tags subscriptions
// @Produce json
// @Success 200 {array} dto.SubscriptionDTO
// @Failure 500 {string} string "Internal server error"
// @Router /subscriptions/list [get]
func (h *Handler) ListSubscriptions(w http.ResponseWriter, r *http.Request) {
	subs, err := h.SubRepo.List()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(dto.SubscriptionsToDTO(subs))
}

// ========================== CREATE ==========================

// CreateSubscription
// @Summary Create subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body dto.SubscriptionDTO true "Subscription data"
// @Success 201 {object} dto.SubscriptionDTO
// @Failure 400 {string} string "Invalid JSON"
// @Failure 500 {string} string "Internal server error"
// @Router /subscriptions/create [post]
func (h *Handler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	var req dto.SubscriptionDTO

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	sub, err := dto.DTOToSubscription(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.SubRepo.Create(sub); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.SubscriptionToDTO(sub))
}

// ========================== GET BY ID ==========================

// GetSubscription
// @Summary Get subscription by ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 200 {object} dto.SubscriptionDTO
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Not found"
// @Router /subscriptions/get/{id} [get]
func (h *Handler) GetSubscription(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/subscriptions/get/")
	id, err := dto.StringToUUID(idStr)
	if err != nil {
		http.Error(w, "invalid id format", http.StatusBadRequest)
		return
	}

	sub, err := h.SubRepo.GetById(id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(dto.SubscriptionToDTO(sub))
}

// ========================== UPDATE ==========================

// UpdateSubscription
// @Summary Update subscription by ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID"
// @Param subscription body dto.SubscriptionDTO true "Updated subscription"
// @Success 200 {object} dto.SubscriptionDTO
// @Failure 400 {string} string "Invalid JSON"
// @Failure 500 {string} string "Internal server error"
// @Router /subscriptions/update/{id} [put]
func (h *Handler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/subscriptions/update/")
	id, err := dto.StringToUUID(idStr)
	if err != nil {
		http.Error(w, "invalid id format", http.StatusBadRequest)
		return
	}

	var req dto.SubscriptionDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	sub, err := dto.DTOToSubscription(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sub.ID = id

	if err := h.SubRepo.Update(sub); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(dto.SubscriptionToDTO(sub))
}

// ========================== DELETE ==========================

// DeleteSubscription
// @Summary Delete subscription by ID
// @Tags subscriptions
// @Param id path string true "Subscription ID"
// @Success 204 {string} string "No Content"
// @Failure 500 {string} string "Internal server error"
// @Router /subscriptions/delete/{id} [delete]
func (h *Handler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/subscriptions/delete/")
	id, err := dto.StringToUUID(idStr)
	if err != nil {
		http.Error(w, "invalid id format", http.StatusBadRequest)
		return
	}

	if err := h.SubRepo.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ========================== TOTAL COST ==========================

// TotalCost
// @Summary Total subscription cost
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "User ID"
// @Param service_name query string false "Service name"
// @Param from query string false "From date (YYYY-MM-DD)"
// @Param to query string false "To date (YYYY-MM-DD)"
// @Success 200 {object} map[string]uint32
// @Failure 500 {string} string "Internal server error"
// @Router /subscriptions/total [get]
func (h *Handler) TotalCost(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	service := r.URL.Query().Get("service_name")
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	total, err := h.SubRepo.TotalCost(userID, service, from, to)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]uint32{"total": total})
}
