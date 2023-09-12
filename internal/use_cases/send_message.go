package use_cases

import (
	"chatroom/internal/domain"
	"chatroom/internal/logger"
)

type sendMessageUseCase struct {
	chatRepository domain.ChatRepository
	userRepository domain.UserRepository
	log            *logger.CustomLogger
}

func NewSendMessageUseCase(chatRepository domain.ChatRepository, userRepository domain.UserRepository, log *logger.CustomLogger) domain.SendMessageUseCase {
	return &sendMessageUseCase{
		chatRepository: chatRepository,
		userRepository: userRepository,
		log:            log,
	}
}

func (u *sendMessageUseCase) SendMessage(userID, roomID, content string) error {
	err := u.chatRepository.IsUserInRoom(userID, roomID)
	if err != nil {
		return err
	}

	err = u.chatRepository.SendMessage(userID, roomID, content)
	if err != nil {
		return err
	}

	return nil
}
