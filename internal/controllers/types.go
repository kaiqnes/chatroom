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
	Message string `json:"message"`
}

type messagesOutputDto struct {
	Messages []MessageOutputDto `json:"messages"`
}

type MessageOutputDto struct {
	Message   string    `json:"message"`
	User      string    `json:"user"`
	Timestamp time.Time `json:"timestamp"`
}
