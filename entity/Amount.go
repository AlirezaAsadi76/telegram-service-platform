package entity

type Amount int64

func (a Amount) HasSufficient(amount Amount) bool {
	if a < amount {
		return false
	}
	return true
}
