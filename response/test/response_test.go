package response_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vixyninja/go-blocks/response"
)

func TestRespondOK(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	err := response.RespondOK(w, r, "test")
	if err != nil {
		t.Fatalf("RespondOK() error = %v", err)
	}

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp response.Response[string]
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Data != "test" {
		t.Fatalf("expected data 'test', got %q", resp.Data)
	}
}

func TestRespondCreated(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)

	err := response.RespondCreated(w, r, map[string]string{"id": "123"})
	if err != nil {
		t.Fatalf("RespondCreated() error = %v", err)
	}

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, w.Code)
	}
}

func TestRespondAccepted(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)

	err := response.RespondAccepted(w, r, map[string]string{"status": "processing"})
	if err != nil {
		t.Fatalf("RespondAccepted() error = %v", err)
	}

	if w.Code != http.StatusAccepted {
		t.Fatalf("expected status %d, got %d", http.StatusAccepted, w.Code)
	}
}

func TestRespondNoContent(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/", nil)

	err := response.RespondNoContent(w, r)
	if err != nil {
		t.Fatalf("RespondNoContent() error = %v", err)
	}

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, w.Code)
	}
}

func TestRespondPaged(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	data := []string{"a", "b", "c"}
	meta := response.PageMeta{Limit: 10, Offset: 0, Count: 3}

	err := response.RespondPaged(w, r, data, meta)
	if err != nil {
		t.Fatalf("RespondPaged() error = %v", err)
	}

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp response.PageResponse[[]string]
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Meta.Count != 3 {
		t.Fatalf("expected count 3, got %d", resp.Meta.Count)
	}
}

func TestRespondError(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	err := response.RespondError(w, r, http.StatusBadRequest, "test_code", "test message", nil)
	if err != nil {
		t.Fatalf("RespondError() error = %v", err)
	}

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var resp response.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Code != "test_code" {
		t.Fatalf("expected code 'test_code', got %q", resp.Code)
	}
}

func TestBadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	err := response.BadRequest(w, r, map[string]string{"field": "error"})
	if err != nil {
		t.Fatalf("BadRequest() error = %v", err)
	}

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUnauthorized(t *testing.T) {
	cases := []struct {
		name    string
		message string
		wantMsg string
	}{
		{"with_message", "Custom unauthorized", "Custom unauthorized"},
		{"empty_message", "", "Unauthorized"},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)

			err := response.Unauthorized(w, r, cs.message)
			if err != nil {
				t.Fatalf("Unauthorized() error = %v", err)
			}

			if w.Code != http.StatusUnauthorized {
				t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
			}

			var resp response.ErrorResponse
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("failed to unmarshal: %v", err)
			}

			if resp.Message != cs.wantMsg {
				t.Fatalf("expected message %q, got %q", cs.wantMsg, resp.Message)
			}
		})
	}
}

func TestNotFound(t *testing.T) {
	cases := []struct {
		name    string
		message string
		wantMsg string
	}{
		{"with_message", "Custom not found", "Custom not found"},
		{"empty_message", "", "Resource not found"},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)

			err := response.NotFound(w, r, cs.message)
			if err != nil {
				t.Fatalf("NotFound() error = %v", err)
			}

			if w.Code != http.StatusNotFound {
				t.Fatalf("expected status %d, got %d", http.StatusNotFound, w.Code)
			}
		})
	}
}

func TestInternalServerError(t *testing.T) {
	cases := []struct {
		name string
		err  error
	}{
		{"with_error", &testError{msg: "test error"}},
		{"nil_error", nil},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)

			err := response.InternalServerError(w, r, cs.err)
			if err != nil {
				t.Fatalf("InternalServerError() error = %v", err)
			}

			if w.Code != http.StatusInternalServerError {
				t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
			}
		})
	}
}

func TestUnprocessableEntity(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	fieldErrors := map[string][]string{
		"email": {"invalid format"},
		"name":  {"required"},
	}

	err := response.UnprocessableEntity(w, r, fieldErrors)
	if err != nil {
		t.Fatalf("UnprocessableEntity() error = %v", err)
	}

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected status %d, got %d", http.StatusUnprocessableEntity, w.Code)
	}

	var resp response.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if resp.Details == nil {
		t.Fatal("expected details to be set")
	}
}

