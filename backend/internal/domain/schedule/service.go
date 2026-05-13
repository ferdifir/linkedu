package schedule

import (
	"errors"
	"linkedu/internal/domain"
)

var (
	ErrNotFound = errors.New("data tidak ditemukan")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// Schedules
func (s *Service) GetSchedules(tenantID uint) ([]domain.Schedule, error) {
	return s.repo.GetSchedules(tenantID)
}

func (s *Service) GetScheduleByID(id, tenantID uint) (*domain.Schedule, error) {
	return s.repo.GetScheduleByID(id, tenantID)
}

func (s *Service) CreateSchedule(tenantID, academicYearID, classroomID, subjectID, teacherID uint, dayOfWeek int, startTime, endTime string) (*domain.Schedule, error) {
	schedule := &domain.Schedule{
		TenantID:       tenantID,
		AcademicYearID: academicYearID,
		ClassroomID:    classroomID,
		SubjectID:      subjectID,
		TeacherID:      teacherID,
		DayOfWeek:      dayOfWeek,
		StartTime:      startTime,
		EndTime:        endTime,
	}
	if err := s.repo.CreateSchedule(schedule); err != nil {
		return nil, err
	}
	return schedule, nil
}

func (s *Service) UpdateSchedule(id, tenantID, academicYearID, classroomID, subjectID, teacherID uint, dayOfWeek int, startTime, endTime string) (*domain.Schedule, error) {
	schedule, err := s.repo.GetScheduleByID(id, tenantID)
	if err != nil {
		return nil, ErrNotFound
	}
	schedule.AcademicYearID = academicYearID
	schedule.ClassroomID = classroomID
	schedule.SubjectID = subjectID
	schedule.TeacherID = teacherID
	schedule.DayOfWeek = dayOfWeek
	schedule.StartTime = startTime
	schedule.EndTime = endTime
	if err := s.repo.UpdateSchedule(schedule); err != nil {
		return nil, err
	}
	return schedule, nil
}

func (s *Service) DeleteSchedule(id, tenantID uint) error {
	_, err := s.repo.GetScheduleByID(id, tenantID)
	if err != nil {
		return ErrNotFound
	}
	return s.repo.DeleteSchedule(id, tenantID)
}

func (s *Service) GetSchedulesByClassroom(classroomID, tenantID uint) ([]domain.Schedule, error) {
	return s.repo.GetSchedulesByClassroom(classroomID, tenantID)
}

// Events
func (s *Service) GetEvents(tenantID uint) ([]domain.Event, error) {
	return s.repo.GetEvents(tenantID)
}

func (s *Service) GetEventByID(id, tenantID uint) (*domain.Event, error) {
	return s.repo.GetEventByID(id, tenantID)
}

func (s *Service) CreateEvent(tenantID uint, name, eventType, date, startTime, endTime, description string) (*domain.Event, error) {
	event := &domain.Event{
		TenantID:    tenantID,
		Name:        name,
		EventType:   eventType,
		Date:        date,
		StartTime:   startTime,
		EndTime:     endTime,
		Description: description,
	}
	if err := s.repo.CreateEvent(event); err != nil {
		return nil, err
	}
	return event, nil
}

func (s *Service) UpdateEvent(id, tenantID uint, name, eventType, date, startTime, endTime, description string) (*domain.Event, error) {
	event, err := s.repo.GetEventByID(id, tenantID)
	if err != nil {
		return nil, ErrNotFound
	}
	event.Name = name
	event.EventType = eventType
	event.Date = date
	event.StartTime = startTime
	event.EndTime = endTime
	event.Description = description
	if err := s.repo.UpdateEvent(event); err != nil {
		return nil, err
	}
	return event, nil
}

func (s *Service) DeleteEvent(id, tenantID uint) error {
	_, err := s.repo.GetEventByID(id, tenantID)
	if err != nil {
		return ErrNotFound
	}
	return s.repo.DeleteEvent(id, tenantID)
}
