package domain

type ChatRepository interface {
	ListActiveRooms() ([]Room, error)
	ListRoomDetailsByRoomID(roomID string) (RoomDetails, error)
	JoinRoom(userID string, roomID string) error
	LeaveRoom(userID string) error
	SendMessage(userID, roomID, content string) error
}

type UserRepository interface {
	GetUserByUsername(username string) (*User, error)
	SaveUser(user *User) error
}
