package helper

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/unicode/norm"
)

var (
	reNonAlphaNumeric = regexp.MustCompile(`[^a-zA-Z0-9]+`)
	reAlpha           = regexp.MustCompile(`^[a-zA-Z]+$`)
	reNumeric         = regexp.MustCompile(`^[0-9]+$`)
)

func IsEmpty(s string) bool                 { return s == "" }
func IsNotEmpty(s string) bool              { return !IsEmpty(s) }
func ConcatStrings(values ...string) string { return strings.Join(values, "") }

func ConcatStringsWithSeparator(separator string, values ...string) string {
	var buffer bytes.Buffer
	for index, value := range values {
		buffer.WriteString(value)
		if index != len(values)-1 {
			buffer.WriteString(separator)
		}
	}
	return buffer.String()
}

func StringToUint64(value string) (uint64, error) {
	if IsEmpty(value) {
		return 0, fmt.Errorf("[pkg.helperx.StringToUint64] value is empty")
	}

	num, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("[pkg.helperx.StringToUint64] failed to convert string to uint64: %w", err)
	}

	return num, nil
}

func StringToInt64(value string) (int64, error) {
	if IsEmpty(value) {
		return 0, fmt.Errorf("[pkg.helperx.StringToInt64] value is empty")
	}
	num, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("[pkg.helperx.StringToInt64] failed to convert string to int64: %w", err)
	}
	return num, nil
}

func StringToInt(value string) (int, error) {
	value = strings.TrimSpace(value)
	if IsEmpty(value) {
		return 0, fmt.Errorf("[pkg.helperx.StringToInt] value is empty")
	}

	num, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("[pkg.helperx.StringToInt] failed to convert string to int: %w", err)
	}
	return num, nil
}

func BytesToString(b []byte) string {
	return string(b)
}

func StringToBytes(s string) []byte {
	return []byte(s)
}

func ReverseString(s string) string {
	runes := []rune(s)
	result := make([]rune, len(runes))
	for index, value := range runes {
		result[len(runes)-index-1] = value
	}
	return string(result)
}

func Slugify(s string) string {
	s = removeAccents(s)
	s = strings.ToLower(s)
	s = reNonAlphaNumeric.ReplaceAllString(s, "-")
	return s
}

func IsAlpha(s string) bool {
	return reAlpha.MatchString(s)
}

func IsNumeric(s string) bool {
	return reNumeric.MatchString(s)
}

func CamelCase(s string) string {
	return cases.Title(language.English).String(s)
}

func SnakeCase(s string) string {
	return strings.ToLower(s)
}

func removeAccents(s string) string {
	t := norm.NFD.String(s)
	sb := strings.Builder{}
	for _, r := range t {
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		sb.WriteRune(r)
	}
	return sb.String()
}
