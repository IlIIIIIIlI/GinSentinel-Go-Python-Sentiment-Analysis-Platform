package services

import (
	"sentiment-service/internal/models"
	"sentiment-service/internal/repositories"
	"github.com/go-xorm/xorm"
)

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(engine *xorm.Engine) *UserService {
	return &UserService{userRepo: repositories.NewUserRepository(engine)}
}

func (us *UserService) GetUsers() ([]*models.User, error) {
	return us.userRepo.GetUsers()
}