func TestTooManyRequests(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	err := response.TooManyRequests(w, r, "Rate limit exceeded")
	if err != nil {
		t.Fatalf("TooManyRequests() error = %v", err)
	}

	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("expected status %d, got %d", http.StatusTooManyRequests, w.Code)
	}
}

func TestServiceUnavailable(t *testing.T) {
	cases := []struct {
		name    string
		message string
		wantMsg string
	}{
		{"with_message", "Custom unavailable", "Custom unavailable"},
		{"empty_message", "", "Service unavailable"},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)

			err := response.ServiceUnavailable(w, r, cs.message)
			if err != nil {
				t.Fatalf("ServiceUnavailable() error = %v", err)
			}

			if w.Code != http.StatusServiceUnavailable {
				t.Fatalf("expected status %d, got %d", http.StatusServiceUnavailable, w.Code)
			}
		})
	}
}

func TestAllErrorStatusCodes(t *testing.T) {
	tests := []struct {
		name     string
		fn       func(http.ResponseWriter, *http.Request, string) error
		expected int
	}{
		{"PaymentRequired", response.PaymentRequired, http.StatusPaymentRequired},
		{"Forbidden", response.Forbidden, http.StatusForbidden},
		{"MethodNotAllowed", response.MethodNotAllowed, http.StatusMethodNotAllowed},
		{"NotAcceptable", response.NotAcceptable, http.StatusNotAcceptable},
		{"ProxyAuthRequired", response.ProxyAuthRequired, http.StatusProxyAuthRequired},
		{"RequestTimeout", response.RequestTimeout, http.StatusRequestTimeout},
		{"Conflict", response.Conflict, http.StatusConflict},
		{"Gone", response.Gone, http.StatusGone},
		{"LengthRequired", response.LengthRequired, http.StatusLengthRequired},
		{"PreconditionFailed", response.PreconditionFailed, http.StatusPreconditionFailed},
		{"RequestEntityTooLarge", response.RequestEntityTooLarge, http.StatusRequestEntityTooLarge},
		{"RequestURITooLong", response.RequestURITooLong, http.StatusRequestURITooLong},
		{"UnsupportedMediaType", response.UnsupportedMediaType, http.StatusUnsupportedMediaType},
		{"RequestedRangeNotSatisfiable", response.RequestedRangeNotSatisfiable, http.StatusRequestedRangeNotSatisfiable},
		{"ExpectationFailed", response.ExpectationFailed, http.StatusExpectationFailed},
		{"Teapot", response.Teapot, http.StatusTeapot},
		{"MisdirectedRequest", response.MisdirectedRequest, http.StatusMisdirectedRequest},
		{"Locked", response.Locked, http.StatusLocked},
		{"FailedDependency", response.FailedDependency, http.StatusFailedDependency},
		{"TooEarly", response.TooEarly, http.StatusTooEarly},
		{"UpgradeRequired", response.UpgradeRequired, http.StatusUpgradeRequired},
		{"PreconditionRequired", response.PreconditionRequired, http.StatusPreconditionRequired},
		{"RequestHeaderFieldsTooLarge", response.RequestHeaderFieldsTooLarge, http.StatusRequestHeaderFieldsTooLarge},
		{"UnavailableForLegalReasons", response.UnavailableForLegalReasons, http.StatusUnavailableForLegalReasons},
		{"NotImplemented", response.NotImplemented, http.StatusNotImplemented},
		{"BadGateway", response.BadGateway, http.StatusBadGateway},
		{"GatewayTimeout", response.GatewayTimeout, http.StatusGatewayTimeout},
		{"HTTPVersionNotSupported", response.HTTPVersionNotSupported, http.StatusHTTPVersionNotSupported},
		{"VariantAlsoNegotiates", response.VariantAlsoNegotiates, http.StatusVariantAlsoNegotiates},
		{"InsufficientStorage", response.InsufficientStorage, http.StatusInsufficientStorage},
		{"LoopDetected", response.LoopDetected, http.StatusLoopDetected},
		{"NotExtended", response.NotExtended, http.StatusNotExtended},
		{"NetworkAuthenticationRequired", response.NetworkAuthenticationRequired, http.StatusNetworkAuthenticationRequired},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)

			err := tt.fn(w, r, "")
			if err != nil {
				t.Fatalf("%s() error = %v", tt.name, err)
			}

			if w.Code != tt.expected {
				t.Fatalf("expected status %d, got %d", tt.expected, w.Code)
			}
		})
	}
}

type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}
