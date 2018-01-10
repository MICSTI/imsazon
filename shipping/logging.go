package shipping

import (
	"github.com/go-kit/kit/log"
	orderModel "github.com/MICSTI/imsazon/models/order"
	"time"
)

type loggingService struct {
	logger log.Logger
	Service
}

// NewLoggingService returns a new instance of a logging Service.
func NewLoggingService (logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) Ship(orderId orderModel.OrderId) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "Ship",
			"orderId", orderId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Ship(orderId)
}