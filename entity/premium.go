package entity

import "time"

type PremiumPlan struct {
	ID        uint64
	Duration  PremiumDuration
	Active    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PremiumDuration uint8

const (
	PremiumDurationOneMonth    PremiumDuration = 1
	PremiumDurationThreeMonth  PremiumDuration = 3
	PremiumDurationSixMonth    PremiumDuration = 6
	PremiumDurationTwelveMonth PremiumDuration = 12
)

func (d PremiumDuration) Months() uint8 {
	return uint8(d)
}

func (d PremiumDuration) String() string {

	switch d {

	case PremiumDurationOneMonth:
		return "1 month"

	case PremiumDurationThreeMonth:
		return "3 months"

	case PremiumDurationSixMonth:
		return "6 months"

	case PremiumDurationTwelveMonth:
		return "12 months"

	default:
		return "unknown"
	}
}

func (d PremiumDuration) IsValid() bool {

	switch d {

	case PremiumDurationOneMonth,
		PremiumDurationThreeMonth,
		PremiumDurationSixMonth,
		PremiumDurationTwelveMonth:

		return true

	default:
		return false
	}
}
