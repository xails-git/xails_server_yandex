package service

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/pkg/morse"
)

func AutoConvert(input string) (string, error) {
	if input == "" {
		return "", fmt.Errorf("пустой ввод")
	}

	isMorse := regexp.MustCompile(`^[\s.\-/]+$`).MatchString(input)

	if isMorse {
		text := morse.ToText(input)
		if text == "" {
			return "", fmt.Errorf("некорректный код Морзе")
		}
		return text, nil
	}

	upperInput := strings.ToUpper(input)
	morseCode := morse.ToMorse(upperInput)
	if morseCode == "" {
		return "", fmt.Errorf("некорректный текст")
	}
	return morseCode, nil
}
