package attendance

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

func (r *Repository) CreateSession(session *domain.AttendanceSession) error {
	return r.db.Create(session).Error
}

func (r *Repository) UpdateSession(session *domain.AttendanceSession) error {
	return r.db.Save(session).Error
}

func (r *Repository) GetSessionByID(id uint, tenantID uint) (*domain.AttendanceSession, error) {
	var session domain.AttendanceSession
	err := r.db.Where("id = ? AND tenant_id = ?", id, tenantID).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *Repository) GetOpenSessionsByTeacher(teacherID uint, tenantID uint) ([]domain.AttendanceSession, error) {
	var sessions []domain.AttendanceSession
	err := r.db.Where("teacher_id = ? AND tenant_id = ? AND status = ?", teacherID, tenantID, "open").Find(&sessions).Error
	return sessions, err
}

func (r *Repository) CreateRecord(record *domain.AttendanceRecord) error {
	return r.db.Create(record).Error
}

func (r *Repository) GetRecordBySessionAndStudent(sessionID uint, studentID uint) (*domain.AttendanceRecord, error) {
	var record domain.AttendanceRecord
	err := r.db.Where("session_id = ? AND student_id = ?", sessionID, studentID).First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *Repository) GetRecordsBySession(sessionID uint, tenantID uint) ([]domain.AttendanceRecord, error) {
	var records []domain.AttendanceRecord
	err := r.db.Where("session_id = ? AND tenant_id = ?", sessionID, tenantID).Find(&records).Error
	return records, err
}

func (r *Repository) FindStudentByRFID(rfidUID string, tenantID uint) (*domain.Student, error) {
	var student domain.Student
	err := r.db.Where("rfid_uid = ? AND tenant_id = ?", rfidUID, tenantID).First(&student).Error
	if err != nil {
		return nil, err
	}
	return &student, nil
}

func (r *Repository) GetStudentsByClassroom(classroomID uint, tenantID uint) ([]domain.Student, error) {
	var students []domain.Student
	err := r.db.Where("classroom_id = ? AND tenant_id = ?", classroomID, tenantID).Find(&students).Error
	return students, err
}

func (r *Repository) GetScheduleByID(id uint, tenantID uint) (*domain.Schedule, error) {
	var schedule domain.Schedule
	err := r.db.Where("id = ? AND tenant_id = ?", id, tenantID).First(&schedule).Error
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (r *Repository) UpdateRecordStatus(studentID uint, date string, status string) error {
	// Find attendance records for the student on that date and update to permit
	// This requires joining with session to get the date
	return r.db.Model(&domain.AttendanceRecord{}).
		Joins("JOIN attendance_sessions ON attendance_records.session_id = attendance_sessions.id").
		Where("attendance_records.student_id = ? AND DATE(attendance_sessions.opened_at) = ? AND attendance_records.status = ?", 
			studentID, date, "absent").
		Update("status", status).Error
}
