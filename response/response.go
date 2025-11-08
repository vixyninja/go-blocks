package response

import (
	"encoding/json"
	"net/http"
)

// writeJSON writes a JSON response with the given status code.
func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

// RespondOK sends a 200 OK response with the given data.
func RespondOK[T any](w http.ResponseWriter, r *http.Request, data T) error {
	return writeJSON(w, http.StatusOK, Response[T]{Data: data})
}

// RespondCreated sends a 201 Created response with the given data.
func RespondCreated[T any](w http.ResponseWriter, r *http.Request, data T) error {
	return writeJSON(w, http.StatusCreated, Response[T]{Data: data})
}

// RespondAccepted sends a 202 Accepted response with the given data.
func RespondAccepted[T any](w http.ResponseWriter, r *http.Request, data T) error {
	return writeJSON(w, http.StatusAccepted, Response[T]{Data: data})
}

// RespondNoContent sends a 204 No Content response.
func RespondNoContent(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusNoContent)
	return nil
}

// RespondPaged sends a 200 OK response with paginated data.
func RespondPaged[T any](w http.ResponseWriter, r *http.Request, data T, meta PageMeta) error {
	return writeJSON(w, http.StatusOK, PageResponse[T]{Data: data, Meta: meta})
}

// RespondError sends an error response with the given status code and error details.
func RespondError(w http.ResponseWriter, r *http.Request, status int, code, message string, details any) error {
	return writeJSON(w, status, ErrorResponse{
		Code:    code,
		Message: message,
		Details: details,
	})
}
