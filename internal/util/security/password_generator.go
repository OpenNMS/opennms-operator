package security

import (
	"math/rand"
	"strings"
)

const (
	LowerChars   = "abcdefghijklmnopqrstuvwxyz"
	UpperChars   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	SpecialChars = "!@#$%"
	Numbers      = "0123456789"
	Alphanumeric = LowerChars + UpperChars + Numbers
	AllChars     = Alphanumeric + SpecialChars
	Length       = 20
)

//GeneratePassword - generates a random password of 20 length
func GeneratePassword(specialChars bool) string {
	var password strings.Builder

	//Set special character
	if specialChars {
		for i := 0; i < 2; i++ {
			random := rand.Intn(len(SpecialChars))
			password.WriteString(string(SpecialChars[random]))
		}
	}

	//Set numeric
	for i := 0; i < 2; i++ {
		random := rand.Intn(len(Numbers))
		password.WriteString(string(Numbers[random]))
	}

	//Set uppercase
	for i := 0; i < 2; i++ {
		random := rand.Intn(len(UpperChars))
		password.WriteString(string(UpperChars[random]))
	}

	remainingLength := Length - password.Len()
	charSet := AllChars
	if !specialChars {
		charSet = Alphanumeric
	}
	for i := 0; i < remainingLength; i++ {
		random := rand.Intn(len(charSet))
		password.WriteString(string(charSet[random]))
	}
	inRune := []rune(password.String())
	rand.Shuffle(len(inRune), func(i, j int) {
		inRune[i], inRune[j] = inRune[j], inRune[i]
	})
	return string(inRune)
}
