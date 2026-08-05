package entity

type AdsPlan struct {
	View  int64   `json:"view"`
	CPM   float64 `json:"cpm"`
	Price float64 `json:"price"`
}
