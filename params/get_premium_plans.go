package params

type PremiumPlanInfo struct {
	ID     uint64
	Months uint8
}

type GetPremiumPlansResponse struct {
	Plans []PremiumPlanInfo
}
