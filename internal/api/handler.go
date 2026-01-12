package api

import (
	"encoding/json"
	"net/http"

	"github.com/javinjeetoo/pack-calculator-api/internal/packs"
)

const maxItems = 1_000_000

type CalculateRequest struct {
	Items     int   `json:"items"`
	PackSizes []int `json:"pack_sizes,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type Handler struct {
	DefaultPackSizes []int
}

func NewHandler(defaultPackSizes []int) *Handler {
	return &Handler{DefaultPackSizes: defaultPackSizes}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Handler) Calculate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "method not allowed"})
		return
	}

	var req CalculateRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid JSON: " + err.Error()})
		return
	}

	// Input validation / guardrail (protect DP solver)
	if req.Items <= 0 {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "items must be > 0"})
		return
	}
	if req.Items > maxItems {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "items must be <= 1000000"})
		return
	}

	sizes := h.DefaultPackSizes
	if len(req.PackSizes) > 0 {
		sizes = req.PackSizes
	}

	res, err := packs.Solve(req.Items, sizes)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, res)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
