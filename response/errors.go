package response

import "net/http"

// BadRequest sends a 400 Bad Request response.
func BadRequest(w http.ResponseWriter, r *http.Request, details any) error {
	return RespondError(w, r, http.StatusBadRequest, "bad_request", "Invalid request", details)
}

// Unauthorized sends a 401 Unauthorized response.
func Unauthorized(w http.ResponseWriter, r *http.Request, message string) error {
	if message == "" {
		message = "Unauthorized"
	}
	return RespondError(w, r, http.StatusUnauthorized, "unauthorized", message, nil)
}

// PaymentRequired sends a 402 Payment Required response.
func PaymentRequired(w http.ResponseWriter, r *http.Request, message string) error {
	if message == "" {
		message = "Payment required"
	}
	return RespondError(w, r, http.StatusPaymentRequired, "payment_required", message, nil)
}

// Forbidden sends a 403 Forbidden response.
func Forbidden(w http.ResponseWriter, r *http.Request, message string) error {
	if message == "" {
		message = "Forbidden"
	}
	return RespondError(w, r, http.StatusForbidden, "forbidden", message, nil)
}

// NotFound sends a 404 Not Found response.
func NotFound(w http.ResponseWriter, r *http.Request, message string) error {
	if message == "" {
		message = "Resource not found"
	}
	return RespondError(w, r, http.StatusNotFound, "not_found", message, nil)
}

// MethodNotAllowed sends a 405 Method Not Allowed response.
func MethodNotAllowed(w http.ResponseWriter, r *http.Request, message string) error {
	if message == "" {
		message = "Method not allowed"
	}
	return RespondError(w, r, http.StatusMethodNotAllowed, "method_not_allowed", message, nil)
}

// NotAcceptable sends a 406 Not Acceptable response.
func NotAcceptable(w http.ResponseWriter, r *http.Request, message string) error {
	if message == "" {
		message = "Not acceptable"
	}
	return RespondError(w, r, http.StatusNotAcceptable, "not_acceptable", message, nil)
}

// ProxyAuthRequired sends a 407 Proxy Authentication Required response.
func ProxyAuthRequired(w http.ResponseWriter, r *http.Request, message string) error {
	if message == "" {
		message = "Proxy authentication required"
	}
	return RespondError(w, r, http.StatusProxyAuthRequired, "proxy_auth_required", message, nil)
}

// RequestTimeout sends a 408 Request Timeout response.
func RequestTimeout(w http.ResponseWriter, r *http.Request, message string) error {
	if message == "" {
		message = "Request timeout"
	}
	return RespondError(w, r, http.StatusRequestTimeout, "request_timeout", message, nil)
}

// Conflict sends a 409 Conflict response.
func Conflict(w http.ResponseWriter, r *http.Request, message string) error {
	if message == "" {
		message = "Conflict"
	}
	return RespondError(w, r, http.StatusConflict, "conflict", message, nil)
}

// Gone sends a 410 Gone response.
func Gone(w http.ResponseWriter, r *http.Request, message string) error {
	if message == "" {
		message = "Resource is gone"
	}
	return RespondError(w, r, http.StatusGone, "gone", message, nil)
}

// LengthRequired sends a 411 Length Required response.
func LengthRequired(w http.ResponseWriter, r *http.Request, message string) error {
	if message == "" {
		message = "Length required"
	}
	return RespondError(w, r, http.StatusLengthRequired, "length_required", message, nil)
}

// PreconditionFailed sends a 412 Precondition Failed response.
func PreconditionFailed(w http.ResponseWriter, r *http.Request, message string) error {
	if message == "" {
		message = "Precondition failed"
	}
	return RespondError(w, r, http.StatusPreconditionFailed, "precondition_failed", message, nil)
}

// RequestEntityTooLarge sends a 413 Request Entity Too Large response.
func RequestEntityTooLarge(w http.ResponseWriter, r *http.Request, message string) error {
	if message == "" {
		message = "Request entity too large"
	}
	return RespondError(w, r, http.StatusRequestEntityTooLarge, "request_entity_too_large", message, nil)
}

