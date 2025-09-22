package main

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

func Unpack(s string) (string, error) {
	var (
		result   strings.Builder
		lastRune rune
	)

	escaped := false
	for _, r := range s {
		switch {
		case escaped:

			result.WriteRune(r)
			lastRune = r
			escaped = false

		case r == '\\':

			escaped = true

		case unicode.IsDigit(r):

			if result.Len() == 0 {
				return "", errors.New("incorrect string")
			}

			if r == '0' {
				return "", errors.New("incorrect digit '0'")
			}

			count := int(r - '0')
			result.WriteString(strings.Repeat(string(lastRune), count-1))

		default:
			result.WriteRune(r)
			lastRune = r
		}
	}

	if escaped {

		return "", errors.New("incomplete escape sequence")

	}

	return result.String(), nil
}

func main() {
	tests := []string{"a4bc2d5e", "abcd", "45", "", "qwe\\4\\5", "qwe\\45"}

	for _, test := range tests {
		result, err := Unpack(test)
		if err != nil {
			fmt.Printf("Вход: %q -> Ошибка: %v\n", test, err)
		} else {
			fmt.Printf("Вход: %q -> Выход: %q\n", test, result)
		}
	}
}
