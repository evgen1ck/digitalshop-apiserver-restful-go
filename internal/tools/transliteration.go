package tools

import (
	"strings"
)

var (
	latinToCyrillicMap = map[string]string{
		"a": "а", "b": "б", "v": "в", "g": "г", "d": "д", "e": "е",
		"yo": "ё", "zh": "ж", "z": "з", "i": "и", "j": "й", "k": "к",
		"l": "л", "m": "м", "n": "н", "o": "о", "p": "п", "r": "р",
		"s": "с", "t": "т", "u": "у", "f": "ф", "h": "х", "c": "ц",
		"ch": "ч", "sh": "ш", "sch": "щ", "y": "ы", "yu": "ю", "ya": "я",
	}
	cyrillicToLatinMap = map[string]string{
		"а": "a", "б": "b", "в": "v", "г": "g", "д": "d", "е": "e",
		"ё": "yo", "ж": "zh", "з": "z", "и": "i", "й": "j", "к": "k",
		"л": "l", "м": "m", "н": "n", "о": "o", "п": "p", "р": "r",
		"с": "s", "т": "t", "у": "u", "ф": "f", "х": "h", "ц": "c",
		"ч": "ch", "ш": "sh", "щ": "sch", "ы": "y", "э": "e", "ю": "yu", "я": "ya",
	}

	engLayout = []rune("qwertyuiop[]asdfghjkl;'zxcvbnm,./`")
	rusLayout = []rune("йцукенгшщзхъфывапролджэячсмитьбю.ё")
)

func Transliterate(text string) string {
	if isLatin(text) {
		return latinToCyrillic(text)
	}
	return cyrillicToLatin(text)
}

func isLatin(text string) bool {
	for _, r := range text {
		if r >= 'а' && r <= 'я' {
			return false
		}
	}
	return true
}

func latinToCyrillic(text string) string {
	var result strings.Builder
	for _, r := range text {
		mapped := latinToCyrillicMap[string(r)]
		if mapped != "" {
			result.WriteString(mapped)
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func cyrillicToLatin(text string) string {
	var result strings.Builder
	for _, r := range text {
		mapped := cyrillicToLatinMap[string(r)]
		if mapped != "" {
			result.WriteString(mapped)
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func EngToRus(input string) string {
	var result strings.Builder

	for _, char := range input {
		engIndex := strings.IndexRune(string(engLayout), char)
		if engIndex != -1 {
			result.WriteRune(rusLayout[engIndex])
		} else {
			result.WriteRune(char)
		}
	}

	return result.String()
}

func RusToEng(input string) string {
	var result strings.Builder

	for _, char := range input {
		rusIndex := strings.IndexRune(string(rusLayout), char)
		if rusIndex != -1 {
			result.WriteRune(engLayout[rusIndex])
		} else {
			result.WriteRune(char)
		}
	}

	return result.String()
}
