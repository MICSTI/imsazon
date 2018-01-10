package stock

import (
	"github.com/go-kit/kit/log"
	productModel "github.com/MICSTI/imsazon/models/product"
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

func (s *loggingService) GetItems() (products []*productModel.Product) {
	defer func(begin time.Time) {
		s.logger.Log("method", "GetItems", "took", time.Since(begin), "err", nil)
	}(time.Now())
	return s.Service.GetItems()
}

func (s *loggingService) Add(productToAdd *productModel.Product) (updatedProduct *productModel.Product, err error) {
	defer func(begin time.Time) {
		s.logger.Log("method", "Add", "product_id", updatedProduct.Id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return s.Service.Add(productToAdd)
}

func (s *loggingService) Withdraw(productToWithdraw *productModel.Product) (updatedProduct *productModel.Product, err error) {
	defer func(begin time.Time) {
		s.logger.Log("method", "Withdraw", "product_id", updatedProduct.Id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return s.Service.Withdraw(productToWithdraw)
}

