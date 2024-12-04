package goauth

// SecurityLevelAccessible - Lower security level but easier to use
//
//	This will give 60 seconds for the user to write 6 digits. Its aim is to help
//	people with difficulties using a computer by allowing them more time to type.
//	The actual security implications are minimal.
const SecurityLevelAccessible = 0

// SecurityLevelDefault - Default security level used everywhere
//
//	This will give 30 seconds for the user to write 6 digits. This is the same as
//	the default for Google Authenticator.
const SecurityLevelDefault = 1

// SecurityLevelCozyHigh - Higher security but more time for the user
//
//	This will give 60 seconds for the user to write 10 digits. This offers a
//	slight security improvement. While remaining easier to use than other levels.
const SecurityLevelCozyHigh = 2

// SecurityLevelHigh - Same as default but with double the digits
//
//	With this you will get 30 seconds to write 12 digits. This is significantly
//	harder to brute force than the default level. This is the highest level
//	recommended for general use.
const SecurityLevelHigh = 3

// SecurityLevelRealTime - Extremely long password with very little time
//
//	Gives 5 seconds to write 15 digits. Intended for computers and application
//	interface use for automated task. May be useful for authentication of
//	unattended systems with less risk of being intercepted.
const SecurityLevelRealTime = 4

// securityLevel is a bitfield of the security levels
// each level is 8 bits long and contains another bitfield for the length
// of the password and the time allowed to enter it, allowing customisation
// ease of use and security
const securityLevel = uint64(uint64((15<<2)|0b00)<<(8*4) | uint64((12<<2)|0b01)<<(8*3) | uint64((10<<2)|0b10)<<(8*2) | uint64((6<<2)|0b01)<<8 | uint64((6<<2)|0b10))

// getSecurityLevel returns the number of digits and seconds allowed for a security level
func getSecurityLevel(i int) (int, int, error) {

	if i > 4 || i < 0 {
		return 0, 0, Error("invalid security level")
	}
	u := (securityLevel & (255 << (8 * uint64(i)))) >> (8 * uint64(i))
	x := int(u & 0b11)
	return int(u >> 2), (x*x*25)/10 + (225*x)/10 + 5, nil
}

func Error(text string) error {
	return &errorString{text}
}

// errorString is a trivial implementation of error.
type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}
