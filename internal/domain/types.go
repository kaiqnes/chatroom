package domain

import (
	"database/sql"
	"time"
)

type User struct {
	ID        string `gorm:"default:uuid_generate_v3()"`
	Username  string
	Password  string
	Email     string
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP()"`
}

type Room struct {
	ID          string `gorm:"default:uuid_generate_v3()"`
	Name        string
	Description string
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP()"`
}

type UserRoom struct {
	ID       string `gorm:"default:uuid_generate_v3()"`
	UserID   string
	RoomID   string
	JoinedAt time.Time `gorm:"default:CURRENT_TIMESTAMP()"`
	LeftAt   *time.Time
}

func (UserRoom) TableName() string {
	return "users_rooms"
}

type Message struct {
	ID         string `gorm:"default:uuid_generate_v3()"`
	UserRoomID string
	Body       string
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP()"`
}

type UserRoomDetail struct {
	UserID   string
	Username string
	LeftAt   time.Time
}

type MessageDetail struct {
	MessageID        string
	MessageContent   string
	MessageCreatedAt time.Time
}

type RoomDetailsDB struct {
	RoomID           string       `gorm:"column:room_id"`
	RoomName         string       `gorm:"column:room_name"`
	RoomDescription  string       `gorm:"column:room_description"`
	RoomCreatedAt    time.Time    `gorm:"column:room_created_at"`
	UserID           string       `gorm:"column:user_id"`
	Username         string       `gorm:"column:username"`
	UserLeftAt       sql.NullTime `gorm:"column:user_left_at"`
	MessageID        string       `gorm:"column:message_id"`
	MessageContent   string       `gorm:"column:message_content"`
	MessageCreatedAt time.Time    `gorm:"column:message_created_at"`
}

type RoomDetails struct {
	RoomID          string
	RoomName        string
	RoomDescription string
	RoomCreatedAt   time.Time
	UsersRoomsID    string
	Users           []UserRoomDetail
	Messages        []MessageDetail
}