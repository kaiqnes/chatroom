package controllers

import (
	"net/http"
	"time"

	"chatroom/internal/db"
	"chatroom/internal/logger"
	"chatroom/internal/middlewares"
	"github.com/gin-gonic/gin"
)

var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

type authController struct {
	routes         *gin.Engine
	db             *db.Database
	authMiddleware middlewares.AuthenticationMiddleware
	log            logger.CustomLogger
}

func NewAuthController(routes *gin.Engine, authMiddleware middlewares.AuthenticationMiddleware) Controller {
	return &authController{routes: routes, authMiddleware: authMiddleware}
}

func (c *authController) SetupEndpoints() {
	c.routes.POST("/signin", c.SignIn)
	c.routes.POST("/signup", c.SignUp)
	c.routes.GET("/signout", c.SignOut)
}

func (c *authController) SignIn(ctx *gin.Context) {
	var creds signInputDto
	err := ctx.BindJSON(&creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}

	// TODO: move it to be executed in DB
	expectedPassword, ok := users[creds.Username]

	if !ok || expectedPassword != creds.Password {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"err": err.Error(),
		})
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)

	err = c.authMiddleware.SetToken(ctx, creds.Username, expirationTime)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"response": "signed in",
	})
}

func (c *authController) SignUp(ctx *gin.Context) {
	var creds signInputDto
	err := ctx.BindJSON(&creds)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}

	var errs []string
	if len(creds.Username) == 0 {
		errs = append(errs, "empty username")
	}
	if len(creds.Password) == 0 {
		errs = append(errs, "empty password")
	}
	if len(errs) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": errs,
		})
		return
	}

	// TODO: move it to be executed in DB
	if _, ok := users[creds.Username]; ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": "username already exists",
		})
		return
	}

	// TODO: move it to be executed in DB
	users[creds.Username] = creds.Password

	ctx.JSON(http.StatusOK, gin.H{
		"response": "signed up",
	})
	return
}

func (c *authController) SignOut(ctx *gin.Context) {
	ctx.SetCookie("token", "", -1, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, gin.H{
		"response": "Logged out",
	})
}
