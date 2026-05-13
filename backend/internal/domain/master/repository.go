package master

import (
	"linkedu/internal/database"
	"linkedu/internal/domain"
)

type Repository interface {
	// Classrooms
	GetClassrooms(tenantID uint) ([]domain.Classroom, error)
	GetClassroomByID(id, tenantID uint) (*domain.Classroom, error)
	CreateClassroom(classroom *domain.Classroom) error
	UpdateClassroom(classroom *domain.Classroom) error
	DeleteClassroom(id, tenantID uint) error

	// Subjects
	GetSubjects(tenantID uint) ([]domain.Subject, error)
	GetSubjectByID(id, tenantID uint) (*domain.Subject, error)
	CreateSubject(subject *domain.Subject) error
	UpdateSubject(subject *domain.Subject) error
	DeleteSubject(id, tenantID uint) error

	// Teachers
	GetTeachers(tenantID uint) ([]domain.Teacher, error)
	GetTeacherByID(id, tenantID uint) (*domain.Teacher, error)
	CreateTeacher(teacher *domain.Teacher) error
	UpdateTeacher(teacher *domain.Teacher) error
	DeleteTeacher(id, tenantID uint) error

	// Students
	GetStudents(tenantID uint) ([]domain.Student, error)
	GetStudentByID(id, tenantID uint) (*domain.Student, error)
	CreateStudent(student *domain.Student) error
	UpdateStudent(student *domain.Student) error
	DeleteStudent(id, tenantID uint) error

	// Parents
	GetParents(tenantID uint) ([]domain.Parent, error)
	GetParentByID(id, tenantID uint) (*domain.Parent, error)
	CreateParent(parent *domain.Parent) error
	UpdateParent(parent *domain.Parent) error
	DeleteParent(id, tenantID uint) error

	// Academic Years
	GetAcademicYears(tenantID uint) ([]domain.AcademicYear, error)
	GetAcademicYearByID(id, tenantID uint) (*domain.AcademicYear, error)
	CreateAcademicYear(academicYear *domain.AcademicYear) error
	UpdateAcademicYear(academicYear *domain.AcademicYear) error
	DeleteAcademicYear(id, tenantID uint) error
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

// Classrooms
func (r *repository) GetClassrooms(tenantID uint) ([]domain.Classroom, error) {
	var classrooms []domain.Classroom
	err := database.GetDB().Where("tenant_id = ?", tenantID).Find(&classrooms).Error
	return classrooms, err
}

func (r *repository) GetClassroomByID(id, tenantID uint) (*domain.Classroom, error) {
	var classroom domain.Classroom
	err := database.GetDB().Where("id = ? AND tenant_id = ?", id, tenantID).First(&classroom).Error
	if err != nil {
		return nil, err
	}
	return &classroom, nil
}

func (r *repository) CreateClassroom(classroom *domain.Classroom) error {
	return database.GetDB().Create(classroom).Error
}

func (r *repository) UpdateClassroom(classroom *domain.Classroom) error {
	return database.GetDB().Save(classroom).Error
}

func (r *repository) DeleteClassroom(id, tenantID uint) error {
	return database.GetDB().Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&domain.Classroom{}).Error
}

// Subjects
func (r *repository) GetSubjects(tenantID uint) ([]domain.Subject, error) {
	var subjects []domain.Subject
	err := database.GetDB().Where("tenant_id = ?", tenantID).Find(&subjects).Error
	return subjects, err
}

func (r *repository) GetSubjectByID(id, tenantID uint) (*domain.Subject, error) {
	var subject domain.Subject
	err := database.GetDB().Where("id = ? AND tenant_id = ?", id, tenantID).First(&subject).Error
	if err != nil {
		return nil, err
	}
	return &subject, nil
}

func (r *repository) CreateSubject(subject *domain.Subject) error {
	return database.GetDB().Create(subject).Error
}

func (r *repository) UpdateSubject(subject *domain.Subject) error {
	return database.GetDB().Save(subject).Error
}

func (r *repository) DeleteSubject(id, tenantID uint) error {
	return database.GetDB().Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&domain.Subject{}).Error
}

