package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"chatroom/internal/config"
	"chatroom/internal/domain"
	"chatroom/internal/logger"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type authenticationMiddleware struct {
	log logger.CustomLogger
	cfg *config.Config
}

type AuthenticationMiddleware interface {
	ValidateToken() gin.HandlerFunc
}

func NewAuthenticationMiddleware(cfg *config.Config, log logger.CustomLogger) AuthenticationMiddleware {
	return &authenticationMiddleware{cfg: cfg, log: log}
}

func (m *authenticationMiddleware) ValidateToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, status, err := m.validateToken(ctx)
		if err != nil {
			fmt.Printf("redirecting with status %d due err: %+v", status, err.Error())
			ctx.Redirect(http.StatusTemporaryRedirect, "http://localhost:8080/sign")
			ctx.Abort()
			return
		}

		if time.Until(claims.ExpiresAt.Time) > 30*time.Second {
			ctx.Next()
			return
		}

		expirationTime := time.Now().Add(5 * time.Minute)
		claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
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

func (m *authenticationMiddleware) validateToken(ctx *gin.Context) (*domain.Claims, int, error) {
	var claims domain.Claims

	tknStr, status, err := m.hasToken(ctx)
	if err != nil {
		return nil, status, err
	}

	tkn, err := jwt.ParseWithClaims(tknStr, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.cfg.JwtKey), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, http.StatusUnauthorized, fmt.Errorf("signature invalid 1")
		}
		return nil, http.StatusBadRequest, fmt.Errorf("signature invalid 2: err: %w", err)
	}
	if !tkn.Valid {
		return nil, http.StatusUnauthorized, fmt.Errorf("signature invalid 3")
	}
	return &claims, http.StatusOK, nil
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
