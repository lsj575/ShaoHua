package repositories

import (
	"gorm.io/gorm"
	"srvs/user_srv/models"
)

var UserRepository = newUserRepository()

func newUserRepository() *userRepository {
	return &userRepository{}
}

type userRepository struct {
}

func (r *userRepository) GetUserById(db *gorm.DB, id uint64) *models.User {
	ret := &models.User{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *userRepository) Take(db *gorm.DB, where ...interface{}) *models.User {
	ret := &models.User{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *userRepository) Create(db *gorm.DB, u *models.User) error {
	err := db.Create(u).Error
	return err
}

func (r *userRepository) Update(db *gorm.DB, t *models.User) (err error) {
	err = db.Save(t).Error
	return
}

func (r *userRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&models.User{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *userRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&models.User{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *userRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.User{}, "id = ?", id)
}

func (r *userRepository) GetByEmail(db *gorm.DB, email string) *models.User {
	return r.Take(db, "email = ?", email)
}

func (r *userRepository) GetByUsername(db *gorm.DB, username string) *models.User {
	return r.Take(db, "username = ?", username)
}
