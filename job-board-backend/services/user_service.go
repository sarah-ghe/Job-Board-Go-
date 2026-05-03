package services

import (
	"errors"
	"job-board/models"
	"job-board/utils"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByEmail(email string) (*models.User, error)
	GetByID(id int) (*models.User, error)
}

type UserService struct {
	Repo UserRepository
}

func (s *UserService) GetByEmail(email string) (*models.User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	return s.Repo.GetByEmail(email)
}

func (s *UserService) Register(user *models.User) error {
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))

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
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || password == "" {
		return "", errors.New("invalid credentials")
	}

	user, err := s.Repo.GetByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *UserService) GetUserByID(id int) (*models.User, error) {
	return s.Repo.GetByID(id)
}
