package random

import (
	"crypto/rand"
	"encoding/base64"
	"math/big"
	"regexp"
	"strings"
)

// String returns a random string consisting only of letters of the length n
func String(n uint) (string, error) {
	result := make([]string, n)
	bytes := make([]byte, n*100)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	chars := strings.Split(base64.URLEncoding.EncodeToString(bytes), "")

	var i uint
	for i = 0; i < n; i++ {
		// we have to use big Int since crypto/rand expects it
		// however max is basically just the maximum amount available
		// characters
		max := big.NewInt(int64(len(chars) - 1))
		rnd, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}

		c := chars[rnd.Int64()]
		isChar, err := regexp.MatchString("[A-Za-z]", c)
		if err != nil {
			return "", err
		}

		if !isChar {
			i--
			continue
		}

		result[i] = c
	}

	return strings.Join(result, ""), nil
}
