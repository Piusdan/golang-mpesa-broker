package paymentsgateway

import (
	"log"
	"encoding/json"
	"net/http"
)

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
