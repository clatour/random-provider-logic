package main

import (
	"crypto/rand"
	"log"
	"math/big"
	"sort"
	"strings"

	"github.com/hashicorp/errwrap"
)

var DefaultOptions = &options{
	length:          32,
	special:         true,
	upper:           true,
	lower:           true,
	number:          true,
	minNumeric:      0,
	minUpper:        0,
	minLower:        0,
	minSpecial:      0,
	overrideSpecial: "",
}

func main() {
	o := DefaultOptions

	o.length = 32
	o.special = true

	for n := 0; n < 1_000_000; n++ {
		s, err := generatePassword(o)
		if err != nil {
			log.Fatalf(err.Error())
		}

		if strings.Contains(s, "'") {
			log.Printf("string contains \"'\", %s\n", s)
		}

		if n%10000 == 0 {
			log.Println(s)
		}
	}
}

type options struct {
	length          int
	upper           bool
	minUpper        int
	lower           bool
	minLower        int
	number          bool
	minNumeric      int
	special         bool
	minSpecial      int
	overrideSpecial string
}

func generateRandomBytes(charSet *string, length int) ([]byte, error) {
	bytes := make([]byte, length)
	setLen := big.NewInt(int64(len(*charSet)))
	for i := range bytes {
		idx, err := rand.Int(rand.Reader, setLen)
		if err != nil {
			return nil, err
		}
		bytes[i] = (*charSet)[idx.Int64()]
	}
	return bytes, nil
}

func generatePassword(o *options) (string, error) {
	const numChars = "0123456789"
	const lowerChars = "abcdefghijklmnopqrstuvwxyz"
	const upperChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var specialChars = "!@#$%&*()-_=+[]{}<>:?"

	length := o.length
	upper := o.upper
	minUpper := o.minUpper
	lower := o.lower
	minLower := o.minLower
	number := o.number
	minNumeric := o.minNumeric
	special := o.special
	minSpecial := o.minSpecial
	overrideSpecial := o.overrideSpecial

	if overrideSpecial != "" {
		specialChars = overrideSpecial
	}

	var chars = string("")
	if upper {
		chars += upperChars
	}
	if lower {
		chars += lowerChars
	}
	if number {
		chars += numChars
	}
	if special {
		chars += specialChars
	}

	minMapping := map[string]int{
		numChars:     minNumeric,
		lowerChars:   minLower,
		upperChars:   minUpper,
		specialChars: minSpecial,
	}

	var result = make([]byte, 0, length)
	for k, v := range minMapping {
		s, err := generateRandomBytes(&k, v)
		if err != nil {
			return "", errwrap.Wrapf("error generating random bytes: {{err}}", err)
		}
		result = append(result, s...)
	}
	s, err := generateRandomBytes(&chars, length-len(result))
	if err != nil {
		return "", errwrap.Wrapf("error generating random bytes: {{err}}", err)
	}
	result = append(result, s...)
	order := make([]byte, len(result))
	if _, err := rand.Read(order); err != nil {
		return "", errwrap.Wrapf("error generating random bytes: {{err}}", err)
	}
	sort.Slice(result, func(i, j int) bool {
		return order[i] < order[j]
	})

	return string(result), nil
}
