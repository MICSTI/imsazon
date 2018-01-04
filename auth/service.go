/*
	The auth service issues a JWT auth token for authentication inside the microservice architecture
 */

package auth

import (
	"errors"
	"github.com/MICSTI/imsazon/user"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// JWT secret - this should definitely be stored more securely
var jwtSecret = []byte("0lLXg7jzoM1a9rEx8AXqv0SMFr4OCRaChlNxgmqIxzb6OLWvh6t96oFwrnQHCvVl")

// ErrInvalidArgument is returned when one or more arguments are invalid
var ErrInvalidArgument = errors.New("Invalid argument")

// ErrExpiredJwt is returned when the JWT itself is valid but it has already expired
var ErrExpired = errors.New("Expired JWT")

// ErrInvalid is returned when the JWT is not valid
var ErrInvalid = errors.New("Invalid JWT")

// Service is the interface that provides the methods for obtaining an auth token
type Service interface {
	// Login checks the passed credentials and issues a JWT auth token in case they are valid
	Login(username string, password string) (string, error)

	// Check checks if the passed JWT auth token is valid
	Check(token string) (user.UserId, error)
}

type service struct {
	users		user.Repository
}

func (s *service) Login(username string, password string) (string, error) {
	if username == "" || password == "" {
		return "", ErrInvalidArgument
	}

	u, err := s.users.CheckLogin(username, password)
	if err != nil {
		return "", err
	}

	// create the JWT auth token for the user
	token := jwt.New(jwt.SigningMethodHS256)

	// create a map to store the claims
	claims := token.Claims.(jwt.MapClaims)

	// set token claims
	claims["sub"] = u.Id
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	claims["role"] = u.Role.String()

	// sign the token with the secret
	signedToken, _ := token.SignedString(jwtSecret)

	return signedToken, nil
}

func (s *service) Check(tokenString string) (user.UserId, error) {
	if tokenString == "" {
		return "", ErrInvalidArgument
	}

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalid
		}

		return jwtSecret, nil
	})

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		//return claims["sub"].(user.UserId), nil
		return user.U0003, nil
	} else {
		return "", ErrExpired
	}
}

// NewService returns a new instance of the auth service
func NewService(users user.Repository) Service {
	return &service{
		users:		users,
	}
}