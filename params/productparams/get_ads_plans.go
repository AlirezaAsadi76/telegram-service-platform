package productparams

type AdsPlanInfo struct {
	ID             uint64
	Views          uint64
	CPM            float64
	DailyViewLimit uint64
}

type GetAdsPlansResponse struct {
	Plans []AdsPlanInfo
}
