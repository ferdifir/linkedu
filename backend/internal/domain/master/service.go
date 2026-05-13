package master

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

// Classrooms
func (s *Service) GetClassrooms(tenantID uint) ([]domain.Classroom, error) {
	return s.repo.GetClassrooms(tenantID)
}

func (s *Service) GetClassroomByID(id, tenantID uint) (*domain.Classroom, error) {
	return s.repo.GetClassroomByID(id, tenantID)
}

func (s *Service) CreateClassroom(tenantID uint, name, gradeLevel string) (*domain.Classroom, error) {
	classroom := &domain.Classroom{
		TenantID:   tenantID,
		Name:       name,
		GradeLevel: gradeLevel,
	}
	if err := s.repo.CreateClassroom(classroom); err != nil {
		return nil, err
	}
	return classroom, nil
}

func (s *Service) UpdateClassroom(id, tenantID uint, name, gradeLevel string) (*domain.Classroom, error) {
	classroom, err := s.repo.GetClassroomByID(id, tenantID)
	if err != nil {
		return nil, ErrNotFound
	}
	classroom.Name = name
	classroom.GradeLevel = gradeLevel
	if err := s.repo.UpdateClassroom(classroom); err != nil {
		return nil, err
	}
	return classroom, nil
}

func (s *Service) DeleteClassroom(id, tenantID uint) error {
	_, err := s.repo.GetClassroomByID(id, tenantID)
	if err != nil {
		return ErrNotFound
	}
	return s.repo.DeleteClassroom(id, tenantID)
}

// Subjects
func (s *Service) GetSubjects(tenantID uint) ([]domain.Subject, error) {
	return s.repo.GetSubjects(tenantID)
}

func (s *Service) GetSubjectByID(id, tenantID uint) (*domain.Subject, error) {
	return s.repo.GetSubjectByID(id, tenantID)
}

func (s *Service) CreateSubject(tenantID uint, name, code string) (*domain.Subject, error) {
	subject := &domain.Subject{
		TenantID: tenantID,
		Name:     name,
		Code:     code,
	}
	if err := s.repo.CreateSubject(subject); err != nil {
		return nil, err
	}
	return subject, nil
}

func (s *Service) UpdateSubject(id, tenantID uint, name, code string) (*domain.Subject, error) {
	subject, err := s.repo.GetSubjectByID(id, tenantID)
	if err != nil {
		return nil, ErrNotFound
	}
	subject.Name = name
	subject.Code = code
	if err := s.repo.UpdateSubject(subject); err != nil {
		return nil, err
	}
	return subject, nil
}

func (s *Service) DeleteSubject(id, tenantID uint) error {
	_, err := s.repo.GetSubjectByID(id, tenantID)
	if err != nil {
		return ErrNotFound
	}
	return s.repo.DeleteSubject(id, tenantID)
}

// Teachers
func (s *Service) GetTeachers(tenantID uint) ([]domain.Teacher, error) {
	return s.repo.GetTeachers(tenantID)
}

func (s *Service) GetTeacherByID(id, tenantID uint) (*domain.Teacher, error) {
	return s.repo.GetTeacherByID(id, tenantID)
}

func (s *Service) CreateTeacher(tenantID, userID uint, nip, phone string) (*domain.Teacher, error) {
	teacher := &domain.Teacher{
		TenantID: tenantID,
		UserID:   userID,
		NIP:      nip,
		Phone:    phone,
	}
	if err := s.repo.CreateTeacher(teacher); err != nil {
		return nil, err
	}
	return teacher, nil
}

func (s *Service) UpdateTeacher(id, tenantID uint, nip, phone string) (*domain.Teacher, error) {
	teacher, err := s.repo.GetTeacherByID(id, tenantID)
	if err != nil {
		return nil, ErrNotFound
	}
	teacher.NIP = nip
	teacher.Phone = phone
	if err := s.repo.UpdateTeacher(teacher); err != nil {
		return nil, err
	}
	return teacher, nil
}

func (s *Service) DeleteTeacher(id, tenantID uint) error {
	_, err := s.repo.GetTeacherByID(id, tenantID)
	if err != nil {
		return ErrNotFound
	}
	return s.repo.DeleteTeacher(id, tenantID)
}

// Students
func (s *Service) GetStudents(tenantID uint) ([]domain.Student, error) {
	return s.repo.GetStudents(tenantID)
}

func (s *Service) GetStudentByID(id, tenantID uint) (*domain.Student, error) {
	return s.repo.GetStudentByID(id, tenantID)
}

