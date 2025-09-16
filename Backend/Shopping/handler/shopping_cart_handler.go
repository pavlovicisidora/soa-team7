package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pavlovicisidora/soa-team7/Backend/Shopping/service"
)

type ShoppingCartHandler struct {
	service *service.ShoppingCartService
}

func NewShoppingCartHandler(service *service.ShoppingCartService) *ShoppingCartHandler {
	return &ShoppingCartHandler{service: service}
}

func (h *ShoppingCartHandler) GetCartHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "X-User-ID header is missing", http.StatusUnauthorized)
		return
	}

	cart, err := h.service.GetCart(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to retrieve cart", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cart)
}

func (h *ShoppingCartHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "X-User-ID header is missing", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	tourIDStr, ok := vars["tourId"]
	if !ok || tourIDStr == "" {
		http.Error(w, "Tour ID is missing in the URL path", http.StatusBadRequest)
		return
	}

	tourId, err := strconv.ParseInt(tourIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid Tour ID format in URL path", http.StatusBadRequest)
		return
	}

	cart, err := h.service.AddItemToCart(r.Context(), userID, int(tourId))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cart)
}

func (h *ShoppingCartHandler) RemoveFromCartHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	tourIDStr := mux.Vars(r)["tourId"]
	tourID, err := strconv.ParseInt(tourIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid Tour ID format in URL path", http.StatusBadRequest)
		return
	}
	cart, err := h.service.RemoveItemFromCart(r.Context(), userID, int(tourID))
	if err != nil {
		http.Error(w, "Failed to remove item from cart", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cart)
}

func (h *ShoppingCartHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")

	tokens, err := h.service.Checkout(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tokens)
}

func (h *ShoppingCartHandler) CheckTokenHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	tourIDStr := mux.Vars(r)["tourId"]
	tourID, err := strconv.ParseInt(tourIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid Tour ID format in URL path", http.StatusBadRequest)
		return
	}

	hasToken, err := h.service.CheckToken(r.Context(), userID, int(tourID))
	if err != nil {
		http.Error(w, "Failed to check token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]bool{"hasToken": hasToken})
}
