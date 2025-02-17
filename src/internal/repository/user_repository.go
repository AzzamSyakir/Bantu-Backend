package repository

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	UserRepository := &UserRepository{}
	return UserRepository
}
