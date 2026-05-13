package permit

import (
	"errors"
	"linkedu/internal/database"
	"linkedu/internal/domain"
	"time"
)

var (
	ErrNotFound      = errors.New("izin tidak ditemukan")
	ErrInvalidStatus = errors.New("status tidak valid")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetPermits(tenantID uint) ([]domain.Permit, error) {
	return s.repo.GetPermits(tenantID)
}

func (s *Service) GetPermitByID(id, tenantID uint) (*domain.Permit, error) {
	return s.repo.GetPermitByID(id, tenantID)
}

func (s *Service) GetPermitsByStudent(studentID, tenantID uint) ([]domain.Permit, error) {
	return s.repo.GetPermitsByStudent(studentID, tenantID)
}

func (s *Service) GetPermitsByParent(parentID, tenantID uint) ([]domain.Permit, error) {
	return s.repo.GetPermitsByParent(parentID, tenantID)
}

func (s *Service) GetPermitsByStatus(status string, tenantID uint) ([]domain.Permit, error) {
	return s.repo.GetPermitsByStatus(status, tenantID)
}

func (s *Service) CreatePermit(tenantID, studentID, parentID uint, sessionID *uint, date, reason, attachmentPath string) (*domain.Permit, error) {
	permit := &domain.Permit{
		TenantID:       tenantID,
		StudentID:      studentID,
		ParentID:       parentID,
		SessionID:      sessionID,
		Date:           date,
		Reason:         reason,
		AttachmentPath: attachmentPath,
		Status:         "pending",
	}
	if err := s.repo.CreatePermit(permit); err != nil {
		return nil, err
	}
	return permit, nil
}

func (s *Service) ReviewPermit(id, tenantID, reviewedBy uint, status, notes string) (*domain.Permit, error) {
	// Validate status
	if status != "approved" && status != "rejected" {
		return nil, ErrInvalidStatus
	}

	// Get permit first to verify it exists
	permit, err := s.repo.GetPermitByID(id, tenantID)
	if err != nil {
		return nil, ErrNotFound
	}

	// Can only review pending permits
	if permit.Status != "pending" {
		return nil, errors.New("izin sudah diproses")
	}

	reviewedAt := time.Now().Format(time.RFC3339)
	if err := s.repo.ReviewPermit(id, tenantID, status, &reviewedBy, reviewedAt, notes); err != nil {
		return nil, err
	}

	// If approved and there's an absent record for that date, update it to permit
	if status == "approved" {
		s.updateAttendanceRecordToPermit(permit.StudentID, permit.Date, tenantID)
	}

	// Return updated permit
	return s.repo.GetPermitByID(id, tenantID)
}

func (s *Service) updateAttendanceRecordToPermit(studentID uint, date string, tenantID uint) {
	// Find attendance sessions for that date
	var sessions []domain.AttendanceSession
	database.GetDB().Where("tenant_id = ? AND DATE(opened_at) = ?", tenantID, date).Find(&sessions)

	for _, session := range sessions {
		// Find absent record for this student in this session
		var record domain.AttendanceRecord
		result := database.GetDB().Where("session_id = ? AND student_id = ? AND status = 'absent'").
			First(&record)
		if result.Error == nil {
			// Update to permit
			record.Status = "permit"
			database.GetDB().Save(&record)
		}
		_ = session // avoid unused variable
	}
}
