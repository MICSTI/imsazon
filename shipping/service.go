package shipping

import (
	"errors"
	orderModel "github.com/MICSTI/imsazon/models/order"
	"time"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"bytes"
	"github.com/MICSTI/imsazon/mail"
)

const getSingleOrderApiUrl = "http://localhost:8605/order/single/"
const updateOrderStatusApiUrl = "http://localhost:8605/order/update/"
const sendMailApiUrl = "http://localhost:8605/mail/send"

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
	testMailRecipient		string
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

	// call the mail service to send out an email that the order was shipped successfully
	mailToSend := mail.New(s.testMailRecipient, "Your order has been shipped", successfulShippingMailBody, "text/html")

	err = sendMail(mailToSend)

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

func sendMail(mail *mail.Email) error {
	message := map[string]interface{}{
		"to": mail.To,
		"subject": mail.Subject,
		"body": mail.Body,
		"contentType": mail.ContentType,
	}

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		return err
	}

	resp, err := http.Post(sendMailApiUrl, "application/json", bytes.NewBuffer(bytesRepresentation))

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
func NewService(testMailRecipient string) Service {
	return &service{
		testMailRecipient: testMailRecipient,
	}
}

const successfulShippingMailBody = `
	<div style="font-size: 18pt; font-weight: bold; text-align: center; margin-bottom: 16px;">IMSazon</div>
	<div style="font-size: 14pt; margin-bottom: 16px;">Your order has been shipped succcessfully!</div>
	<div style="font-size: 12pt; margin-bottom: 10px;">Thank you so much for ordering from us, we hope you will be delighted with your new things</div> 
	<div style="font-size: 12pt; margin-bottom: 10px;">- Michael from <b>IMSazon</b>
`