package freekassa

const (
	mainUrl         = "https://pay.freekassa.ru/"
	balanceUrl      = "https://api.freekassa.ru/v1/balance"
	currenciesUrl   = "https://api.freekassa.ru/v1/currencies"
	ordersCreateUrl = "https://api.freekassa.ru/v1/orders/create"
)

var (
	AllowedFreekassaIPs = []string{"168.119.157.136", "168.119.60.227", "138.201.88.124", "178.154.197.79"}
)

const (
	CurrencyRUB = "RUB"
	CurrencyUSD = "USD"
	CurrencyEUR = "EUR"
	CurrencyKZT = "KZT"
	CurrencyUAH = "UAH"
)

type Config struct {
	ShopId           *uint32 `yaml:"shopId" json:"shop_id"`
	ApiKey           *string `yaml:"apiKey" json:"api_key"`
	FirstSecretWord  *string `yaml:"firstSecretWord" json:"first_secret_word"`
	SecondSecretWord *string `yaml:"secondSecretWord" json:"second_secret_word"`
}

func NewConfig(shopId uint32, apiKey, firstSecretWord, SecondSecretWord string) *Config {
	return &Config{
		ShopId:           &shopId,
		ApiKey:           &apiKey,
		FirstSecretWord:  &firstSecretWord,
		SecondSecretWord: &SecondSecretWord,
	}
}
