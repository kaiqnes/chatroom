package repositories

import (
	"time"

	"chatroom/internal/db"
	"chatroom/internal/domain"
	"chatroom/internal/logger"
)

type chatRepository struct {
	db  *db.Database
	log logger.CustomLogger
}

func NewChatRepository(db *db.Database, log *logger.CustomLogger) domain.ChatRepository {
	return &chatRepository{
		db:  db,
		log: log,
	}
}

func (r *chatRepository) IsUserInRoom(userID, roomID string) error {
	var userRoom domain.UserRoom
	tx := r.db.DBInstance.Where("user_id = ? AND room_id = ? AND left_at IS NULL", userID, roomID).First(&userRoom)
	if tx.Error != nil {
		return tx.Error
	}

	if userRoom.ID == "" {
		return domain.ErrUserNotInRoom
	}

	return nil
}

func (r *chatRepository) SendMessage(userID, roomID, content string) error {
	tx := r.db.DBInstance.Create(&domain.Message{
		UserID: userID,
		RoomID: roomID,
		Body:   content,
	})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (r *chatRepository) ListActiveRooms() ([]domain.Room, error) {
	var rooms []domain.Room
	tx := r.db.DBInstance.Where("deactivated_at IS NULL").Find(&rooms)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return rooms, nil
}

func (r *chatRepository) ListRoomDetailsByRoomID(roomID string) (*domain.RoomDetails, error) {
	var roomDetails []domain.RoomDetailsDB
	tx := r.db.DBInstance.Table("rooms r").
		Select(`r.id as room_id, r.name as room_name, r.description as room_description, 
r.created_at as room_created_at, u.id as user_id, u.username as username, ur.left_at as user_left_at, 
m.id as message_id, m.body as message_content, m.created_at as message_created_at`).
		Joins("LEFT JOIN users_rooms ur ON r.id = ur.room_id").
		Joins("LEFT JOIN users u ON ur.user_id = u.id").
		Joins("LEFT JOIN messages m ON ur.id = m.user_room_id").
		Where("r.deactivated_at IS NULL AND r.id = ?", roomID).
		Order("m.created_at DESC").
		Limit(50).
		Scan(&roomDetails)
	if tx.Error != nil {
		return nil, tx.Error
	}

	var results domain.RoomDetails
	var addedUsers = make(map[string]domain.UserRoomDetail)
	for i, row := range roomDetails {
		if i == 0 { // on first row, get room details
			results.RoomID = row.RoomID
			results.RoomName = row.RoomName
			results.RoomDescription = row.RoomDescription
			results.RoomCreatedAt = row.RoomCreatedAt
		}

		if row.MessageID != "" {
			var md domain.MessageDetail
			md.MessageID = row.MessageID
			md.MessageContent = row.MessageContent
			md.MessageCreatedAt = row.MessageCreatedAt
			results.Messages = append(results.Messages, md)
		}

		if urd, ok := addedUsers[row.UserID]; !ok {
			addedUsers[row.UserID] = urd
			if row.UserLeftAt.Valid {
				continue
			}
			urd.UserID = row.UserID
			urd.Username = row.Username
			results.Users = append(results.Users, urd)
		}
	}

	return &results, nil
}

func (r *chatRepository) JoinRoom(userID string, roomID string) error {
	tx := r.db.DBInstance.Create(&domain.UserRoom{UserID: userID, RoomID: roomID})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (r *chatRepository) LeaveRoom(userID string) error {
	tx := r.db.DBInstance.Table("users_rooms").
		Where("user_id = ? AND left_at IS NULL", userID).
		UpdateColumn("left_at", time.Now().UTC())
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
