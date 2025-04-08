package middleware

import (
	"log"
	"net"
	"net/http"

	"github.com/jozvahedi/loadbalancer/loadbalancer/internal/auth"
)

type Middleware interface {
	Wrap(http.HandlerFunc) http.HandlerFunc
}

// LoggingMiddleware implements the Middleware interface
type LoggingMiddleware struct{}

func (m LoggingMiddleware) Wrap(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	}
}

// IPWhitelistMiddleware implements the Middleware interface
type IPWhitelistMiddleware struct {
	Whitelist []string
}

func (m IPWhitelistMiddleware) Wrap(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "Invalid IP", http.StatusForbidden)
			return
		}

		allowed := false
		for _, allowedIP := range m.Whitelist {
			if ip == allowedIP {
				allowed = true
				break
			}
		}

		if !allowed {
			http.Error(w, "IP not allowed", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}
}

// BasicAuthMiddleware implements the Middleware interface
type BasicAuthMiddleware struct {
	AuthService auth.Authenticator
}

func (m BasicAuthMiddleware) Wrap(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok || !m.AuthService.Authenticate(username, password) {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

// Chain applies a list of middleware to a http.HandlerFunc
func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m.Wrap(f)
	}
	return f
}
