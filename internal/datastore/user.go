package datastore

type (
	User struct {
		ID   uint64
		Name string
	}

	UserFilter struct {
		IDs   []uint64
		Names []string
	}
)

type UserRepo interface {
	CRUD[User, UserFilter]
}
