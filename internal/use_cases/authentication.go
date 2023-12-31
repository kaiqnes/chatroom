package use_cases

import (
	"chatroom/internal/domain"
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type authUseCase struct {
	jwtSecret      []byte
	userRepository domain.UserRepository
}

func NewAuthUseCase(secret string, userRepository domain.UserRepository) domain.AuthUseCase {
	return &authUseCase{
		jwtSecret:      []byte(secret),
		userRepository: userRepository,
	}
}

func (u *authUseCase) SignIn(ctx context.Context, username, password string) (string, int, error) {
	// check if user exists
	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return "", 0, fmt.Errorf("[authUseCase.SignIn] error getting user: %w", err)
	}

	// compare password
	err = u.compareHash(password, user.Password)
	if err != nil {
		return "", 0, fmt.Errorf("[authUseCase.SignIn] error comparing hash: %w", err)
	}

	// generate token
	expirationTime := time.Now().Add(5 * time.Minute)
	tokenString, err := u.genToken(user, expirationTime)
	if err != nil {
		return "", 0, fmt.Errorf("[authUseCase.SignIn] error signing token: %w", err)
	}

	return tokenString, int(expirationTime.Unix()), nil
}

func (u *authUseCase) SignUp(ctx context.Context, username, password string) error {
	// check if user exists
	_, err := u.userRepository.GetUserByUsername(username)
	if err == nil {
		return fmt.Errorf("[authUseCase.SignUp] user already exists")
	}

	// hash password
	hashedPassword, err := u.hashAndSalt(password)
	if err != nil {
		return fmt.Errorf("[authUseCase.SignUp] error hashing password: %w", err)
	}

	// save user
	user := &domain.User{
		Username: username,
		Password: hashedPassword,
	}
	err = u.userRepository.SaveUser(user)
	if err != nil {
		return fmt.Errorf("[authUseCase.SignUp] error saving user: %w", err)
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

	return token.SignedString(u.jwtSecret)
}

func (u *authUseCase) compareHash(received, stored string) error {
	return bcrypt.CompareHashAndPassword([]byte(stored), []byte(received))
}

func (u *authUseCase) hashAndSalt(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("[authUseCase.hashAndSalt] error generating hash: %w", err)
	}
	return string(hash), nil
}
