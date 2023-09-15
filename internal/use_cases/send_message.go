package use_cases

import (
	"chatroom/internal/domain"
	"fmt"
	"regexp"
	"time"
)

type sendMessageUseCase struct {
	chatRepository domain.ChatRepository
	userRepository domain.UserRepository
	botClient      domain.StockBotClient
	mq             domain.MessageQueue
}

func NewSendMessageUseCase(chatRepository domain.ChatRepository, userRepository domain.UserRepository,
	botClient domain.StockBotClient, mq domain.MessageQueue) domain.SendMessageUseCase {
	return &sendMessageUseCase{
		chatRepository: chatRepository,
		userRepository: userRepository,
		botClient:      botClient,
		mq:             mq,
	}
}

func (u *sendMessageUseCase) SendMessage(username, roomID, content string) ([]domain.MessageResponseDto, error) {
	response := []domain.MessageResponseDto{
		{
			RoomID:    roomID,
			Message:   content,
			Username:  username,
			Timestamp: time.Now(),
		},
	}
	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("[sendMessageUseCase.SendMessage] error getting user by username. username %s. Err: %w", username, err)
	}

	err = u.chatRepository.SendMessage(user.ID, roomID, content)
	if err != nil {
		return nil, fmt.Errorf("[sendMessageUseCase.SendMessage] error saving message. Err: %w", err)
	}

	if u.isACommand(content) {
		botMsg, err := u.SendToBotMessage(roomID, content)
		if err != nil {
			return nil, fmt.Errorf("[sendMessageUseCase.SendMessage] error sending message to bot. Err: %w", err)
		}
		response = append(response, domain.MessageResponseDto{
			RoomID:    roomID,
			Message:   botMsg,
			Username:  "StockBot",
			Timestamp: time.Now(),
		})
	}

	return response, nil
}

func (u *sendMessageUseCase) SendToBotMessage(roomID, content string) (string, error) {
	stockCode, err := u.getStockFromMessage(content)
	if err != nil {
		return "", fmt.Errorf("[sendMessageUseCase.SendToBotMessage] error getting stock code from message. Err: %w", err)
	}

	request := domain.StockBotRequest{StockCode: stockCode}

	response, err := u.botClient.Call(request)
	if err != nil {
		return "", fmt.Errorf("[sendMessageUseCase.SendToBotMessage] error calling bot. Err: %w", err)
	}

	// FIX: send to rabbitMQ
	//botRequest := domain.MessageRequestDto{
	//	Username: "stockBot",
	//	RoomID:   roomID,
	//	Message:  response.GetFormattedResponse(),
	//}
	//
	//botRequestBytes, err := json.Marshal(botRequest)
	//if err != nil {
	//	return "", err
	//}
	//
	//err = u.mq.Send(botRequestBytes)
	//if err != nil {
	//	return err
	//}

	return response.GetFormattedResponse(), nil
}

func (u *sendMessageUseCase) getStockFromMessage(content string) (string, error) {
	re := regexp.MustCompile(domain.StockRequestTemplate)
	match := re.FindStringSubmatch(content)

	if len(match) > 0 {
		return match[1], nil
	}
	return "", fmt.Errorf("[sendMessageUseCase.getStockFromMessage] error getting stock code from message. content %s", content)
}

func (u *sendMessageUseCase) isACommand(content string) bool {
	for _, command := range domain.CommandList {
		re := regexp.MustCompile(command)
		match := re.FindStringSubmatch(content)
		if len(match) > 0 {
			return true
		}
	}
	return false
}
