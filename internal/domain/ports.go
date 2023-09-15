package domain

import (
	"context"

	socketio "github.com/googollee/go-socket.io"

	"github.com/gin-gonic/gin"
)

type ChatRepository interface {
	ListActiveRooms() ([]Room, error)
	ListRoomDetailsByRoomID(roomID string) (*RoomDetails, error)
	IsUserInRoom(userID, roomID string) error
	JoinRoom(userID string, roomID string) error
	LeaveRoom(userID string) error
	SendMessage(userID, roomID, content string) error
}

type UserRepository interface {
	GetUserByUsername(username string) (*User, error)
	GetUserByID(userID string) (*User, error)
	SaveUser(user *User) error
}

type SendMessageUseCase interface {
	SendMessage(username, roomID, content string) ([]MessageResponseDto, error)
	SendToBotMessage(roomID, content string) (string, error)
}

type AuthUseCase interface {
	SignIn(ctx context.Context, username, password string) (string, int, error)
	SignUp(ctx context.Context, username, password string) error
}

type StockBotClient interface {
	Call(req StockBotRequest) (*StockBotResponseDto, error)
}

type MessageQueue interface {
	Send(message []byte) error
	Listen() error
}

type AuthController interface {
	SetupEndpoints()
	SignIn(ctx *gin.Context)
	SignUp(ctx *gin.Context)
}

type ChatController interface {
	SetupEndpoints()
	SendMessageFromBot(req MessageRequestDto)
	OnMessage(socket socketio.Conn, req MessageRequestDto)
	ListMessages(ctx *gin.Context)
}

type HealthController interface {
	SetupEndpoints()
	HealthCheck(ctx *gin.Context)
}

type AuthenticationMiddleware interface {
	ValidateToken() gin.HandlerFunc
}
