package handlers

import (
	"net/http"
	"strconv"
)

const (
	defaultPageLimit = 20
	maxPageLimit     = 100
)

func parsePagination(r *http.Request) (limit, offset int) {
	limit = defaultPageLimit
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= maxPageLimit {
			limit = n
		}
	}
	if v := r.URL.Query().Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			offset = n
		}
	}
	return limit, offset
}

func trimPage[T any](rows []T, limit int) (items []T, hasMore bool) {
	if len(rows) > limit {
		return rows[:limit], true
	}
	if rows == nil {
		return []T{}, false
	}
	return rows, false
}
