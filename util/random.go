package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	// Seed the random number generator with the current time (ensures different sequences each time)
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates a random integer between a specified minimum and maximum value (inclusive)
func RandomInt(min, max int64) int64 {
	// Use rand.Int63n to generate a random number within the range [0, max-min]
	// Add min to the result to get the final value within the range [min, max]
	return min + rand.Int63n(max-min+1)
}

// RandomString generates a random string of a specified length
func RandomString(n int) string {
	// Create a StringBuilder object to efficiently build the string
	var sb strings.Builder
	k := len(alphabet) // Get the length of the alphabet string

	// Loop n times to generate random characters and append them to the StringBuilder
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)] // Get a random character from the alphabet
		sb.WriteByte(c)             // Append the character to the StringBuilder
	}

	// Return the final string built from the random characters
	return sb.String()
}

// RandomOwner generates a random owner name (assuming a simple 6-character string)
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney generates a random amount of money between 0 and 1000 (inclusive)
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// RandomCurrency generates a random currency from a pre-defined list
func RandomCurrency() string {
	currencies := []string{"EUR", "USD", "CAD"} // Define a list of possible currencies
	n := len(currencies)                        // Get the length of the currency list

	// Return a random element from the currency list
	return currencies[rand.Intn(n)]
}

// RandomEmail generates a random email address with a 6-character username and gmail.com domain
func RandomEmail() string {
	// Use fmt.Sprintf to format a string with the username and domain
	return fmt.Sprintf("%s@gmail.com", RandomString(6))
}
