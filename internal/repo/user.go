package repo

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"matheus.com/vgs/internal/model"
)

type UserRepo interface {
	Create(user model.User) error
	GetById(id uuid.UUID) (model.User, error)
	GetAll() ([]model.User, error)
}

type userRepo struct {
	*gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepo {
	db.AutoMigrate(&model.User{})
	return &userRepo{DB: db}
}

func (r *userRepo) GetById(id uuid.UUID) (model.User, error) {
	u := model.User{}
	if r.DB.Where("id = ?", id).Find(&u).RowsAffected == 0 {
		return u, gorm.ErrRecordNotFound
	}
	return u, nil
}

// gorm does not export a unique violation error
// detecting it here requires assertions and code comparisons
func (r *userRepo) Create(user model.User) error {
	dbResult := r.DB.Create(&user)
	return dbResult.Error
}

func (r *userRepo) GetAll() ([]model.User, error) {
	var users []model.User
	dbResult := r.DB.Find(&users)
	return users, dbResult.Error
}
