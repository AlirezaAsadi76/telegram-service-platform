package smmentity

type Category string

var (
	Reaction Category = "reaction"
	Member   Category = "member"
	View     Category = "view"
	Comment  Category = "comment"
	Follower Category = "follower"
	Share    Category = "share"
)
