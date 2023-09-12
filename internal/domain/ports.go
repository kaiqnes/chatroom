package domain

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
	SendMessage(userID, roomID, content string) error
}
