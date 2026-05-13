package permit

import (
	"linkedu/internal/database"
	"linkedu/internal/domain"
)

type Repository interface {
	GetPermits(tenantID uint) ([]domain.Permit, error)
	GetPermitByID(id, tenantID uint) (*domain.Permit, error)
	GetPermitsByStudent(studentID, tenantID uint) ([]domain.Permit, error)
	GetPermitsByParent(parentID, tenantID uint) ([]domain.Permit, error)
	GetPermitsByStatus(status string, tenantID uint) ([]domain.Permit, error)
	CreatePermit(permit *domain.Permit) error
	UpdatePermit(permit *domain.Permit) error
	ReviewPermit(id, tenantID uint, status string, reviewedBy *uint, reviewedAt, notes string) error
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) GetPermits(tenantID uint) ([]domain.Permit, error) {
	var permits []domain.Permit
	err := database.GetDB().Preload("Student").Preload("Parent").Preload("Reviewer").
		Where("tenant_id = ?", tenantID).Order("created_at DESC").Find(&permits).Error
	return permits, err
}

func (r *repository) GetPermitByID(id, tenantID uint) (*domain.Permit, error) {
	var permit domain.Permit
	err := database.GetDB().Preload("Student").Preload("Parent").Preload("Reviewer").
		Where("id = ? AND tenant_id = ?", id, tenantID).First(&permit).Error
	if err != nil {
		return nil, err
	}
	return &permit, nil
}

func (r *repository) GetPermitsByStudent(studentID, tenantID uint) ([]domain.Permit, error) {
	var permits []domain.Permit
	err := database.GetDB().Preload("Student").Preload("Parent").
		Where("tenant_id = ? AND student_id = ?", tenantID, studentID).
		Order("created_at DESC").Find(&permits).Error
	return permits, err
}

func (r *repository) GetPermitsByParent(parentID, tenantID uint) ([]domain.Permit, error) {
	var permits []domain.Permit
	err := database.GetDB().Preload("Student").Preload("Parent").
		Where("tenant_id = ? AND parent_id = ?", tenantID, parentID).
		Order("created_at DESC").Find(&permits).Error
	return permits, err
}

func (r *repository) GetPermitsByStatus(status string, tenantID uint) ([]domain.Permit, error) {
	var permits []domain.Permit
	err := database.GetDB().Preload("Student").Preload("Parent").
		Where("tenant_id = ? AND status = ?", tenantID, status).
		Order("created_at DESC").Find(&permits).Error
	return permits, err
}

func (r *repository) CreatePermit(permit *domain.Permit) error {
	return database.GetDB().Create(permit).Error
}

func (r *repository) UpdatePermit(permit *domain.Permit) error {
	return database.GetDB().Save(permit).Error
}

func (r *repository) ReviewPermit(id, tenantID uint, status string, reviewedBy *uint, reviewedAt, notes string) error {
	result := database.GetDB().Model(&domain.Permit{}).
		Where("id = ? AND tenant_id = ?", id, tenantID).
		Updates(map[string]interface{}{
			"status":      status,
			"reviewed_by": reviewedBy,
			"reviewed_at": reviewedAt,
			"notes":       notes,
		})
	return result.Error
}
