package controllers

type Controller interface {
	SetupEndpoints()
}

type signInputDto struct {
	Password string `json:"password"`
	Username string `json:"username"`
}
