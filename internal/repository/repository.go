package repository

import (
	"cursebot/internal/models"
	"fmt"

	"gorm.io/gorm"
)

type Repository interface {
	CreateUser(user *models.User) error
	GetUserByTelegramID(telegramID int64) (*models.User, error)
	UpdateUser(user *models.User) error
	GetCourse(id uint) (*models.Course, error)
	ListCourses() ([]models.Course, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateUser(user *models.User) error {
	return r.db.Session(&gorm.Session{PrepareStmt: false}).Create(user).Error
}

func (r *repository) GetUserByTelegramID(telegramID int64) (*models.User, error) {
	var user models.User
	err := r.db.Session(&gorm.Session{PrepareStmt: false}).
		Where("telegram_id = ?", telegramID).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) UpdateUser(user *models.User) error {
	return r.db.Session(&gorm.Session{PrepareStmt: false}).Save(user).Error
}

func (r *repository) GetCourse(id uint) (*models.Course, error) {
	var course models.Course
	err := r.db.Session(&gorm.Session{PrepareStmt: false}).
		First(&course, id).Error
	if err != nil {
		return nil, err
	}
	return &course, nil
}

func (r *repository) ListCourses() ([]models.Course, error) {
	var courses []models.Course
	err := r.db.Session(&gorm.Session{PrepareStmt: false}).
		Table("courses").
		Select("*").
		Where("deleted_at IS NULL").
		Find(&courses).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list courses: %w", err)
	}

	return courses, nil
}
