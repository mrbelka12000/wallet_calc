package v1

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/mrbelka12000/wallet_calc/internal/dto"
)

var jwtSecretKey []byte

type Token struct {
	JWT string `json:"jwt"`
}

func (c *Controller) buildJWT(user dto.User) (string, error) {
	payload := jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	t, err := token.SignedString(jwtSecretKey)
	if err != nil {
		c.log.With("error", err).Error("error signing token")
		return "", err
	}

	return t, nil
}

func (c *Controller) jwtPayloadFromRequest(tokenString string) (jwt.MapClaims, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})

	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	payload, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, err
	}

	return payload, nil
}

func castClaims(claims map[string]interface{}) (dto.User, error) {
	id, ok := claims["id"].(uuid.UUID)
	if !ok {
		return dto.User{}, errors.New("invalid token")
	}

	return dto.User{
		ID: id,
	}, nil
}
