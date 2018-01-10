package cart

import (
	"github.com/go-kit/kit/log"
	productModel "github.com/MICSTI/imsazon/models/product"
	userModel "github.com/MICSTI/imsazon/models/user"
	"time"
)

type loggingService struct {
	logger log.Logger
	Service
}

// NewLoggingService returns a new instance of a logging service
func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) GetCart(userId userModel.UserId) (cartItems []*productModel.SimpleProduct, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "GetCart",
			"userId", userId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.GetCart(userId)
}

func (s *loggingService) Put(userId userModel.UserId, productId productModel.ProductId, quantity int) (updatedCart []*productModel.SimpleProduct, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "Put",
			"userId", userId,
			"productId", productId,
			"quantity", quantity,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Put(userId, productId, quantity)
}

func (s *loggingService) Remove(userId userModel.UserId, productId productModel.ProductId) (updatedCart []*productModel.SimpleProduct, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "Remove",
			"userId", userId,
			"productId", productId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Remove(userId, productId)
}