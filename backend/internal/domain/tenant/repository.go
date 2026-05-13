package tenant

import (
	"linkedu/internal/database"
	"linkedu/internal/domain"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository() *Repository {
	return &Repository{db: database.GetDB()}
}

func (r *Repository) Create(tenant *domain.Tenant) error {
	return r.db.Create(tenant).Error
}

func (r *Repository) FindByEmail(email string) (*domain.Tenant, error) {
	var tenant domain.Tenant
	err := r.db.Where("email = ?", email).First(&tenant).Error
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *Repository) FindByID(id uint) (*domain.Tenant, error) {
	var tenant domain.Tenant
	err := r.db.First(&tenant, id).Error
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *Repository) FindBySubdomain(subdomain string) (*domain.Tenant, error) {
	var tenant domain.Tenant
	err := r.db.Where("subdomain = ?", subdomain).First(&tenant).Error
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *Repository) Update(tenant *domain.Tenant) error {
	return r.db.Save(tenant).Error
}
