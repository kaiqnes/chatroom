package controllers

import (
	"chatroom/internal/domain"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type authController struct {
	routes  *gin.Engine
	useCase domain.AuthUseCase
}

func NewAuthController(routes *gin.Engine, authUseCase domain.AuthUseCase) domain.AuthController {
	return &authController{routes: routes, useCase: authUseCase}
}

func (c *authController) SetupEndpoints() {
	c.routes.POST("/signin", c.SignIn)
	c.routes.POST("/signup", c.SignUp)
	c.routes.GET("/signout", c.SignOut)
}

func (c *authController) SignIn(ctx *gin.Context) {
	var req domain.SignRequestDto
	err := ctx.BindJSON(&req)
	if err != nil {
		fmt.Printf("error binding json: %v\n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}

	errs := c.validateRequest(req)
	if len(errs) > 0 {
		fmt.Printf("missing mandatory fields: %v\n", errs)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": errs,
		})
		return
	}

	token, expiration, err := c.useCase.SignIn(ctx, req.Username, req.Password)
	if err != nil {
		fmt.Printf("error signing in: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}

	ctx.SetCookie("token", token, expiration, "/", "localhost", false, true)

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (c *authController) SignUp(ctx *gin.Context) {
	var creds domain.SignRequestDto
	err := ctx.BindJSON(&creds)
	if err != nil {
		fmt.Printf("error binding json: %v\n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}

	errs := c.validateRequest(creds)
	if len(errs) > 0 {
		fmt.Printf("missing mandatory fields: %v\n", errs)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": errs,
		})
		return
	}

	err = c.useCase.SignUp(ctx, creds.Username, creds.Password)
	if err != nil {
		fmt.Printf("error signing up: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"response": "signed up",
	})
}

func (c *authController) validateRequest(creds domain.SignRequestDto) []string {
	var errs []string
	if len(creds.Username) == 0 {
		errs = append(errs, "empty username")
	}
	if len(creds.Password) == 0 {
		errs = append(errs, "empty password")
	}
	return errs
}

func (c *authController) SignOut(ctx *gin.Context) {
	ctx.SetCookie("token", "", -1, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, gin.H{
		"response": "Logged out",
	})
}
