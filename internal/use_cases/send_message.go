package use_cases

import (
	"chatroom/internal/domain"
	"chatroom/internal/logger"
	"encoding/json"
	"fmt"
	"regexp"
)

type sendMessageUseCase struct {
	chatRepository domain.ChatRepository
	userRepository domain.UserRepository
	botClient      domain.StockBotClient
	mq             domain.MessageQueue
	log            logger.CustomLogger
}

func NewSendMessageUseCase(chatRepository domain.ChatRepository, userRepository domain.UserRepository,
	botClient domain.StockBotClient, log logger.CustomLogger) domain.SendMessageUseCase {
	return &sendMessageUseCase{
		chatRepository: chatRepository,
		userRepository: userRepository,
		botClient:      botClient,
		log:            log,
	}
}

func (u *sendMessageUseCase) SendMessage(username, roomID, content string) error {
	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return err
	}

	err = u.chatRepository.SendMessage(user.ID, roomID, content)
	if err != nil {
		return err
	}

	if u.isACommand(content) {
		err = u.SendToBotMessage(roomID, content)
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *sendMessageUseCase) SendToBotMessage(roomID, content string) error {
	stockCode, err := u.getStockFromMessage(content)
	if err != nil {
		return err
	}

	request := domain.StockBotRequest{StockCode: stockCode}

	response, err := u.botClient.Call(request)
	if err != nil {
		return err
	}

	botRequest := domain.MessageRequestDto{
		Username: "stockBot",
		RoomID:   roomID,
		Message:  response.GetFormattedResponse(),
	}

	botRequestBytes, err := json.Marshal(botRequest)
	if err != nil {
		return err
	}

	err = u.mq.Send(botRequestBytes)
	if err != nil {
		return err
	}

	return nil
}

func (u *sendMessageUseCase) SendFromBotMessage(roomID, content string) error {

	return nil
}

func (u *sendMessageUseCase) getStockFromMessage(content string) (string, error) {
	re := regexp.MustCompile(domain.StockRequestTemplate)
	match := re.FindStringSubmatch(content)

	if len(match) > 0 {
		return match[1], nil
	}
	return "", fmt.Errorf("invalid stock code")
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
