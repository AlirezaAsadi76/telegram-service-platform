package smmentity

type Category string

const (
	ReactionCategory Category = "reaction"
	MemberCategory   Category = "member"
	ViewCategory     Category = "view"
	CommentCategory  Category = "comment"
	FollowerCategory Category = "follower"
	ShareCategory    Category = "share"
	SaveCategory     Category = "save"
)
