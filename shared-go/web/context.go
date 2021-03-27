package web

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "rest context value " + k.name
}

func NewContextKey(key string) *contextKey {
	return &contextKey{key}
}

func RequestCtxValue(key string) (func(next http.Handler) http.Handler, func(r *http.Request) string) {
	ctxKey := contextKey{key}

	set := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			ctxValue := chi.URLParam(request, key)

			next.ServeHTTP(writer, request.WithContext(context.WithValue(request.Context(), ctxKey, ctxValue)))
		})
	}
	get := func(r *http.Request) string {
		val := r.Context().Value(ctxKey)

		if val == nil {
			return ""
		}

		return val.(string)
	}

	return set, get
}
