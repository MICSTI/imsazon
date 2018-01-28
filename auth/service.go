/*
	The auth service issues a JWT auth token for authentication inside the microservice architecture
 */

package auth

import (
	"errors"
	userModel "github.com/MICSTI/imsazon/models/user"
	"github.com/dgrijalva/jwt-go"
	"time"
)

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
	Check(token string) (userModel.UserId, error)
}

type service struct {
	jwtSecret	[]byte
	users		userModel.Repository
}

// create a custom JWT claims struct
type CustomClaims struct {
	Role 					string	 	`json:"role"`
	Name					string		`json:"name,omitempty"`
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
		u.Name,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer: "imsazon",
			Subject: u.Id.String(),
		},
	}

	// create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString(s.jwtSecret)

	return signedToken, nil
}

func (s *service) Check(tokenString string) (userModel.UserId, error) {
	if tokenString == "" {
		return "", ErrInvalidArgument
	}

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	// for the case we just received a random string
	if token == nil {
		return "", ErrInvalid
	}

	if token.Valid {
		if claims, ok := token.Claims.(*CustomClaims); ok {
			userId := userModel.UserId(claims.Subject)

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
func NewService(jwtSecret []byte, users userModel.Repository) Service {
	return &service{
		jwtSecret:	jwtSecret,
		users:		users,
	}
}