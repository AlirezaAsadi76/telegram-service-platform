package params

type StarPlanInfo struct {
	ID     uint64
	Amount uint32
	//Price  float64
}

type GetStarPlansResponse struct {
	Plans []StarPlanInfo
}
