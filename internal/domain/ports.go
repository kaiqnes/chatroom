package domain

import "context"

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
	SendMessage(username, roomID, content string) error
	SendToBotMessage(roomID, content string) error
	SendFromBotMessage(roomID, content string) error
}

type AuthUseCase interface {
	SignIn(ctx context.Context, username, password string) (string, int, error)
	SignUp(ctx context.Context, username, password string) error
}

type StockBotClient interface {
	Call(req StockBotRequest) (*StockBotResponse, error)
}
