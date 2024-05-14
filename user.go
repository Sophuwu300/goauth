package goauth

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"github.com/pquerna/otp/totp"
)

// User struct for storing user data
type User struct {
	Name string
	Hash []byte
	Salt []byte
	OtpS string
}

func salt() []byte {
	var s = make([]byte, 32)
	_, err := rand.Read(s)
	if err != nil {
		panic(err)
	}
	return s
}
func hash(b ...[]byte) []byte {
	h := sha256.New()
	for _, v := range b {
		h.Write(v)
	}
	return h.Sum(nil)
}

func (u *User) SetPass(pass string) {
	u.Salt = salt()
	u.Hash = hash([]byte(pass), u.Salt)
}

func (u *User) CheckPass(pass string) bool {
	return subtle.ConstantTimeCompare(u.Hash, hash([]byte(pass), u.Salt)) == 1
}

// NewUser creates a new user
func NewUser(name, pass string) *User {
	var u User
	key, err := totp.Generate(totp.GenerateOpts{AccountName: name, Issuer: "soph.local"})
	if err != nil {
		panic(err)
	}
	key.Secret()
	u.Name = name
	u.SetPass(pass)

}

// CheckPassword checks if the password is correct
func CheckPassword(name, pass string) bool {
	return false
}

// subtle.ConstantTimeCompare([]byte(expectedPass), []byte(pass)) == 1
