package response

import (
	"net/http"
	"strconv"
)

type Pagination struct {
	Limit  int
	Offset int
}

// ParsePagination parses limit and offset from query parameters with sane defaults and bounds.
// - limit default: 20, min: 1, max: 100
// - offset default: 0, min: 0
func ParsePagination(r *http.Request) Pagination {
	const (
		defaultLimit = 20
		maxLimit     = 100
	)

	limit := defaultLimit
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			if n < 1 {
				limit = 1
			} else if n > maxLimit {
				limit = maxLimit
			} else {
				limit = n
			}
		}
	}

	offset := 0
	if v := r.URL.Query().Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			if n < 0 {
				offset = 0
			} else {
				offset = n
			}
		}
	}

	return Pagination{Limit: limit, Offset: offset}
}
