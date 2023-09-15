package repositories

import (
	"chatroom/internal/db"
	"chatroom/internal/domain"
	"chatroom/internal/logger"
)

type userRepository struct {
	db  *db.Database
	log logger.CustomLogger
}

func NewUserRepository(db *db.Database, log logger.CustomLogger) domain.UserRepository {
	return &userRepository{
		db:  db,
		log: log,
	}
}

func (r *userRepository) GetUserByUsername(username string) (*domain.User, error) {
	var user domain.User
	tx := r.db.DBInstance.Where("username = ?", username).First(&user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &user, nil
}

func (r *userRepository) GetUserByID(userID string) (*domain.User, error) {
	var user domain.User
	tx := r.db.DBInstance.Where("id = ?", userID).First(&user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &user, nil
}

func (r *userRepository) SaveUser(user *domain.User) error {
	tx := r.db.DBInstance.Select("username", "password").Create(&user)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
