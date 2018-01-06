package stock

import (
	"github.com/go-kit/kit/log"
	"github.com/MICSTI/imsazon/product"
	"time"
)

type loggingService struct {
	logger log.Logger
	Service
}

// NewLoggingService returns a new instace of a logging service
func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) GetItems() (products []product.Product, err error) {
	defer func(begin time.Time) {
		s.logger.Log("method", "GetItems", "took", time.Since(begin), "err", err)
	}(time.Now())
	return s.Service.GetItems()
}

func (s *loggingService) Add(productToAdd *product.Product) (updatedProduct product.Product, err error) {
	defer func(begin time.Time) {
		s.logger.Log("method", "Add", "product_id", updatedProduct.Id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return s.Service.Add(productToAdd)
}

func (s *loggingService) Withdraw(productToWithdraw *product.Product) (updatedProduct product.Product, err error) {
	defer func(begin time.Time) {
		s.logger.Log("method", "Withdraw", "product_id", updatedProduct.Id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return s.Service.Withdraw(productToWithdraw)
}

