package paymentsgateway

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"
)

func ConfirmC2BTransactionEndpoint(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		FirstName string
		MiddleName string
		LastName string
		TransactionType string
		TransID string
		TransTime string
		TransAmount string
		BusinessShortCode int
		BillRefNumber int
		OrgAccountBalance float64
		ThirdPartyTransID string
		MSISDN string
	}
	request := Req{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&request)
	if err != nil {
		log.Fatal(err)
	}
	// TODO Confirm that this is a valid transaction
	res := Result{
		ResultCode: 0,
		ResultDesc: "The service was accepted successfully",
		ThirdPartyTransID: request.ThirdPartyTransID,
	}
	json.NewEncoder(w).Encode(res)
}

func LipaNaMpesaOnlineEndpoint(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		PhoneNumber string `json:"phone_number"`
		Amount int32 `json:"amount"`
		ApiKey string `json:"api_key"`
	}
	req := Request{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&req)
	if err != nil {
		log.Fatal(err)
	}
	// use Api key to load credentails and initalise the gateway
	gateway := MpesaClient{}

	response := gateway.InitiateLipaNaMpesaRequest(req.Amount, req.PhoneNumber, "Online Payment")
	
	result := Result{
		ResultCode: fmt.Sprintf("%d", response.ResultCode),
		ResultDesc: response.ResponseDescription,
		ThirdPartyTransID: response.CheckoutRequestID,
	}

	json.NewEncoder(w).Encode(&result)

}

func ValidateTransactionEndpoint(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		TransID string
		TransTime string 
		TransAmount float32 
		BusinessShortCode int
		BillRefNumber string 
		FirstName string 
		MiddleName string
		LastName string
	}
	request := Request{}
	dec := json.NewDecoder(r.Body)
	dec.Decode(&request)
	// Validate that this is a valid transaction

	res := Result{
		ResultCode: 0,
    	ResultDesc: "The service was accepted successfully",
        ThirdPartyTransID: "1234567890",
	}
	json.NewEncoder(w).Encode(&res)
}
