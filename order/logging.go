package order

import (
	"github.com/go-kit/kit/log"
	orderModel "github.com/MICSTI/imsazon/models/order"
	userModel "github.com/MICSTI/imsazon/models/user"
	"time"
)

type loggingService struct {
	logger		log.Logger
	Service
}

// NewLoggingService returns an instance of a logging Service.
func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) Create(newOrder *orderModel.Order) (order *orderModel.Order, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "Create",
			"orderId", newOrder.Id,
			"userId", newOrder.UserId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Create(newOrder)
}

func (s *loggingService) UpdateStatus(id orderModel.OrderId, newStatus orderModel.OrderStatus) (order *orderModel.Order, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "UpdateStatus",
			"orderId", id,
			"newStatus", newStatus.String(),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.UpdateStatus(id, newStatus)
}

func (s *loggingService) GetById(id orderModel.OrderId) (order *orderModel.Order, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "GetById",
			"orderId", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.GetById(id)
}

func (s *loggingService) GetAll() []*orderModel.Order {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "GetAll",
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.GetAll()
}

func (s *loggingService) GetAllForUser(userId userModel.UserId) []*orderModel.Order {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "GetAllForUser",
			"userId", userId,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.GetAllForUser(userId)
}