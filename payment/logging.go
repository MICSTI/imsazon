package payment

import (
	"github.com/go-kit/kit/log"
	"time"
)

type loggingService struct {
	logger log.Logger
	Service
}

// NewLoggingService returns an instance of a logging service
func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) Charge(charge CreditCardCharge) (status CreditCardChargeStatus, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "Charge",
			"successStatus", status,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Charge(charge)
}