package util

const (
	USD = "USD"
	EUR = "EUR"
	AUD = "AUD"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, AUD:
		return true
	}

	return false
}
