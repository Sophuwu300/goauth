package http

import (
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pquerna/otp/totp"
	"net/http"
	"os"
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

func newSession(w http.ResponseWriter) error {
	s := session{
		V: make([]byte, 32),
		T: time.Now().Add(time.Hour),
	}
	n, err := rand.Read(s.V)
	if err != nil || 32 != len(s.V) || n != 32 {
		return fmt.Errorf("failed to generate session: %w", err)
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    s.Hash(),
		HttpOnly: true,
		Secure:   true,
		Expires:  s.T,
		SameSite: http.SameSiteStrictMode,
	})
	sessionID.Store(&s)
	return nil
}

func validateSession(r *http.Request) bool {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return false
	}
	s := sessionID.Load()
	if s == nil || cookie.Value != s.Hash() {
		return false
	}
	if time.Now().After((*s).T) || (*s).T.After(time.Now().Add(time.Hour)) {
		sessionID.Store(nil)
		return false
	}
	return true
}

var sessionID atomic.Pointer[session]

func Authenticate(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if validateSession(r) || loginHandler(w, r) {
			next.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(html))
	})
}

type User struct {
	Username string // `storm:"id,index,unique"`
	OTP      string
	Hash     []byte
	Salt     []byte
}

func randomBytes(b *[]byte, n int) error {
	*b = make([]byte, n)
	nn, err := rand.Read(*b)
	if err != nil || nn != n {
		return err
	}
	return nil
}

func (u *User) HashPassword(password string) []byte {
	h := sha512.New()
	h.Write([]byte(password))
	h.Write(u.Salt)
	return h.Sum(nil)
}

func (u *User) CheckPassword(password string) bool {
	return subtle.ConstantTimeCompare(u.Hash, u.HashPassword(password)) == 1
}

func Sh1(username string) string {
	h := sha1.New()
	h.Write([]byte(username))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

func LoadUser(username string) (*User, error) {
	var u User
	b, err := os.ReadFile("./.passwd/" + Sh1(username))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func NewUser(username, password string) error {
	var u User
	u.Username = username

	if err := randomBytes(&u.Salt, 32); err != nil {
		return err
	}
	u.Hash = u.HashPassword(password)

	host, err := os.Hostname()
	if err != nil {
		host = "unknown.host"
	}
	k, err := totp.Generate(totp.GenerateOpts{
		Issuer:      host,
		AccountName: username,
	})
	if err != nil {
		return err
	}
	u.OTP = k.Secret()
	var b []byte
	b, err = json.Marshal(&u)
	if err != nil {
		return err
	}
	fmt.Println("user", username, "created with secret", u.OTP)
	return os.WriteFile("./.passwd/"+Sh1(username), b, 0600)

}

func (u *User) ValidateOtp(otp string) bool {
	return totp.Validate(otp, u.OTP)
}

func loginHandler(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodPost {
		return false
	}

	err := r.ParseForm()
	if err != nil {
		return false
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	otp := r.FormValue("otp")

	if username == "" || password == "" || otp == "" {
		return false
	}

	u, err := LoadUser(username)
	if err != nil {
		return false
	}
	if !(u.CheckPassword(password) && u.ValidateOtp(otp)) {
		return false
	}

	return newSession(w) == nil
}

func init() {
	if err := os.MkdirAll("./.passwd", 0700); err != nil {
		panic(err)
	}
	dr, err := os.ReadDir("./.passwd")
	if err != nil {
		panic(err)
	}
	if len(dr) == 0 {
		if err = NewUser("user", "password"); err != nil {
			panic(err)
		}
	}
}

func main() {

	http.ListenAndServe(":8000", Authenticate(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	}))
}

var html = `<html><body><form method="post" style="width: min-content;height: min-content;transform: translate(-50%,-50%);top: 50%;position: absolute;left: 50%;"><input type="text" placeholder="username" autocomplete="off" name="username">
<input type="password" autocomplete="off" name="password" placeholder="password">
<input type="text" autocomplete="off" placeholder="otp" name="otp">
  <input type="submit" value="submit">
</form></body></html>`
