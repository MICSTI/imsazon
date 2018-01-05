/*
	The auth service issues a JWT auth token for authentication inside the microservice architecture
 */

package auth

import (
	"errors"
	"github.com/MICSTI/imsazon/user"
	"github.com/dgrijalva/jwt-go"
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

	_, err := s.users.CheckLogin(username, password)
	if err != nil {
		return "", err
	}

	return "", nil

	// create the claims
	/*claims := CustomClaims{
		u.Role.String(),
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer: "imsazon",
		}
	}*/

	/* 			OLD METHOD
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

	return signedToken, nil*/
}

func (s *service) Check(tokenString string) (user.UserId, error) {
	if tokenString == "" {
		return "", ErrInvalidArgument
	}

	// validate the token
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	switch err.(type) {
	case nil:
		// no error

		// we still have to check the token's validity
		if !token.Valid {
			return "", ErrInvalid
		}

		return user.U0001, nil

	case *jwt.ValidationError:
		// something went wrong during the validation
		vErr := err.(*jwt.ValidationError)

		switch vErr.Errors {
		case jwt.ValidationErrorExpired:
			// the token has expired
			return "", ErrExpired

		default:
			return "", ErrInvalid
		}

	default:
		// something else went wrong
		return "", ErrInvalid
	}
}

// NewService returns a new instance of the auth service
func NewService(users user.Repository) Service {
	return &service{
		users:		users,
	}
}