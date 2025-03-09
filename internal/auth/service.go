package auth

import (
	"errors"
	"urlshortener/internal/user"
	"urlshortener/pkg/di"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepository di.IUserRepository
}

func NewAuthService(userRepository di.IUserRepository) *AuthService {
	return &AuthService{UserRepository: userRepository}
}

func (service *AuthService) Register(email string, password string, name string) (string, error) {
	existedUser, _ := service.UserRepository.FindByEmail(email)
	if existedUser != nil {
		return "", errors.New(ErrUserExists)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", nil
	}

	user := &user.User{
		Email:    email,
		Password: string(hashedPassword),
		Name:     name,
	}

	_, err = service.UserRepository.Create(user)
	if err != nil {
		return "", err
	}

	return user.Email, nil
}

func (service *AuthService) Login(email string, password string) (string, error) {
	existedUser, _ := service.UserRepository.FindByEmail(email)
	if existedUser == nil {
		return "", errors.New(ErrUserDoesNotExist)
	}

	err := bcrypt.CompareHashAndPassword([]byte(existedUser.Password), []byte(password))
	if err != nil {
		return "", errors.New(ErrWrongCreds)
	}

	return existedUser.Email, nil
}
