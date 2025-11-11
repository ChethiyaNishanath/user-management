package middleware

import (
	"context"
	"net/http"
	"strconv"
)

type contextKey string

const (
	PageKey  contextKey = "page"
	LimitKey contextKey = "limit"
)

func Paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := 1
		limit := 10

		q := r.URL.Query()

		if p := q.Get("page"); p != "" {
			if val, err := strconv.Atoi(p); err == nil && val > 0 {
				page = val
			}
		}

		if l := q.Get("limit"); l != "" {
			if val, err := strconv.Atoi(l); err == nil && val > 0 && val <= 100 {
				limit = val
			}
		}

		ctx := context.WithValue(r.Context(), PageKey, page)
		ctx = context.WithValue(ctx, LimitKey, limit)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