// Teachers
func (r *repository) GetTeachers(tenantID uint) ([]domain.Teacher, error) {
	var teachers []domain.Teacher
	err := database.GetDB().Preload("User").Where("tenant_id = ?", tenantID).Find(&teachers).Error
	return teachers, err
}

func (r *repository) GetTeacherByID(id, tenantID uint) (*domain.Teacher, error) {
	var teacher domain.Teacher
	err := database.GetDB().Preload("User").Where("id = ? AND tenant_id = ?", id, tenantID).First(&teacher).Error
	if err != nil {
		return nil, err
	}
	return &teacher, nil
}

func (r *repository) CreateTeacher(teacher *domain.Teacher) error {
	return database.GetDB().Create(teacher).Error
}

func (r *repository) UpdateTeacher(teacher *domain.Teacher) error {
	return database.GetDB().Save(teacher).Error
}

func (r *repository) DeleteTeacher(id, tenantID uint) error {
	return database.GetDB().Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&domain.Teacher{}).Error
}

// Students
func (r *repository) GetStudents(tenantID uint) ([]domain.Student, error) {
	var students []domain.Student
	err := database.GetDB().Preload("Classroom").Preload("User").Where("tenant_id = ?", tenantID).Find(&students).Error
	return students, err
}

func (r *repository) GetStudentByID(id, tenantID uint) (*domain.Student, error) {
	var student domain.Student
	err := database.GetDB().Preload("Classroom").Preload("User").Where("id = ? AND tenant_id = ?", id, tenantID).First(&student).Error
	if err != nil {
		return nil, err
	}
	return &student, nil
}

func (r *repository) CreateStudent(student *domain.Student) error {
	return database.GetDB().Create(student).Error
}

func (r *repository) UpdateStudent(student *domain.Student) error {
	return database.GetDB().Save(student).Error
}

func (r *repository) DeleteStudent(id, tenantID uint) error {
	return database.GetDB().Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&domain.Student{}).Error
}

// Parents
func (r *repository) GetParents(tenantID uint) ([]domain.Parent, error) {
	var parents []domain.Parent
	err := database.GetDB().Preload("User").Preload("Student").Where("tenant_id = ?", tenantID).Find(&parents).Error
	return parents, err
}

func (r *repository) GetParentByID(id, tenantID uint) (*domain.Parent, error) {
	var parent domain.Parent
	err := database.GetDB().Preload("User").Preload("Student").Where("id = ? AND tenant_id = ?", id, tenantID).First(&parent).Error
	if err != nil {
		return nil, err
	}
	return &parent, nil
}

func (r *repository) CreateParent(parent *domain.Parent) error {
	return database.GetDB().Create(parent).Error
}

func (r *repository) UpdateParent(parent *domain.Parent) error {
	return database.GetDB().Save(parent).Error
}

func (r *repository) DeleteParent(id, tenantID uint) error {
	return database.GetDB().Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&domain.Parent{}).Error
}

// Academic Years
func (r *repository) GetAcademicYears(tenantID uint) ([]domain.AcademicYear, error) {
	var academicYears []domain.AcademicYear
	err := database.GetDB().Where("tenant_id = ?", tenantID).Find(&academicYears).Error
	return academicYears, err
}

func (r *repository) GetAcademicYearByID(id, tenantID uint) (*domain.AcademicYear, error) {
	var academicYear domain.AcademicYear
	err := database.GetDB().Where("id = ? AND tenant_id = ?", id, tenantID).First(&academicYear).Error
	if err != nil {
		return nil, err
	}
	return &academicYear, nil
}

func (r *repository) CreateAcademicYear(academicYear *domain.AcademicYear) error {
	return database.GetDB().Create(academicYear).Error
}

func (r *repository) UpdateAcademicYear(academicYear *domain.AcademicYear) error {
	return database.GetDB().Save(academicYear).Error
}

func (r *repository) DeleteAcademicYear(id, tenantID uint) error {
	return database.GetDB().Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&domain.AcademicYear{}).Error
}
