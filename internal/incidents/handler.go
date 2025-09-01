package incidents

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	Svc *Service
}

func RegisterRoutes(r chi.Router, h *Handler) {
	r.Route("/incidents", func(r chi.Router) {
		r.Post("/create", h.create)
		r.Get("/get", h.list)
	})
}

// @Summary Create an incident
// @Description Create a new incident with AI classification
// @Tags incidents
// @Accept json
// @Produce json
// @Param request body Incident true "Incident data"
// @Success 200 {object} Incident
// @Router /incidents/create [post]
func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var req CreateIncidentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	inc, err := h.Svc.Create(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(inc)
}

// @Summary List incidents
// @Description Get a list of all incidents
// @Tags incidents
// @Produce json
// @Success 200 {array} Incident
// @Router /incidents/get [get]
func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	incs, err := h.Svc.List(r.Context())
	if err != nil {
		http.Error(w, "failed to fetch", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(incs)
}
