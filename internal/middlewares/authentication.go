package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"chatroom/internal/config"
	"chatroom/internal/logger"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type authenticationMiddleware struct {
	log logger.CustomLogger
	cfg config.Config
}

type AuthenticationMiddleware interface {
	SetToken(ctx *gin.Context, username string, expirationTime time.Time) error
	ValidateToken() gin.HandlerFunc
}

func NewAuthenticationMiddleware(log logger.CustomLogger) AuthenticationMiddleware {
	return &authenticationMiddleware{log: log}
}

func (m *authenticationMiddleware) SetToken(ctx *gin.Context, username string, expirationTime time.Time) error {
	claim := &claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString([]byte(m.cfg.JwtKey))
	if err != nil {
		return err
	}

	ctx.SetCookie("token", tokenString, expirationTime.Second(), "/", "localhost", false, true)
	return nil
}

func (m *authenticationMiddleware) ValidateToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claim, status, err := m.validateToken(ctx)
		if err != nil {
			ctx.JSON(status, gin.H{
				"err": err.Error(),
			})
			ctx.Abort()
			return
		}

		if time.Until(claim.ExpiresAt.Time) > 30*time.Second {
			ctx.Next()
			return
		}

		expirationTime := time.Now().Add(5 * time.Minute)
		claim.ExpiresAt = jwt.NewNumericDate(expirationTime)
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
		tokenString, err := token.SignedString([]byte(m.cfg.JwtKey))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"err": err.Error(),
			})
			ctx.Abort()
			return
		}

		ctx.SetCookie("token", tokenString, expirationTime.Second(), "/", "localhost", false, true)

		ctx.Next()
	}
}

func (m *authenticationMiddleware) validateToken(ctx *gin.Context) (*claims, int, error) {
	var claim claims

	tknStr, status, err := m.hasToken(ctx)
	if err != nil {
		return nil, status, err
	}

	tkn, err := jwt.ParseWithClaims(tknStr, &claim, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.cfg.JwtKey), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, http.StatusUnauthorized, fmt.Errorf("signature invalid 1")
		}
		return nil, http.StatusBadRequest, fmt.Errorf("signature invalid 2: err: %v", err)
	}
	if !tkn.Valid {
		return nil, http.StatusUnauthorized, fmt.Errorf("signature invalid 3")
	}
	return &claim, http.StatusOK, nil
}

func (m *authenticationMiddleware) hasToken(ctx *gin.Context) (string, int, error) {
	tknStr, err := ctx.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			return "", http.StatusUnauthorized, fmt.Errorf("cookie not found")
		}
		return "", http.StatusBadRequest, fmt.Errorf("cookie not found")
	}
	return tknStr, http.StatusOK, nil
}
