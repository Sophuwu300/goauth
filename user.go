package goauth

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"github.com/pquerna/otp/totp"
	"os"
	"time"
)

// User struct for storing user data
type User struct {
	Name string
	Hash []byte
	Salt []byte
	OtpS string
}

func salt(n int) []byte {
	var s = make([]byte, n)
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
	u.Salt = salt(16)
	u.Hash = hash([]byte(pass), u.Salt)
}

func (u *User) CheckPass(pass string) bool {
	return subtle.ConstantTimeCompare(u.Hash, hash([]byte(pass), u.Salt)) == 1
}

// NewUser creates a new user
func NewUser(name, pass string) (User, error) {
	var u User
	u.Name = name
	u.SetPass(pass)
	otp, e := totp.Generate(totp.GenerateOpts{Issuer: "soph.local", AccountName: u.Name})
	if e != nil {
		return u, e
	}
	u.OtpS = otp.Secret()
	q, e := GenQR(otp.URL(),
		"One Time Password:",
		"  User  : "+u.Name,
		"  Issuer: "+"local",
		"  Secret: "+otp.Secret(),
		"  Period: "+fmt.Sprint(time.Duration(otp.Period()*10e8).String()),
		"  Digits: "+fmt.Sprintf("%d", otp.Digits()),
		"  Algo  : "+fmt.Sprintf("%s", otp.Algorithm()),
		"Scan the QR code with your authenticator app.",
	)
	if e != nil {
		return u, e
	}
	fmt.Println(q.String())
	os.WriteFile(u.Name+".png", q.Png(), 0644)
	return u, nil
}

// CheckPassword checks if the password is correct
func CheckPassword(name, pass string) bool {
	return false
}

// subtle.ConstantTimeCompare([]byte(expectedPass), []byte(pass)) == 1
