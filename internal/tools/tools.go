package tools

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"math/big"
	"net"
	"net/url"
	"strconv"
	"strings"
	"unicode"
)

// GenerateSixDigitNumber generates a six-digit random number.
// It does this by generating 3 random bytes and converting them to a big integer.
// It then limits the big integer to a six-digit number and ensures that the first digit is not zero.
// Finally, it returns the string representation of the generated number.
func GenerateSixDigitNumber() (string, error) {
	// Generate 3 random bytes
	b := make([]byte, 3)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	// Convert bytes to a number
	n := new(big.Int).SetBytes(b)

	// Limit the number to 6 digits and ensure that the first digit is not zero
	n.Mod(n, big.NewInt(900000))
	n.Add(n, big.NewInt(100000))

	return n.String(), nil
}

// GenerateRandomString generates a random string of the specified length.
// It generates a slice of bytes with the specified length and fills it with random bytes using the crypto/rand package.
// It returns the byte slice converted to a string.
func GenerateRandomString(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// GenerateURLToken generates a URL-safe token of the specified length.
func GenerateURLToken(length int) (string, error) {
	nonce := make([]byte, length)
	_, err := io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", err
	}

	secret, err := GenerateRandomClassicString(64)
	if err != nil {
		return "", err
	}

	h := hmac.New(sha256.New, []byte(secret))
	h.Write(nonce)

	hash := h.Sum(nil)
	token := base64.RawURLEncoding.EncodeToString(hash)

	for len(token) < length {
		token = base64.RawURLEncoding.EncodeToString([]byte(token))
	}

	if len(token) > length {
		token = token[:length]
	}

	return token, nil
}

// GenerateRandomClassicString generates a random string of the specified length with letters and digits.
// It generates a byte slice with the specified length and fills it with random bytes using the crypto/rand package.
// It iterates over each byte in the byte slice and replaces it with a character from the specified character set.
// It returns the byte slice converted to a string.
func GenerateRandomClassicString(length int) (string, error) {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	for i := 0; i < length; i++ {
		b[i] = chars[int(b[i])%len(chars)]
	}

	return string(b), nil
}

// ToInt64 converts a value of any type to an int64.
// It supports the following types: int, int8, int16, int32, int64, uint, uint8, uint16, uint32, and string.
// If the value is a string, it attempts to parse it as an int64 with base 10.
// If the value is not a supported type, it returns an error.
// If the value is an uint64, it is not supported due to the risk of truncation.
func ToInt64(i any) (int64, error) {
	switch v := i.(type) {
	case int:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case uint:
		return int64(v), nil
	case uint8:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	// Note: uint64 not supported due to risk of truncation.
	case string:
		return strconv.ParseInt(v, 10, 64)
	}

	return 0, fmt.Errorf("unable to convert type %T to int", i)
}

// Slugify generates a Url-friendly version of the input string.
// It converts all letters to lowercase and replaces spaces with hyphens.
// It allows alphanumeric characters, hyphens, and underscores to remain unchanged.
// It ignores non-ASCII characters.
func Slugify(s string) string {
	var buf bytes.Buffer

	for _, r := range s {
		switch {
		case r > unicode.MaxASCII:
			continue
		case unicode.IsLetter(r):
			buf.WriteRune(unicode.ToLower(r))
		case unicode.IsDigit(r), r == '_', r == '-':
			buf.WriteRune(r)
		case unicode.IsSpace(r):
			buf.WriteRune('-')
		}
	}

	return buf.String()
}

// UrlSetParam sets or updates a query parameter in a Url string.
// It parses the input Url string and sets the specified key-value pair in the query parameters.
// If the Url string cannot be parsed, an error is returned.
// The function returns the modified Url string with the updated query parameters.
func UrlSetParam(u string, key string, value interface{}) (string, error) {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return "", err
	}

	values := parsedURL.Query()
	values.Set(key, fmt.Sprintf("%v", value))

	parsedURL.RawQuery = values.Encode()
	return parsedURL.String(), nil
}

// UrlDelParam removes a query parameter with the specified key from the provided Url string.
// It parses the Url string, removes the specified parameter, encodes the updated query string, and returns the updated Url string.
// If there's an error parsing the Url, it returns an empty string and the error.
func UrlDelParam(u string, key string) (string, error) {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return "", err
	}

	values := parsedURL.Query()
	values.Del(key)

	parsedURL.RawQuery = values.Encode()
	return parsedURL.String(), nil
}

// CheckEmailDomainExistence checks if an email domain exists using SPF (Sender Policy Framework) record.
// It extracts the domain from the email address, performs a DNS lookup to get the TXT record of the domain,
// and checks if the TXT record contains the "v=spf1" flag. It returns true if the flag is found, and false otherwise.
func CheckEmailDomainExistence(addr string) (bool, error) {
	// Extract the domain from the email address
	domain := strings.Split(addr, "@")[1]

	// Get the TXT record of the domain
	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		return false, err
	}

	// Search for the "v=spf1" flag in the TXT record
	for _, txt := range txtRecords {
		if strings.Contains(txt, "v=spf1") {
			return true, nil
		}
	}

	return false, nil
}
