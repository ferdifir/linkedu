package domain

import "time"

type Tenant struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"size:100;not null"`
	SchoolName  string `gorm:"size:200;not null"`
	Email       string `gorm:"size:100;uniqueIndex;not null"`
	Subdomain   string `gorm:"size:50;uniqueIndex"`
	Status      string `gorm:"size:20;default:'active'"` // active, suspended
	CreatedAt   time.Time
}

type User struct {
	ID           uint   `gorm:"primaryKey"`
	TenantID     *uint  `gorm:"index"`
	Name         string `gorm:"size:100;not null"`
	Email        string `gorm:"size:100;uniqueIndex;not null"`
	PasswordHash string `gorm:"size:255;not null"`
	Role         string `gorm:"size:20;not null"` // super_admin, school_admin, teacher, student, parent
	CreatedAt    time.Time
}

type AcademicYear struct {
	ID        uint `gorm:"primaryKey"`
	TenantID  uint `gorm:"index;not null"`
	Name      string `gorm:"size:50;not null"`
	StartDate string `gorm:"not null"` // YYYY-MM-DD
	EndDate   string `gorm:"not null"` // YYYY-MM-DD
	IsActive  bool   `gorm:"default:false"`
}

type Classroom struct {
	ID        uint `gorm:"primaryKey"`
	TenantID  uint `gorm:"index;not null"`
	Name      string `gorm:"size:50;not null"`
	GradeLevel string `gorm:"size:20"`
}

type Subject struct {
	ID       uint `gorm:"primaryKey"`
	TenantID uint `gorm:"index;not null"`
	Name     string `gorm:"size:100;not null"`
	Code     string `gorm:"size:20"`
}

type Teacher struct {
	ID       uint `gorm:"primaryKey"`
	TenantID uint `gorm:"index;not null"`
	UserID   uint `gorm:"uniqueIndex;not null"`
	NIP      string `gorm:"size:50"`
	Phone    string `gorm:"size:20"`
}

type Student struct {
	ID          uint `gorm:"primaryKey"`
	TenantID    uint `gorm:"index;not null"`
	UserID      uint `gorm:"uniqueIndex"`
	NIS         string `gorm:"size:50;uniqueIndex"`
	ClassroomID uint `gorm:"index"`
	RFIDUID     string `gorm:"size:50;index"`
}

type Parent struct {
	ID        uint `gorm:"primaryKey"`
	TenantID  uint `gorm:"index;not null"`
	UserID    uint `gorm:"uniqueIndex;not null"`
	StudentID uint `gorm:"index;not null"`
}

type Schedule struct {
	ID            uint `gorm:"primaryKey"`
	TenantID      uint `gorm:"index;not null"`
	AcademicYearID uint `gorm:"index;not null"`
	ClassroomID   uint `gorm:"index;not null"`
	SubjectID     uint `gorm:"index;not null"`
	TeacherID     uint `gorm:"index;not null"`
	DayOfWeek     int  `gorm:"not null"` // 0=Sunday, 1=Monday, etc.
	StartTime     string `gorm:"not null"` // HH:MM
	EndTime       string `gorm:"not null"` // HH:MM
}

type Event struct {
	ID          uint `gorm:"primaryKey"`
	TenantID    uint `gorm:"index;not null"`
	Name        string `gorm:"size:200;not null"`
	EventType   string `gorm:"size:50;not null"` // exam, holiday, etc.
	Date        string `gorm:"not null"` // YYYY-MM-DD
	StartTime   string // HH:MM
	EndTime     string // HH:MM
	Description string `gorm:"type:text"`
}

type AttendanceSession struct {
	ID        uint `gorm:"primaryKey"`
	TenantID  uint `gorm:"index;not null"`
	ScheduleID *uint `gorm:"index"`
	EventID   *uint `gorm:"index"`
	TeacherID uint `gorm:"index;not null"`
	OpenedAt  string `gorm:"not null"`
	ClosedAt  *string
	Status    string `gorm:"size:20;default:'open'"` // open, closed
}

type AttendanceRecord struct {
	ID        uint `gorm:"primaryKey"`
	TenantID  uint `gorm:"index;not null"`
	SessionID uint `gorm:"index;not null"`
	StudentID uint `gorm:"index;not null"`
	Status    string `gorm:"size:20;not null"` // present, absent, permit
	TappedAt  *string
	Method    *string // nfc, manual
}

type Permit struct {
	ID          uint `gorm:"primaryKey"`
	TenantID    uint `gorm:"index;not null"`
	StudentID   uint `gorm:"index;not null"`
	ParentID    uint `gorm:"index;not null"`
	SessionID   *uint `gorm:"index"`
	Date        string `gorm:"not null"` // YYYY-MM-DD
	Reason      string `gorm:"type:text;not null"`
	AttachmentPath string `gorm:"size:255"`
	Status      string `gorm:"size:20;default:'pending'"` // pending, approved, rejected
	ReviewedBy  *uint
	ReviewedAt  *string
	Notes       string `gorm:"type:text"`
}
