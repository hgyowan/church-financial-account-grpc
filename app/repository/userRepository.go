package repository

type userRepository struct {
	repository *repository
}

func registerUserRepository(r *repository) {
	r.UserRepository = &userRepository{repository: r}
}
