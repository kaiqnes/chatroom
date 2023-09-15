package domain

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const StockRequestTemplate = `^\\/stock=(.)$`

var CommandList = []string{StockRequestTemplate}

type ClaimsDto struct {
	Username string `json:"username"`
	UserID   string `json:"user_id"`
	jwt.RegisteredClaims
}

type User struct {
	ID        string `gorm:"default:uuid_generate_v3()"`
	Username  string
	Password  string
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
}

func (UserRoom) TableName() string {
	return "users_rooms"
}

type Message struct {
	ID        string `gorm:"default:uuid_generate_v3()"`
	UserID    string
	RoomID    string
	Body      string
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP()"`
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

type StockBotRequest struct {
	StockCode string
}

type StockBotResponseDto struct {
	Symbol string `csv:"Symbol"`
	Date   string `csv:"Date"`
	Time   string `csv:"Time"`
	Open   string `csv:"Open"`
	High   string `csv:"High"`
	Low    string `csv:"Low"`
	Close  string `csv:"Close"`
	Volume string `csv:"Volume"`
}

func (s *StockBotResponseDto) GetFormattedResponse() string {
	return fmt.Sprintf("%s quote is $%s per share", s.Symbol, s.Close)
}

type SignRequestDto struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type MessageRequestDto struct {
	Username string `json:"username"`
	RoomID   string `json:"room_id"`
	Message  string `json:"message"`
}

type MessagesResponseDto struct {
	Messages []MessageResponseDto `json:"messages"`
}

type MessageResponseDto struct {
	RoomID    string    `json:"room_id"`
	Message   string    `json:"message"`
	Username  string    `json:"username"`
	Timestamp time.Time `json:"timestamp"`
}

type Claims struct {
	Username string `json:"username"`
	UserID   string `json:"user_id"`
	jwt.RegisteredClaims
}
