package service

import (
	"fmt"
	"strings"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/pkg/morse"
)

func AutoConvert(input string) (string, error) {
	allowed := ".-/ "
	isMorse := true

	if input == "" {
		return "", fmt.Errorf("пустой ввод")
	}

	for _, char := range input {
		if !strings.ContainsRune(allowed, char) {
			isMorse = false
			break
		}
	}

	if isMorse {
		text := morse.ToText(input)
		if text == "" {
			return "", fmt.Errorf("ошибка конвертации в текст")
		}
		return text, nil
	}

	morseCode := morse.ToMorse(input)
	if morseCode == "" {
		return "", fmt.Errorf("ошибка конвертации в морзе")
	}
	return morseCode, nil

}
