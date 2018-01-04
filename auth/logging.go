package auth

import (
	"github.com/go-kit/kit/log"
	"time"
	"github.com/MICSTI/imsazon/user"
)

type loggingService struct {
	logger log.Logger
	Service
}

// NewLogging Service returns an instance of a logging Service
func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) Login(username string, password string) (signedToken string, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "Login",
			"username", username,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Login(username, password)
}

func (s *loggingService) Check(tokenString string) (userId user.UserId, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "Check",
			"userId", userId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Check(tokenString)
}