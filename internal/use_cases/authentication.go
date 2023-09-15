package use_cases

import (
	"chatroom/internal/domain"
	"chatroom/internal/logger"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type authUseCase struct {
	jwtSecret      []byte
	userRepository domain.UserRepository
	log            logger.CustomLogger
}

func NewAuthUseCase(secret string, userRepository domain.UserRepository, log logger.CustomLogger) domain.AuthUseCase {
	return &authUseCase{
		jwtSecret:      []byte(secret),
		userRepository: userRepository,
		log:            log,
	}
}

func (u *authUseCase) SignIn(ctx context.Context, username, password string) (string, int, error) {
	// check if user exists
	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return "", 0, fmt.Errorf("error getting user: %w", err)
	}

	// compare password
	err = u.compareHash(password, user.Password)
	if err != nil {
		return "", 0, fmt.Errorf("error comparing hash: %w", err)
	}

	// generate token
	expirationTime := time.Now().Add(5 * time.Minute)
	tokenString, err := u.genToken(user, expirationTime)
	if err != nil {
		return "", 0, fmt.Errorf("error signing token: %w", err)
	}

	return tokenString, int(expirationTime.Unix()), nil
}

func (u *authUseCase) SignUp(ctx context.Context, username, password string) error {
	// check if user exists
	_, err := u.userRepository.GetUserByUsername(username)
	if err == nil {
		return fmt.Errorf("user already exists")
	}

	// hash password
	hashedPassword, err := u.hashAndSalt(password)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	// save user
	user := &domain.User{
		Username: username,
		Password: hashedPassword,
	}
	err = u.userRepository.SaveUser(user)
	if err != nil {
		return fmt.Errorf("error saving user: %w", err)
	}

	return nil
}

func (u *authUseCase) genToken(user *domain.User, expirationTime time.Time) (string, error) {
	claims := &domain.Claims{
		Username: user.Username,
		UserID:   user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(u.jwtSecret)
	return tokenString, err
}

func (u *authUseCase) compareHash(received, stored string) error {
	return bcrypt.CompareHashAndPassword([]byte(stored), []byte(received))
}

func (u *authUseCase) hashAndSalt(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error generating hash: %w", err)
	}
	return string(hash), nil
}
