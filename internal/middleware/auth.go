package middleware

import (
	"net/http"

	"homework/internal/config"
)

func AuthMiddleware(auth config.Auth) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			user, pass, ok := request.BasicAuth()
			if !ok || user != auth.User || pass != auth.Password {
				writer.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
				writer.WriteHeader(http.StatusUnauthorized)
				_, err := writer.Write([]byte("Unauthorized\n"))
				if err != nil {
					writer.WriteHeader(http.StatusInternalServerError)
					return
				}
				return
			}
			handler.ServeHTTP(writer, request)
		})
	}
}
