package entity

type Role uint8

const (
	UserRole Role = iota + 1
	AdminRole
)

const (
	UserRoleString  = "user"
	AdminRoleString = "admin"
)

func (r Role) String() string {
	switch r {
	case UserRole:
		return UserRoleString
	case AdminRole:

		return AdminRoleString
	}
	return ""
}
func MapToRoleEntity(role string) Role {
	switch role {
	case UserRoleString:
		return UserRole
	case AdminRoleString:
		return AdminRole
	}
	return Role(0)
}
