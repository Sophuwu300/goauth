package goauth

import (
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"

	"fmt"
	"os"
)

func generate(user string) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "soph.local",
		AccountName: user,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(key.URL())
}

func validate(key *otp.Key) {
	// Now Validate that the user's successfully added the passcode.
	fmt.Println("Validating TOTP...")
	passcode := "123456"
	valid := totp.Validate(passcode, key.Secret())
	if valid {
		println("Valid passcode!")
		os.Exit(0)
	} else {
		println("Invalid passcode!")
		os.Exit(1)
	}
}
