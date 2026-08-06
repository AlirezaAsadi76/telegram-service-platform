package exchangerate

import (
	"encoding/json"
	"strconv"
)

type tonUsdPriceResponse map[string]price

type price struct {
	USD float64 `json:"usd"`
}

type marketsResponse struct {
	Result map[string]market `json:"result"`
}

type market struct {
	Symbol string `json:"symbol"`
	Stats  stats  `json:"stats"`
}

type stats struct {
	LastPrice Number `json:"lastPrice"`
}
type Number float64

func (n *Number) UnmarshalJSON(data []byte) error {

	var num float64

	if err := json.Unmarshal(data, &num); err == nil {
		*n = Number(num)
		return nil
	}

	var str string

	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	if str == "" || str == "-" {
		*n = 0
		return nil
	}

	value, err := strconv.ParseFloat(str, 64)

	if err != nil {
		return err
	}

	*n = Number(value)

	return nil
}
