package freekassa

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

func createSHA256Signature(apiKey string, values map[string]string) string {
	// Get the keys and sort them in alphabetical order
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Concat the key values using the "|" sign
	var dataBuilder strings.Builder
	for i, k := range keys {
		dataBuilder.WriteString(values[k])
		if i < len(keys)-1 {
			dataBuilder.WriteString("|")
		}
	}
	data := dataBuilder.String()

	// Hash the resulting string to sha256 using the API key
	mac := hmac.New(sha256.New, []byte(apiKey))
	mac.Write([]byte(data))
	return hex.EncodeToString(mac.Sum(nil))
}

func createPaymentFormSignature(cfg *Config, amount string, currency string, OrderId string) string {
	data := fmt.Sprintf("%d:%s:%s:%s:%s", *cfg.ShopId, amount, *cfg.FirstSecretWord, currency, OrderId)
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}

func CreateNotificationSignature(cfg *Config, amount string, OrderId string) string {
	data := fmt.Sprintf("%d:%s:%s:%s", *cfg.ShopId, amount, *cfg.SecondSecretWord, OrderId)
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}