// RequestURITooLong sends a 414 Request URI Too Long response.
func RequestURITooLong(w http.ResponseWriter, r *http.Request, message string) error {
	if message == "" {
		message = "Request URI too long"
	}
	return RespondError(w, r, http.StatusRequestURITooLong, "request_uri_too_long", message, nil)
}

// UnsupportedMediaType sends a 415 Unsupported Media Type response.
func UnsupportedMediaType(w http.ResponseWriter, r *http.Request, message string) error {
	if message == "" {
		message = "Unsupported media type"
	}
	return RespondError(w, r, http.StatusUnsupportedMediaType, "unsupported_media_type", message, nil)
}

// RequestedRangeNotSatisfiable sends a 416 Range Not Satisfiable response.
func RequestedRangeNotSatisfiable(w http.ResponseWriter, r *http.Request, message string) error {
	if message == "" {
		message = "Requested range not satisfiable"
	}
	return RespondError(w, r, http.StatusRequestedRangeNotSatisfiable, "range_not_satisfiable", message, nil)
}

// ExpectationFailed sends a 417 Expectation Failed response.
func ExpectationFailed(w http.ResponseWriter, r *http.Request, message string) error {
	if message == "" {
		message = "Expectation failed"
	}
	return RespondError(w, r, http.StatusExpectationFailed, "expectation_failed", message, nil)
}

// Teapot sends a 418 I'm a teapot response.
func Teapot(w http.ResponseWriter, r *http.Request, message string) error {
	if message == "" {
		message = "I'm a teapot"
	}
	return RespondError(w, r, http.StatusTeapot, "teapot", message, nil)
}

// MisdirectedRequest sends a 421 Misdirected Request response.
func MisdirectedRequest(w http.ResponseWriter, r *http.Request, message string) error {
	if message == "" {
		message = "Misdirected request"
	}
	return RespondError(w, r, http.StatusMisdirectedRequest, "misdirected_request", message, nil)
}

// UnprocessableEntity sends a 422 Unprocessable Entity response with validation errors.
func UnprocessableEntity(w http.ResponseWriter, r *http.Request, fieldErrors map[string][]string) error {
	return RespondError(w, r, http.StatusUnprocessableEntity, "validation_error", "Validation error", fieldErrors)
}

// Locked sends a 423 Locked response.
func Locked(w http.ResponseWriter, r *http.Request, message string) error {
	if message == "" {
		message = "Resource is locked"
	}
	return RespondError(w, r, http.StatusLocked, "locked", message, nil)
}

// FailedDependency sends a 424 Failed Dependency response.
func FailedDependency(w http.ResponseWriter, r *http.Request, message string) error {
	if message == "" {
		message = "Failed dependency"
	}
	return RespondError(w, r, http.StatusFailedDependency, "failed_dependency", message, nil)
}

// TooEarly sends a 425 Too Early response.
func TooEarly(w http.ResponseWriter, r *http.Request, message string) error {
	if message == "" {
		message = "Too early"
	}
	return RespondError(w, r, http.StatusTooEarly, "too_early", message, nil)
}

// UpgradeRequired sends a 426 Upgrade Required response.
func UpgradeRequired(w http.ResponseWriter, r *http.Request, message string) error {
	if message == "" {
		message = "Upgrade required"
	}
	return RespondError(w, r, http.StatusUpgradeRequired, "upgrade_required", message, nil)
}

// PreconditionRequired sends a 428 Precondition Required response.
func PreconditionRequired(w http.ResponseWriter, r *http.Request, message string) error {
	if message == "" {
		message = "Precondition required"
	}
	return RespondError(w, r, http.StatusPreconditionRequired, "precondition_required", message, nil)
}

// TooManyRequests sends a 429 Too Many Requests response.
func TooManyRequests(w http.ResponseWriter, r *http.Request, message string) error {
	if message == "" {
		message = "Too many requests"
	}
	return RespondError(w, r, http.StatusTooManyRequests, "too_many_requests", message, nil)
}

