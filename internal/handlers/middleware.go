package handlers

import (
	"fmt"
	"frappuccino/internal/utils"
	"net/http"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func ChainMiddleware(h http.HandlerFunc, m ...Middleware) http.HandlerFunc {
	if len(m) < 1 {
		return h
	}

	wrapped := h
	for i := len(m) - 1; i >= 0; i-- {
		wrapped = m[i](wrapped)
	}

	return wrapped
}

func logRequest(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
		)

		logger := utils.GetLogger()
		logger.Info("received request", "ip", ip, "proto", proto, "method", method, "uri", uri)

		next.ServeHTTP(w, r)
	})
}

func recoverPanic(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				var (
					method = r.Method
					uri    = r.URL.RequestURI()
				)
				logger := utils.GetLogger()
				logger.Error(fmt.Errorf("%s", err).Error(), "method", method, "uri", uri)

				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
