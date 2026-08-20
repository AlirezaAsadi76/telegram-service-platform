package smmentity

type Category struct {
	Name CategoryType
}

type CategoryType string

const (
	ReactionCategory CategoryType = "reaction"
	MemberCategory   CategoryType = "member"
	ViewCategory     CategoryType = "view"
	CommentCategory  CategoryType = "comment"
	FollowerCategory CategoryType = "follower"
	ShareCategory    CategoryType = "share"
	SaveCategory     CategoryType = "save"
)

func (c CategoryType) String() string {
	return string(c)
}
