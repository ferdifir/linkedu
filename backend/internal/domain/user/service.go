package user

import (
	"linkedu/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginOutput struct {
	User       *domain.User `json:"user"`
	AccessToken string      `json:"access_token"`
}

func (s *Service) Login(input LoginInput) (*LoginOutput, error) {
	// Find user by email
	user, err := s.repo.FindByEmail(input.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	return &LoginOutput{
		User: user,
	}, nil
}

var ErrInvalidCredentials = Error("Email atau password salah")

type Error string

func (e Error) Error() string {
	return string(e)
}

// CreateUser creates a new user (exported for use in tenant service)
func CreateUser(user *domain.User) error {
	repo := NewRepository()
	return repo.Create(user)
}
