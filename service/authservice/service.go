package authservice

import (
	"errors"
	"strings"
	"telegram-service-platform/entity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	config Config
}

func New(config Config) *Service {
	return &Service{config: config}
}

func (s *Service) CreateAccessToken(user entity.User) (string, error) {
	return createToken(user.TelegramID, int64(user.ID), user.Role.String(), s.config.SignKey, s.config.AccessSubject, s.config.AccessTokenDuration)
}

func (s *Service) CreateRefreshToken(user entity.User) (string, error) {
	return createToken(user.TelegramID, int64(user.ID), user.Role.String(), s.config.SignKey, s.config.RefreshSubject, s.config.RefreshTokenDuration)
}

func createToken(telegramId int64, userID int64, role string, signKey, subject string, duration time.Duration) (string, error) {

	t := jwt.New(jwt.SigningMethodHS256)

	t.Claims = &Claims{
		jwt.RegisteredClaims{
			Subject:   subject,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		telegramId,
		userID,
		role,
	}

	return t.SignedString([]byte(signKey))
}

func (s *Service) ParseToken(tokenStr string) (*Claims, error) {
	const op = "authservice.ParseToken"
	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
	if tokenStr == "" {
		return nil, richerror.New(op, errors.New(msgerror.ErrMissingBearer))
	}
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

			return nil, richerror.New(op, errors.New(msgerror.ErrInvalidAlgorithm))
		}
		return []byte(s.config.SignKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, richerror.New(op, errors.New(msgerror.ErrTokenExpired))
		}
		return nil, richerror.New(op, errors.New(msgerror.ErrInvalidToken))
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, richerror.New(op, errors.New(msgerror.ErrInvalidToken))
	}

	return claims, nil
}
