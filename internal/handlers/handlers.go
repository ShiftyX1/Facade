package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ShiftyX1/Facade/internal/state"
)


type StateHandler struct {
	stateManager *state.Manager
}


func NewStateHandler(sm *state.Manager) *StateHandler {
	return &StateHandler{stateManager: sm}
}


func (h *StateHandler) GetState(w http.ResponseWriter, r *http.Request) {
	data := h.stateManager.GetAll()
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}


func (h *StateHandler) ClearState(w http.ResponseWriter, r *http.Request) {
	h.stateManager.Clear()
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "cleared"})
}


func (h *StateHandler) GetStateKeys(w http.ResponseWriter, r *http.Request) {
	keys := h.stateManager.Keys()
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"keys": keys})
}


func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"service": "facade",
	})
}