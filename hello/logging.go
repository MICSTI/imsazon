package hello

import (
	"github.com/go-kit/kit/log"
	"time"
)

type loggingService struct {
	logger log.Logger
	Service
}

// NewLoggingService returns a new instance of a logging Service
func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) SayHello() (greeting string) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "SayHello",
			"greeting", greeting,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.SayHello()
}

func (s *loggingService) SayHelloTo(name string) (greeting string, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "SayHelloTo",
			"greeting", greeting,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.SayHelloTo(name)
}