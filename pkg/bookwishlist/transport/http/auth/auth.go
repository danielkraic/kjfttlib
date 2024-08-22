package auth

import "net/http"

type Config struct {
	Username string
	Password string
}

func (a *Config) Middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !a.isAuthorized(r) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		handler.ServeHTTP(w, r)
	})
}

func (a *Config) isAuthorized(r *http.Request) bool {
	if a.Username == "" && a.Password == "" {
		// auth is disabled
		return true
	}

	username, password, ok := r.BasicAuth()
	return ok && username == a.Username && password == a.Password
}