// RequestHeaderFieldsTooLarge sends a 431 Request Header Fields Too Large response.
func RequestHeaderFieldsTooLarge(w http.ResponseWriter, r *http.Request, message string) error {
	if message == "" {
		message = "Request header fields too large"
	}
	return RespondError(w, r, http.StatusRequestHeaderFieldsTooLarge, "request_header_fields_too_large", message, nil)
}

// UnavailableForLegalReasons sends a 451 Unavailable For Legal Reasons response.
func UnavailableForLegalReasons(w http.ResponseWriter, r *http.Request, message string) error {
	if message == "" {
		message = "Unavailable for legal reasons"
	}
	return RespondError(w, r, http.StatusUnavailableForLegalReasons, "unavailable_for_legal_reasons", message, nil)
}

// InternalServerError sends a 500 Internal Server Error response.
func InternalServerError(w http.ResponseWriter, r *http.Request, err error) error {
	msg := "Internal server error"
	if err != nil {
		msg = err.Error()
	}
	return RespondError(w, r, http.StatusInternalServerError, "internal_error", msg, nil)
}

// NotImplemented sends a 501 Not Implemented response.
func NotImplemented(w http.ResponseWriter, r *http.Request, message string) error {
	if message == "" {
		message = "Not implemented"
	}
	return RespondError(w, r, http.StatusNotImplemented, "not_implemented", message, nil)
}

// BadGateway sends a 502 Bad Gateway response.
func BadGateway(w http.ResponseWriter, r *http.Request, message string) error {
	if message == "" {
		message = "Bad gateway"
	}
	return RespondError(w, r, http.StatusBadGateway, "bad_gateway", message, nil)
}

// ServiceUnavailable sends a 503 Service Unavailable response.
func ServiceUnavailable(w http.ResponseWriter, r *http.Request, message string) error {
	if message == "" {
		message = "Service unavailable"
	}
	return RespondError(w, r, http.StatusServiceUnavailable, "service_unavailable", message, nil)
}

// GatewayTimeout sends a 504 Gateway Timeout response.
func GatewayTimeout(w http.ResponseWriter, r *http.Request, message string) error {
	if message == "" {
		message = "Gateway timeout"
	}
	return RespondError(w, r, http.StatusGatewayTimeout, "gateway_timeout", message, nil)
}

// HTTPVersionNotSupported sends a 505 HTTP Version Not Supported response.
func HTTPVersionNotSupported(w http.ResponseWriter, r *http.Request, message string) error {
	if message == "" {
		message = "HTTP version not supported"
	}
	return RespondError(w, r, http.StatusHTTPVersionNotSupported, "http_version_not_supported", message, nil)
}

// VariantAlsoNegotiates sends a 506 Variant Also Negotiates response.
func VariantAlsoNegotiates(w http.ResponseWriter, r *http.Request, message string) error {
	if message == "" {
		message = "Variant also negotiates"
	}
	return RespondError(w, r, http.StatusVariantAlsoNegotiates, "variant_also_negotiates", message, nil)
}

// InsufficientStorage sends a 507 Insufficient Storage response.
func InsufficientStorage(w http.ResponseWriter, r *http.Request, message string) error {
	if message == "" {
		message = "Insufficient storage"
	}
	return RespondError(w, r, http.StatusInsufficientStorage, "insufficient_storage", message, nil)
}

// LoopDetected sends a 508 Loop Detected response.
func LoopDetected(w http.ResponseWriter, r *http.Request, message string) error {
	if message == "" {
		message = "Loop detected"
	}
	return RespondError(w, r, http.StatusLoopDetected, "loop_detected", message, nil)
}

// NotExtended sends a 510 Not Extended response.
func NotExtended(w http.ResponseWriter, r *http.Request, message string) error {
	if message == "" {
		message = "Not extended"
	}
	return RespondError(w, r, http.StatusNotExtended, "not_extended", message, nil)
}

// NetworkAuthenticationRequired sends a 511 Network Authentication Required response.
func NetworkAuthenticationRequired(w http.ResponseWriter, r *http.Request, message string) error {
	if message == "" {
		message = "Network authentication required"
	}
	return RespondError(w, r, http.StatusNetworkAuthenticationRequired, "network_authentication_required", message, nil)
}
