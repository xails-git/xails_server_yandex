// Пакет service
package service

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/pkg/morse"
)

func AutoConvert(input string) (string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", fmt.Errorf("empty input")
	}

	// Проверка на Morse code (допускаем любые комбинации . - / и пробелов)
	if isMorse(input) {
		text := morse.ToText(input)
		if text == "" {
			return "", fmt.Errorf("invalid morse code")
		}
		return text, nil
	}

	// Конвертация текста в верхний регистр
	morseCode := morse.ToMorse(strings.ToUpper(input))
	if morseCode == "" {
		return "", fmt.Errorf("invalid text")
	}

	return morseCode, nil
}

func isMorse(s string) bool {
	return regexp.MustCompile(`^[\s.\-/]+$`).MatchString(s)
}
