package htsesh

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base32"
	"encoding/base64"
	"net/http"
	"sync/atomic"
	"time"
)

type session struct {
	V []byte
	T time.Time
}

func (s *session) Hash() string {
	h := sha512.New()
	h.Write(s.V)
	h.Write([]byte(s.T.String()))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

func (s *session) Validate(cookie *http.Cookie) bool {
	if s == nil {
		return false
	}
	if time.Now().After(s.T) || s.T.After(time.Now().Add(time.Hour)) {
		s.V = nil
		s.T = time.Time{}
		sessionID.Store(nil)
		return false
	}
	if cookie == nil || cookie.Value != s.Hash() {
		return false
	}
	return true
}

func sessionOK(r *http.Request) bool {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return false
	}
	return sessionID.Load().Validate(cookie)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/logout" {
		ji := base32.StdEncoding.EncodeToString(sessionID.Load().V[:5])
		if r.URL.Query().Get("logoutid") == ji {
			sessionID.Store(nil)
			http.SetCookie(w, &http.Cookie{
				Name:     "session_id",
				Value:    "",
				HttpOnly: true,
				Secure:   true,
				Expires:  time.Now(),
				SameSite: http.SameSiteStrictMode,
			})
			http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		} else {
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("<html><body><a href='?logout&logoutid=" + ji + "'>Logout</a></body></html>"))
		}
		return
	}
}

func Authenticate(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if sessionOK(r) {
			next.ServeHTTP(w, r)
		} else if loginHandler(w, r) {
			http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		} else {
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(FORM))
		}
	})
}

func loginHandler(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodPost || r.ParseForm() != nil {
		return false
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	otp := r.FormValue("otp")
	if username == "" || password == "" || otp == "" {
		return false
	}
	if !verify(username, password, otp) {
		return false
	}
	s := session{
		V: make([]byte, 32),
		T: time.Now().Add(time.Hour),
	}
	n, err := rand.Read(s.V)
	if err != nil || 32 != len(s.V) || n != 32 {
		return false
	}
	sessionID.Store(&s)
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    s.Hash(),
		HttpOnly: true,
		Secure:   true,
		Expires:  s.T,
		SameSite: http.SameSiteStrictMode,
	})
	return true
}

var (
	verify = func(username, password, otp string) bool {
		return false
	}
	sessionID atomic.Pointer[session]
)

func SetVerifyFunc(f func(username, password, otp string) bool) {
	verify = f
}

var FORM = `<html><body><form method="post" style="width: min-content;height: min-content;transform: translate(-50%,-50%);top: 50%;position: absolute;left: 50%;">
<input type="text" placeholder="username" autocomplete="off" name="username">
<input type="password" autocomplete="off" name="password" placeholder="password">
<input type="text" autocomplete="off" placeholder="otp" name="otp">
<input type="submit" value="submit">
</form></body></html>`
