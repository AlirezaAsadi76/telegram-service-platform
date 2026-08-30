package helpers

import "github.com/shopspring/decimal"

func FormatPrice(price decimal.Decimal) string {

	str := price.StringFixed(0)

	var result []byte
	for i, c := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, byte(c))
	}
	return string(result)
}
