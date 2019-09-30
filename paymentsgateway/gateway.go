package paymentsgateway

import (
	"time"
)

type PaymentGateway interface {
	SendB2C(amount float32, recipientMSISDN string) (PaymentGatewayResult, error)
	MakeOnlinePayment(amount float64, senderMSISDN string) (PaymentGatewayResult, error)
}

type PaymentGatewayResult struct {
	Status     string `json:"status"`
	Message 	string `json:"message"`
	Timestamp  time.Time `json:"timestamp"`
	StatusCode string  `json:"status_code"`
	Errors     string `json:"errors"`
	ErrorCode string `json:"error_code"`
}

type Result struct {
	ResultCode int
	ResultDesc string
	ThirdPartyTransID string
}
