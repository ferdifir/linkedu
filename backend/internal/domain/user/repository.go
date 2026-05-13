package user

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

func (r *Repository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *Repository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) FindByID(id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) FindByTenantAndRole(tenantID uint, role string) ([]domain.User, error) {
	var users []domain.User
	err := r.db.Where("tenant_id = ? AND role = ?", tenantID, role).Find(&users).Error
	return users, err
}

func (r *Repository) FindByTenant(tenantID uint) ([]domain.User, error) {
	var users []domain.User
	err := r.db.Where("tenant_id = ?", tenantID).Find(&users).Error
	return users, err
}

func (r *Repository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

func (r *Repository) Delete(id uint) error {
	return r.db.Delete(&domain.User{}, id).Error
}
