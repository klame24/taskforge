package service

import (
	"fmt"
	"net/mail"
	"taskforge/internal/auth"
	"taskforge/internal/models"
	"taskforge/internal/repository"
)

type UserService struct {
	userRepo repository.UserRepository
	jwtManager *auth.JWTManager
}

func NewUserService(userRepo repository.UserRepository, jwtManager *auth.JWTManager) *UserService {
	return &UserService{
		userRepo: userRepo,
		jwtManager: jwtManager,
	}
}

func (s *UserService) Register(firstName, lastName, email, password string) (*models.User, error) {
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, ErrInvalidEmail
	}

	if len(password) < 6 {
		return nil, ErrPasswordTooShort
	}

	existing, _ := s.userRepo.GetByEmail(email)
	if existing != nil {
		return nil, ErrEmailAlreadyExist
	}

	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		PasswordHash: hashedPassword,
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(email, password string) (*models.User, string, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}

	if !auth.CheckPasswordHash(password, user.PasswordHash) {
		return nil, "", ErrInvalidCredentials
	}

	token, err := s.jwtManager.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate JWT token: %w", err)
	}

	return user, token, nil
}
