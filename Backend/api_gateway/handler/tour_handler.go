package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/pavlovicisidora/soa-team7/Backend/APIGateway/middleware"
	tour_proto "github.com/pavlovicisidora/soa-team7/Backend/Tour/src/main/proto/tour"
)

// Definišemo strukturu za dolazeći JSON zahtev
type CreateTourRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Difficulty  string `json:"difficulty"`
	Tags        string `json:"tags"`
}

// TourHandler sadrži gRPC klijent za Java tour-service
type TourHandler struct {
	client tour_proto.TourGrpcServiceClient
}

// Konstruktor za TourHandler
func NewTourHandler(client tour_proto.TourGrpcServiceClient) *TourHandler {
	return &TourHandler{client: client}
}

// ServeHTTP registruje rute SAMO za ture
func (h *TourHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := mux.NewRouter()

	// Definišemo samo rute koje se tiču tura
	router.HandleFunc("/create", h.CreateTourHandler).Methods("POST")
	// Ovde biste dodali i ostale rute za ture, npr. /tours/{id} itd.

	router.ServeHTTP(w, r)
}

// CreateTourHandler je funkcija koja obrađuje HTTP zahtev
// Primetite kako je ovaj kod skoro identičan BlogHandler-u.
func (h *TourHandler) CreateTourHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Koristimo postojeći middleware da dobijemo userID
	userID, ok := r.Context().Value(middleware.UserKey).(string)
	if !ok || userID == "" {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	// 2. Čitamo JSON telo zahteva
	var reqBody CreateTourRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 3. Kreiramo gRPC zahtev od HTTP podataka
	grpcRequest := &tour_proto.CreateTourRequest{
		Name:        reqBody.Name,
		Description: reqBody.Description,
		Difficulty:  reqBody.Difficulty,
		Tags:        reqBody.Tags,
	}

	// 4. Pozivamo gRPC metodu na Java servisu
	// U ovom slučaju ne prosleđujemo userID u telu, jer će Java interceptor to izvući iz tokena
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// KAKO PROSLEDITI TOKEN? Koristimo postojeći AuthMiddleware,
	// ali moramo i proslediti header dalje ka Javi.
	// (Vidi izmene u main.go za ovo)

	resp, err := h.client.CreateTour(ctx, grpcRequest)
	if err != nil {
		log.Printf("Failed to create tour via gRPC: %v", err)
		http.Error(w, "Failed to create tour", http.StatusInternalServerError)
		return
	}

	// 5. Vraćamo odgovor klijentu
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Failed to encode tour response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
