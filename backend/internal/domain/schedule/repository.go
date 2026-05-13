package schedule

import (
	"linkedu/internal/database"
	"linkedu/internal/domain"
)

type Repository interface {
	// Schedules
	GetSchedules(tenantID uint) ([]domain.Schedule, error)
	GetScheduleByID(id, tenantID uint) (*domain.Schedule, error)
	CreateSchedule(schedule *domain.Schedule) error
	UpdateSchedule(schedule *domain.Schedule) error
	DeleteSchedule(id, tenantID uint) error
	GetSchedulesByClassroom(classroomID, tenantID uint) ([]domain.Schedule, error)

	// Events
	GetEvents(tenantID uint) ([]domain.Event, error)
	GetEventByID(id, tenantID uint) (*domain.Event, error)
	CreateEvent(event *domain.Event) error
	UpdateEvent(event *domain.Event) error
	DeleteEvent(id, tenantID uint) error
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

// Schedules
func (r *repository) GetSchedules(tenantID uint) ([]domain.Schedule, error) {
	var schedules []domain.Schedule
	err := database.GetDB().Preload("AcademicYear").Preload("Classroom").Preload("Subject").Preload("Teacher").
		Where("tenant_id = ?", tenantID).Find(&schedules).Error
	return schedules, err
}

func (r *repository) GetScheduleByID(id, tenantID uint) (*domain.Schedule, error) {
	var schedule domain.Schedule
	err := database.GetDB().Preload("AcademicYear").Preload("Classroom").Preload("Subject").Preload("Teacher").
		Where("id = ? AND tenant_id = ?", id, tenantID).First(&schedule).Error
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (r *repository) CreateSchedule(schedule *domain.Schedule) error {
	return database.GetDB().Create(schedule).Error
}

func (r *repository) UpdateSchedule(schedule *domain.Schedule) error {
	return database.GetDB().Save(schedule).Error
}

func (r *repository) DeleteSchedule(id, tenantID uint) error {
	return database.GetDB().Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&domain.Schedule{}).Error
}

func (r *repository) GetSchedulesByClassroom(classroomID, tenantID uint) ([]domain.Schedule, error) {
	var schedules []domain.Schedule
	err := database.GetDB().Preload("AcademicYear").Preload("Classroom").Preload("Subject").Preload("Teacher").
		Where("tenant_id = ? AND classroom_id = ?", tenantID, classroomID).Find(&schedules).Error
	return schedules, err
}

// Events
func (r *repository) GetEvents(tenantID uint) ([]domain.Event, error) {
	var events []domain.Event
	err := database.GetDB().Where("tenant_id = ?", tenantID).Order("date DESC").Find(&events).Error
	return events, err
}

func (r *repository) GetEventByID(id, tenantID uint) (*domain.Event, error) {
	var event domain.Event
	err := database.GetDB().Where("id = ? AND tenant_id = ?", id, tenantID).First(&event).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *repository) CreateEvent(event *domain.Event) error {
	return database.GetDB().Create(event).Error
}

func (r *repository) UpdateEvent(event *domain.Event) error {
	return database.GetDB().Save(event).Error
}

func (r *repository) DeleteEvent(id, tenantID uint) error {
	return database.GetDB().Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&domain.Event{}).Error
}
