package auth

import "net/http"

type SingleUserAuth struct {
	password string
}

func NewSingleUserAuth(password string) *SingleUserAuth {
	return &SingleUserAuth{password: password}
}

func (a *SingleUserAuth) Wrap(next http.Handler) http.Handler {
	if a.password == "" {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, pass, ok := r.BasicAuth()
		if !ok || pass != a.password {
			w.Header().Set("WWW-Authenticate", `Basic realm="HabitClaw"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
