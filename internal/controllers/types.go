package controllers

import "time"

type Controller interface {
	SetupEndpoints()
}

type signInputDto struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type messageInputDto struct {
	Username string `json:"username"`
	RoomID   string `json:"room_id"`
	Message  string `json:"message"`
}

type messagesOutputDto struct {
	Messages []MessageOutputDto `json:"messages"`
}

type MessageOutputDto struct {
	RoomID    string    `json:"room_id"`
	Message   string    `json:"message"`
	Username  string    `json:"username"`
	Timestamp time.Time `json:"timestamp"`
}
