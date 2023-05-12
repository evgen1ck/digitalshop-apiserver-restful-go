package freekassa

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type otherPayload struct {
	ShopId    int    `json:"shopId"`
	Nonce     int    `json:"nonce"`
	Signature string `json:"signature"`
}

type _balance struct {
	Currency string `json:"currency"`
	Value    string `json:"value"`
}

type Balance struct {
	Currency string  `json:"currency"`
	Value    float64 `json:"value"`
}

type balanceResponse struct {
	Type    string     `json:"type"`
	Balance []_balance `json:"balance"`
}

func Balances(cfg *Config) ([]Balance, error) {
	var signature string
	var localBalanceResponse balanceResponse
	nonce := int(time.Now().Unix())

	values := map[string]string{
		"shopId": strconv.Itoa(int(*cfg.ShopId)),
		"nonce":  strconv.Itoa(nonce),
	}
	signature = createSHA256Signature(*cfg.ApiKey, values)

	payload := &otherPayload{
		ShopId:    int(*cfg.ShopId),
		Nonce:     nonce,
		Signature: signature,
	}

	body, err := sendQuery(balanceUrl, payload)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(body, &localBalanceResponse); err != nil {
		return nil, fmt.Errorf("error unmarshalling response body: %v", err)
	}

	if localBalanceResponse.Type != "success" {
		return nil, fmt.Errorf("error getting in response success status in json format")
	}

	var newBalance []Balance
	for _, bal := range localBalanceResponse.Balance {

		value, err := strconv.ParseFloat(bal.Value, 64)
		if err != nil {
			return nil, fmt.Errorf("error converting value to float64: %v", err)
		}

		balanceFloat := Balance{
			Currency: bal.Currency,
			Value:    value,
		}

		newBalance = append(newBalance, balanceFloat)
	}

	return newBalance, err
}

type Currency struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Currency   string `json:"currency"`
	IsEnabled  int    `json:"is_enabled"`
	IsFavorite int    `json:"is_favorite"`
}

type currenciesResponse struct {
	Type     string     `json:"type"`
	Currency []Currency `json:"balance"`
}

func Currencies(cfg *Config) ([]Currency, error) {
	var signature string
	var localCurrenciesResponse currenciesResponse
	nonce := int(time.Now().Unix())

	values := map[string]string{
		"shopId": strconv.Itoa(int(*cfg.ShopId)),
		"nonce":  strconv.Itoa(nonce),
	}
	signature = createSHA256Signature(*cfg.ApiKey, values)

	payload := &otherPayload{
		ShopId:    int(*cfg.ShopId),
		Nonce:     nonce,
		Signature: signature,
	}

	body, err := sendQuery(currenciesUrl, payload)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(body, &localCurrenciesResponse); err != nil {
		return nil, fmt.Errorf("error unmarshalling response body: %v", err)
	}

	if localCurrenciesResponse.Type != "success" {
		return nil, fmt.Errorf("error getting in response success status in json format")
	}

	return localCurrenciesResponse.Currency, err
}
