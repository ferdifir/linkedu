package tenant

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

type RegisterTenantInput struct {
	Name       string `json:"name"`
	SchoolName string `json:"school_name"`
	Email      string `json:"email"`
	Subdomain  string `json:"subdomain"`
	Password   string `json:"password"`
}

func (s *Service) RegisterTenant(input RegisterTenantInput) (*domain.Tenant, *domain.User, error) {
	// Check if email already exists
	existingTenant, _ := s.repo.FindByEmail(input.Email)
	if existingTenant != nil {
		return nil, nil, ErrEmailExists
	}

	// Check if subdomain already exists
	if input.Subdomain != "" {
		existingBySubdomain, _ := s.repo.FindBySubdomain(input.Subdomain)
		if existingBySubdomain != nil {
			return nil, nil, ErrSubdomainExists
		}
	}

	// Create tenant
	tenant := &domain.Tenant{
		Name:       input.Name,
		SchoolName: input.SchoolName,
		Email:      input.Email,
		Subdomain:  input.Subdomain,
		Status:     "active",
	}

	if err := s.repo.Create(tenant); err != nil {
		return nil, nil, err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, err
	}

	// Create admin user
	user := &domain.User{
		TenantID:     &tenant.ID,
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: string(hashedPassword),
		Role:         "school_admin",
	}

	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, nil, err
	}

	return tenant, user, nil
}

var (
	ErrEmailExists     = Error("Email sudah terdaftar")
	ErrSubdomainExists = Error("Subdomain sudah digunakan")
)

type Error string

func (e Error) Error() string {
	return string(e)
}
