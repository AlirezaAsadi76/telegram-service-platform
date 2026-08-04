package params

type GetTonPriceRequest struct{}
type GetTonPriceResponse struct {
	Price float64 `json:"price"`
}
