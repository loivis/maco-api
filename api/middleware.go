package api

import (
	"net/http"
	"strings"
)

func allowMethods(methods []string, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var allowed bool

		for _, method := range methods {
			if r.Method == method {
				allowed = true
			}
		}

		if !allowed {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		handler.ServeHTTP(w, r)
	})
}

func validatePath(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paths := strings.Split(r.URL.Path, "/")

		if len(paths) > 3 {
			http.Error(w, "unsupported resource", http.StatusBadRequest)
			return
		}

		handler.ServeHTTP(w, r)
	})
}
