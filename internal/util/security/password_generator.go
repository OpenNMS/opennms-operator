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
	AllChars     = LowerChars + UpperChars + SpecialChars + Numbers
)

func GeneratePassword() string {
	var password strings.Builder

	//Set special character
	for i := 0; i < 2; i++ {
		random := rand.Intn(len(SpecialChars))
		password.WriteString(string(SpecialChars[random]))
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

	remainingLength := 14
	for i := 0; i < remainingLength; i++ {
		random := rand.Intn(len(AllChars))
		password.WriteString(string(AllChars[random]))
	}
	inRune := []rune(password.String())
	rand.Shuffle(len(inRune), func(i, j int) {
		inRune[i], inRune[j] = inRune[j], inRune[i]
	})
	return string(inRune)
}