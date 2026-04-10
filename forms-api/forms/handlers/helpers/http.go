package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"forms/models"

	"github.com/gorilla/mux"
)

// EnsureMethod valida que el método HTTP sea el esperado
func EnsureMethod(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method != method {
		RespondError(w, http.StatusMethodNotAllowed, "Método no permitido")
		return false
	}
	return true
}

// ParseFormIDFromVars parsea el ID de formulario desde las variables de la URL
func ParseFormIDFromVars(w http.ResponseWriter, r *http.Request) (int64, bool) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	formID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "ID de formulario inválido")
		return 0, false
	}
	return formID, true
}

// DecodeJSON decodifica el JSON del request body
func DecodeJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		RespondError(w, http.StatusBadRequest, fmt.Sprintf("Error decodificando JSON: %v", err))
		return false
	}
	return true
}

// RespondJSON responde con un JSON
func RespondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// RespondError responde con un error en formato JSON
func RespondError(w http.ResponseWriter, status int, message string) {
	RespondJSON(w, status, models.ErrorResponse{
		Success: false,
		Message: message,
	})
}














