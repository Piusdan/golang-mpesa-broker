package paymentsgateway

import (
	"log"
	"encoding/json"
	"net/http"
)

type Result struct {
	ResultCode int
	ResultDesc string
	ThirdPartyTransID string
}

func DisburseEndpoint(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Amount   float32 `json:"amount`
		PhoneNumber   string `json:"phone"`
		ApiKey string  `json:"api_key"`
	}
	request := Request{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&request)
	if err != nil {
		log.Fatal(err)
	}

	//TODO use ApiKey to load credentials and Initialise gateway
	gateway := MpesaClient{}
	gatewayResp := gateway.SendB2C(request.Amount, request.PhoneNumber)
	// TODO save the request
	json.NewEncoder(w).Encode(gatewayResp)
}


func B2CResultEndpoint(w http.ResponseWriter, r *http.Request) {
	gateway := MpesaClient{}
	req, err := gateway.Parseb2cResult(r.Body)
	res := Result{}
	if err != nil {
		log.Fatal(err)
	}
	// save transaction details

	// make response
	res.ResultDesc = req.Result.ResultDesc
	res.ResultCode = req.Result.ResultCode
	res.ThirdPartyTransID = req.Result.OriginatorConversationID
	json.NewEncoder(w).Encode(res)

}

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