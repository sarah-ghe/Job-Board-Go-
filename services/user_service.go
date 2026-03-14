package services

import (
	"errors"
	"job-board/models"		
	"job-board/repositories"
)


type UserRepository interface {
	Create(user *models.User) error
	GetByEmail(email string) (*models.User, error)
}

type UserService struct {
	Repo UserRepository
}

func (s *UserService) Register(user *models.User) error {

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	return s.Repo.Create(user)
}


func (s *UserService) Login(email, password string) (string, error) {

	user, err := s.Repo.GetByEmail(email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := GenerateJWT(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}