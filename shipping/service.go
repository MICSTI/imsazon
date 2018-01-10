package shipping

import (
	"errors"
	orderModel "github.com/MICSTI/imsazon/models/order"
	"time"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"bytes"
)

const getSingleOrderApiUrl = "http://localhost:8605/order/single/"
const updateOrderStatusApiUrl = "http://localhost:8605/order/update/"

// ErrInvalidArgument is returned when one or more arguments are invalid.
var ErrInvalidArgument = errors.New("Invalid argument")
var ErrShippingNotPossible = errors.New("Shipping is currently not possible for this order")
var ErrInvalidOperation = errors.New("Invalid operation")
var ErrApi = errors.New("Error response from API")

// Service is the interface that provides the shipping methods
type Service interface {
	// Ships the order from the physical store
	Ship(orderModel.OrderId) error
}

type service struct {

}

type OrderStatusApiResponse struct {
	Order			*orderModel.Order		`json:"order"`
}

func (s *service) Ship(orderId orderModel.OrderId) (error) {
	if orderId == "" {
		return ErrInvalidArgument
	}

	// check the order service if the current order status is "Payment Successful"
	currentOrderStatus, err := getOrderStatus(orderId)

	if err != nil {
		return err
	}

	if currentOrderStatus != orderModel.PaymentSuccessful {
		return ErrInvalidOperation
	}

	// we can't really do anything, so we just add a delay and trigger the sending of an email
	duration := time.Millisecond * 750
	time.Sleep(duration)

	// call order service to mark order as "shipped"
	err = setOrderStatus(orderId, orderModel.Shipped)

	if err != nil {
		return ErrShippingNotPossible
	}

	// TODO call mail service


	return nil
}

func parseOrderStatusResponse(body []byte) (*OrderStatusApiResponse, error) {
	var s = new(OrderStatusApiResponse)
	err := json.Unmarshal(body, &s)
	if err != nil {
		return nil, err
	}
	return s, err
}

func getOrderStatus(id orderModel.OrderId) (orderModel.OrderStatus, error) {
	resp, err := http.Get(getSingleOrderApiUrl + id.String())

	if err != nil {
		return orderModel.Created, ErrApi
	}

	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return orderModel.Created, ErrApi
	}

	parsed, err := parseOrderStatusResponse(bodyBytes)

	if err != nil {
		return orderModel.Created, err
	}

	return parsed.Order.Status, nil
}

func setOrderStatus(id orderModel.OrderId, newStatus orderModel.OrderStatus) error {
	message := map[string]interface{}{
		"status": newStatus,
	}

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		return err
	}

	resp, err := http.Post(updateOrderStatusApiUrl + id.String(), "application/json", bytes.NewBuffer(bytesRepresentation))

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ErrApi
	}

	return nil
}

// NewService creates a shipping service
func NewService() Service {
	return &service{

	}
}