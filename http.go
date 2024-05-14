package goauth

import (
	"crypto/subtle"
	"net/http"
)

func HttpAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		inputUser, inputPass, authOK := r.BasicAuth()
		CheckPassword(inputUser, inputPass)
		if !authOK || !lookupOK ||  {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized.", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
