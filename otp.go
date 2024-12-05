package goauth

import (
	"fmt"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"os"
)

// ServiceHost is the host of the service
var ServiceHost string

// init here gets the hostname used for issuing keys
func init() {
	s, e := os.Hostname()
	if e != nil {
		ServiceHost = "localhost"
	} else {
		ServiceHost = s
	}
}

// getSecurityLevel returns the number of digits and seconds allowed for a security level
func getSecurityLevel(i int) (otp.Digits, uint, uint, error) {
	if i > 4 || i < 0 {
		return 0, 0, 0, Error("invalid security level")
	}
	u := (securityLevel & (255 << (8 * uint64(i)))) >> (8 * uint64(i))
	x := uint(u & 0b11)
	return otp.Digits(u >> 2), (x*x*25)/10 + (225*x)/10 + 5, x, nil
}

// OTP holds the information for the otp
type OTP struct {
	Secret string
	Vals   totp.ValidateOpts
}

// Users is a map of users to their otp keys for validation
var Users map[string]OTP

// NewUser creates a new user
func NewUser(name string, securityLevel ...int) (QRcode, error) {
	if _, ok := Users[name]; ok {
		return QRcode{}, fmt.Errorf("user %s already exists", name)
	}

	if len(securityLevel) == 0 {
		securityLevel = []int{SecurityLevelDefault}
	}

	digits, period, skew, e := getSecurityLevel(securityLevel[0])
	if e != nil {
		return QRcode{}, e
	}

	userOtp, e := totp.Generate(totp.GenerateOpts{
		Issuer:      ServiceHost,
		AccountName: name,
		Digits:      digits,
		Period:      period,
	})
	if e != nil {
		return QRcode{}, e
	}

	Users[name] = OTP{
		Secret: userOtp.Secret(),
		Vals: totp.ValidateOpts{
			Period: period,
			Digits: digits,
			Skew:   skew,
		},
	}

	q, e := GenQR(userOtp.URL(),
		"OTP for "+name+" at "+ServiceHost,
		"Copy the secret below into your OTP app",
		"Secret: "+userOtp.Secret(),
		"Or scan the QR code below",
	)
	return q, e
}

// DelUser deletes a user
func DelUser(name string) error {
	if _, ok := Users[name]; !ok {
		return fmt.Errorf("user %s does not exist", name)
	}
	delete(Users, name)
	return nil
}

// ValidateOtp validates the otp
func ValidateOtp(user, pass string) bool {
	if _, ok := Users[user]; !ok {
		return false
	}

	return totp.Validate(pass, Users[user].Secret)
}
