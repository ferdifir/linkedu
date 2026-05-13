package attendance

import (
	"linkedu/internal/domain"
	"time"
)

type OpenSessionInput struct {
	ScheduleID *uint `json:"schedule_id"`
	EventID    *uint `json:"event_id"`
	TeacherID  uint  `json:"teacher_id"`
	TenantID   uint  `json:"tenant_id"`
}

type TapInput struct {
	RFIDUID string `json:"rfid_uid"`
}

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) OpenSession(input OpenSessionInput) (*domain.AttendanceSession, error) {
	now := time.Now().Format("2006-01-02 15:04:05")
	
	session := &domain.AttendanceSession{
		TenantID:   input.TenantID,
		ScheduleID: input.ScheduleID,
		EventID:    input.EventID,
		TeacherID:  input.TeacherID,
		OpenedAt:   now,
		Status:     "open",
	}

	if err := s.repo.CreateSession(session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *Service) CloseSession(sessionID uint, tenantID uint) error {
	session, err := s.repo.GetSessionByID(sessionID, tenantID)
	if err != nil {
		return err
	}

	if session.Status == "closed" {
		return ErrSessionAlreadyClosed
	}

	// Get schedule to find classroom
	var classroomID uint
	if session.ScheduleID != nil {
		schedule, err := s.repo.GetScheduleByID(*session.ScheduleID, tenantID)
		if err != nil {
			return err
		}
		classroomID = schedule.ClassroomID
	}

	// Get all students in classroom
	students, err := s.repo.GetStudentsByClassroom(classroomID, tenantID)
	if err != nil {
		return err
	}

	// Insert absent records for students who haven't tapped
	now := time.Now().Format("2006-01-02 15:04:05")
	for _, student := range students {
		_, err := s.repo.GetRecordBySessionAndStudent(session.ID, student.ID)
		if err != nil {
			// Student has no record, create absent record
			record := &domain.AttendanceRecord{
				TenantID:  tenantID,
				SessionID: session.ID,
				StudentID: student.ID,
				Status:    "absent",
			}
			if err := s.repo.CreateRecord(record); err != nil {
				return err
			}
		}
	}

	// Update session status
	session.Status = "closed"
	session.ClosedAt = &now
	return s.repo.UpdateSession(session)
}


func (s *Service) TapStudent(sessionID uint, rfidUID string, tenantID uint) (*domain.AttendanceRecord, *domain.Student, error) {
	// Find student by RFID
	student, err := s.repo.FindStudentByRFID(rfidUID, tenantID)
	if err != nil {
		return nil, nil, ErrStudentNotFound
	}

	// Check if session exists and is open
	session, err := s.repo.GetSessionByID(sessionID, tenantID)
	if err != nil {
		return nil, nil, ErrSessionNotFound
	}

	if session.Status != "open" {
		return nil, nil, ErrSessionClosed
	}

	// Check if record already exists
	existingRecord, _ := s.repo.GetRecordBySessionAndStudent(sessionID, student.ID)
	if existingRecord != nil {
		// Return existing record (idempotent)
		return existingRecord, student, nil
	}

	// Create attendance record
	now := time.Now().Format("2006-01-02 15:04:05")
	record := &domain.AttendanceRecord{
		TenantID:  tenantID,
		SessionID: sessionID,
		StudentID: student.ID,
		Status:    "present",
		TappedAt:  &now,
		Method:    strPtr("nfc"),
	}

	if err := s.repo.CreateRecord(record); err != nil {
		return nil, nil, err
	}

	return record, student, nil
}

func strPtr(s string) *string {
	return &s
}

var (
	ErrSessionNotFound      = Error("Sesi presensi tidak ditemukan")
	ErrSessionClosed        = Error("Sesi presensi sudah ditutup")
	ErrSessionAlreadyClosed = Error("Sesi presensi sudah ditutup sebelumnya")
	ErrStudentNotFound      = Error("Siswa tidak ditemukan")
)

type Error string

func (e Error) Error() string {
	return string(e)
}
