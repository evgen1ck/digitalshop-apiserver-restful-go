package freekassa

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func NewOrderUrl(cfg *Config, amount float64, currency, orderName string) string {
	var signature, orderUrl string
	orderName = strings.ReplaceAll(strings.TrimSpace(orderName), " ", "_")

	signature = createPaymentFormSignature(cfg, formatMoney(amount), currency, orderName)

	orderUrl = fmt.Sprintf("%s?m=%d&oa=%s&currency=%s&o=%s&s=%s",
		mainUrl, *cfg.ShopId, formatMoney(amount), currency, orderName, signature)

	return orderUrl
}

type newOrderPlusPayload struct {
	ShopId    int     `json:"shopId"`
	Nonce     int     `json:"nonce"`
	Signature string  `json:"signature"`
	OrderName string  `json:"paymentId"`
	PaymentId int     `json:"i"`
	Email     string  `json:"email"`
	ClientIp  string  `json:"ip"`
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency"`
}

type newOrderPlusResponse struct {
	Type      string `json:"type"`
	OrderId   int    `json:"orderId"`
	OrderHash string `json:"orderHash"`
	Location  string `json:"location"`
}

func CreateOrder(cfg *Config, amount float64, currency, orderName, ip, email string, i int) (int, string, string, error) {
	var signature string
	var localNewOrderPlusResponse newOrderPlusResponse
	nonce := int(time.Now().Unix())

	values := map[string]string{
		"shopId":    strconv.Itoa(int(*cfg.ShopId)),
		"nonce":     strconv.Itoa(nonce),
		"paymentId": orderName,
		"i":         strconv.Itoa(i),
		"email":     email,
		"ip":        ip,
		"amount":    formatMoney(amount),
		"currency":  currency,
	}
	signature = createSHA256Signature(*cfg.ApiKey, values)

	payload := &newOrderPlusPayload{
		ShopId:    int(*cfg.ShopId),
		Nonce:     nonce,
		Signature: signature,
		OrderName: orderName,
		PaymentId: i,
		Email:     email,
		ClientIp:  ip,
		Amount:    amount,
		Currency:  currency,
	}

	body, err := sendQuery(ordersCreateUrl, payload)
	if err != nil {
		return 0, "", "", err
	}

	if err = json.Unmarshal(body, &localNewOrderPlusResponse); err != nil {
		return 0, "", "", fmt.Errorf("error unmarshalling response body: %v", err)
	}

	if localNewOrderPlusResponse.Type != "success" {
		return 0, "", "", fmt.Errorf("error getting in response success status in json format")
	}

	return localNewOrderPlusResponse.OrderId, localNewOrderPlusResponse.OrderHash, localNewOrderPlusResponse.Location, err
}
