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

// create a custom JWT claims struct
type CustomClaims struct {
	Role string `json:"role"`
	jwt.StandardClaims
}

func (s *service) Login(username string, password string) (string, error) {
	if username == "" || password == "" {
		return "", ErrInvalidArgument
	}

	u, err := s.users.CheckLogin(username, password)
	if err != nil {
		return "", err
	}

	// create the claims
	claims := CustomClaims{
		u.Role.String(),
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer: "imsazon",
			Subject: u.Id.String(),
		},
	}

	// create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString(jwtSecret)

	return signedToken, nil
}

func (s *service) Check(tokenString string) (user.UserId, error) {
	if tokenString == "" {
		return "", ErrInvalidArgument
	}

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	// for the case we just received a random string
	if token == nil {
		return "", ErrInvalid
	}

	if token.Valid {
		if claims, ok := token.Claims.(*CustomClaims); ok {
			userId := user.UserId(claims.Subject)

			return userId, nil
		} else {
			return "", ErrInvalid
		}
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors & jwt.ValidationErrorMalformed != 0 {
			return "", ErrInvalid
		} else if ve.Errors & (jwt.ValidationErrorExpired | jwt.ValidationErrorNotValidYet) != 0 {
			return "", ErrExpired
		} else {
			return "", ErrInvalid
		}
	}

	return "", ErrInvalid
}

// NewService returns a new instance of the auth service
func NewService(users user.Repository) Service {
	return &service{
		users:		users,
	}
}