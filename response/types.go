package response

import "encoding/json"

// Response represents a standard API response wrapper.
type Response[T any] struct {
	Data T `json:"data"`
}

// PageMeta contains pagination metadata.
type PageMeta struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Count  int `json:"count"`
}

// PageResponse represents a paginated API response.
type PageResponse[T any] struct {
	Data T        `json:"data"`
	Meta PageMeta `json:"meta"`
}

// ErrorResponse represents an error API response.
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

// Encoder defines the interface for encoding responses.
type Encoder interface {
	Encode(v any) error
}

// JSONEncoder is the default JSON encoder.
type JSONEncoder struct {
	encoder *json.Encoder
}

// NewJSONEncoder creates a new JSON encoder.
func NewJSONEncoder(encoder *json.Encoder) *JSONEncoder {
	return &JSONEncoder{encoder: encoder}
}

// Encode encodes the value to JSON.
func (e *JSONEncoder) Encode(v any) error {
	return e.encoder.Encode(v)
}
