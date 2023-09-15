package repositories

import (
	"chatroom/internal/db"
	"chatroom/internal/domain"
	"fmt"
)

type userRepository struct {
	db *db.Database
}

func NewUserRepository(db *db.Database) domain.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) GetUserByUsername(username string) (*domain.User, error) {
	var user domain.User
	tx := r.db.DBInstance.Where("username = ?", username).First(&user)
	if tx.Error != nil {
		return nil, fmt.Errorf("[userRepository.GetUserByUsername] error getting user by username. username %s. Err: %w", username, tx.Error)
	}
	return &user, nil
}

func (r *userRepository) GetUserByID(userID string) (*domain.User, error) {
	var user domain.User
	tx := r.db.DBInstance.Where("id = ?", userID).First(&user)
	if tx.Error != nil {
		return nil, fmt.Errorf("[userRepository.GetUserByID] error getting user by id. userID %s. Err: %w", userID, tx.Error)
	}
	return &user, nil
}

func (r *userRepository) SaveUser(user *domain.User) error {
	tx := r.db.DBInstance.Select("username", "password").Create(&user)
	if tx.Error != nil {
		return fmt.Errorf("[userRepository.SaveUser] error saving user. Err: %w", tx.Error)
	}
	return nil
}
