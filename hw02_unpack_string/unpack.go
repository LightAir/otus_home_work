package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

const MaxSequenceRepeatingChars = 2

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	var result strings.Builder

	inputLength := len(input)
	maxCharacterIndex := inputLength - 1

	isLetter := false
	repeatCharCounter := 0

	for index, currentRune := range input {
		if index == 0 && unicode.IsDigit(currentRune) {
			return "", ErrInvalidString
		}

		nextCharacterIndex := index + 1

		if string(currentRune) == "\\" && !isLetter {
			if validateCharAfterSlash(nextCharacterIndex, maxCharacterIndex, input) != nil {
				return "", ErrInvalidString
			}
			isLetter = true
			continue
		}

		if unicode.IsDigit(currentRune) && !isLetter {
			repeatCharCounter++
			if repeatCharCounter >= MaxSequenceRepeatingChars {
				return "", ErrInvalidString
			}
		} else {
			repeatCharCounter = 0
		}

		if unicode.IsLetter(currentRune) || isLetter {
			if nextCharacterIndex > maxCharacterIndex {
				result.WriteRune(currentRune)
				break
			}

			nextRune := rune(input[nextCharacterIndex])

			if unicode.IsDigit(nextRune) {
				digit, _ := strconv.Atoi(string(nextRune))
				result.WriteString(strings.Repeat(string(currentRune), digit))
			} else {
				result.WriteRune(currentRune)
			}

			isLetter = false
		}
	}

	return result.String(), nil
}

func validateCharAfterSlash(nextCharacterIndex int, maxCharacterIndex int, input string) error {
	if nextCharacterIndex <= maxCharacterIndex {
		nextRune := rune(input[nextCharacterIndex])
		if string(nextRune) == "\\" || unicode.IsDigit(nextRune) {
			return nil
		}
	}

	return ErrInvalidString
}