func (s *Service) CreateStudent(tenantID, userID, classroomID uint, nis, rfidUID string) (*domain.Student, error) {
	student := &domain.Student{
		TenantID:    tenantID,
		UserID:      userID,
		ClassroomID: classroomID,
		NIS:         nis,
		RFIDUID:     rfidUID,
	}
	if err := s.repo.CreateStudent(student); err != nil {
		return nil, err
	}
	return student, nil
}

func (s *Service) UpdateStudent(id, tenantID, classroomID uint, nis, rfidUID string) (*domain.Student, error) {
	student, err := s.repo.GetStudentByID(id, tenantID)
	if err != nil {
		return nil, ErrNotFound
	}
	student.ClassroomID = classroomID
	student.NIS = nis
	student.RFIDUID = rfidUID
	if err := s.repo.UpdateStudent(student); err != nil {
		return nil, err
	}
	return student, nil
}

func (s *Service) DeleteStudent(id, tenantID uint) error {
	_, err := s.repo.GetStudentByID(id, tenantID)
	if err != nil {
		return ErrNotFound
	}
	return s.repo.DeleteStudent(id, tenantID)
}

// Parents
func (s *Service) GetParents(tenantID uint) ([]domain.Parent, error) {
	return s.repo.GetParents(tenantID)
}

func (s *Service) GetParentByID(id, tenantID uint) (*domain.Parent, error) {
	return s.repo.GetParentByID(id, tenantID)
}

func (s *Service) CreateParent(tenantID, userID, studentID uint) (*domain.Parent, error) {
	parent := &domain.Parent{
		TenantID:  tenantID,
		UserID:    userID,
		StudentID: studentID,
	}
	if err := s.repo.CreateParent(parent); err != nil {
		return nil, err
	}
	return parent, nil
}

func (s *Service) UpdateParent(id, tenantID, studentID uint) (*domain.Parent, error) {
	parent, err := s.repo.GetParentByID(id, tenantID)
	if err != nil {
		return nil, ErrNotFound
	}
	parent.StudentID = studentID
	if err := s.repo.UpdateParent(parent); err != nil {
		return nil, err
	}
	return parent, nil
}

func (s *Service) DeleteParent(id, tenantID uint) error {
	_, err := s.repo.GetParentByID(id, tenantID)
	if err != nil {
		return ErrNotFound
	}
	return s.repo.DeleteParent(id, tenantID)
}

// Academic Years
func (s *Service) GetAcademicYears(tenantID uint) ([]domain.AcademicYear, error) {
	return s.repo.GetAcademicYears(tenantID)
}

func (s *Service) GetAcademicYearByID(id, tenantID uint) (*domain.AcademicYear, error) {
	return s.repo.GetAcademicYearByID(id, tenantID)
}

func (s *Service) CreateAcademicYear(tenantID uint, name, startDate, endDate string, isActive bool) (*domain.AcademicYear, error) {
	// If isActive is true, set all other academic years to inactive
	if isActive {
		years, _ := s.repo.GetAcademicYears(tenantID)
		for _, year := range years {
			if year.IsActive {
				year.IsActive = false
				s.repo.UpdateAcademicYear(&year)
			}
		}
	}

	academicYear := &domain.AcademicYear{
		TenantID:  tenantID,
		Name:      name,
		StartDate: startDate,
		EndDate:   endDate,
		IsActive:  isActive,
	}
	if err := s.repo.CreateAcademicYear(academicYear); err != nil {
		return nil, err
	}
	return academicYear, nil
}

func (s *Service) UpdateAcademicYear(id, tenantID uint, name, startDate, endDate string, isActive bool) (*domain.AcademicYear, error) {
	academicYear, err := s.repo.GetAcademicYearByID(id, tenantID)
	if err != nil {
		return nil, ErrNotFound
	}

	// If setting to active, deactivate others
	if isActive && !academicYear.IsActive {
		years, _ := s.repo.GetAcademicYears(tenantID)
		for _, year := range years {
			if year.ID != id && year.IsActive {
				year.IsActive = false
				s.repo.UpdateAcademicYear(&year)
			}
		}
	}

	academicYear.Name = name
	academicYear.StartDate = startDate
	academicYear.EndDate = endDate
	academicYear.IsActive = isActive
	if err := s.repo.UpdateAcademicYear(academicYear); err != nil {
		return nil, err
	}
	return academicYear, nil
}

func (s *Service) DeleteAcademicYear(id, tenantID uint) error {
	_, err := s.repo.GetAcademicYearByID(id, tenantID)
	if err != nil {
		return ErrNotFound
	}
	return s.repo.DeleteAcademicYear(id, tenantID)
}
