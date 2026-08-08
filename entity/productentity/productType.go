package productentity

type ProductType string

const (
	Stars    ProductType = "stars"
	Premium  ProductType = "premium"
	Ads      ProductType = "ads"
	Reaction ProductType = "reaction"
	Member   ProductType = "member"
	View     ProductType = "view"
	Comment  ProductType = "comment"
	Follower ProductType = "follower"
	Share    ProductType = "share"
)
