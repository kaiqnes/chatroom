package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"chatroom/internal/config"
	"chatroom/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type authenticationMiddleware struct {
	cfg *config.Config
}

func NewAuthenticationMiddleware(cfg *config.Config) domain.AuthenticationMiddleware {
	return &authenticationMiddleware{cfg: cfg}
}

func (m *authenticationMiddleware) ValidateToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, status, err := m.validateToken(ctx)
		if err != nil {
			fmt.Printf("[AuthenticationMiddleware.ValidateToken] redirecting with status %d due err: %+v", status, err.Error())
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

func (m *authenticationMiddleware) validateToken(ctx *gin.Context) (*domain.ClaimsDto, int, error) {
	var claims domain.ClaimsDto

	tknStr, status, err := m.hasToken(ctx)
	if err != nil {
		return nil, status, fmt.Errorf("[AuthenticationMiddleware.validateToken] hasToken: %w", err)
	}

	tkn, err := jwt.ParseWithClaims(tknStr, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.cfg.JwtKey), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, http.StatusUnauthorized, fmt.Errorf("[AuthenticationMiddleware.validateToken] signature invalid")
		}
		return nil, http.StatusBadRequest, fmt.Errorf("[AuthenticationMiddleware.validateToken] parsing token invalid: err: %w", err)
	}
	if !tkn.Valid {
		return nil, http.StatusUnauthorized, fmt.Errorf("[AuthenticationMiddleware.validateToken] token not valid")
	}
	return &claims, http.StatusOK, nil
}

func (m *authenticationMiddleware) hasToken(ctx *gin.Context) (string, int, error) {
	tknStr, err := ctx.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			return "", http.StatusUnauthorized, fmt.Errorf("[AuthenticationMiddleware.hasToken] token not found")
		}
		return "", http.StatusBadRequest, fmt.Errorf("[AuthenticationMiddleware.hasToken] failed to get cookie from ctx: %w", err)
	}
	return tknStr, http.StatusOK, nil
}
