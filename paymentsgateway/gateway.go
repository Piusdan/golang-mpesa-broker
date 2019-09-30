package paymentsgateway

import (
	"time"
)

type PaymentGateway interface {
	SendB2C(amount float32, recipientMSISDN string) (PaymentGatewayResult, error)
	MakeOnlinePayment(amount float64, senderMSISDN string) (PaymentGatewayResult, error)
}

type PaymentGatewayResult struct {
	Status     string
	Message 	string
	Timestamp  time.Time
	StatusCode string
	Errors     string
	ErrorCode string
}
