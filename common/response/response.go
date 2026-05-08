package response

import (
	"encoding/json"
	"net/http"
)

// APIResponse adalah format standar untuk seluruh balasan HTTP di sistem ini
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

// Success mengirimkan respons JSON standar untuk status berhasil
func Success(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error mengirimkan respons JSON standar untuk status gagal/error
func Error(w http.ResponseWriter, statusCode int, message string, errors interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(APIResponse{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}
