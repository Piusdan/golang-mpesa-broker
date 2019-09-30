package paymentsgateway

import (
	"encoding/json"
	"net/http"
)

func DisburseEndpoint(w http.ResponseWriter, req *http.Request) {
	type Request struct {
		amount   float32
		msisdn   string
		clientID string
	}
	reqBody := Request{}
	type Response struct {
		Message string
	}
	var resp Response
	resp.Message = "Hello world"
	json.NewEncoder(w).Encode(resp)
	gateway := MpesaClient{}
	gatewayResp := gateway.SendB2C(reqBody.amount, reqBody.msisdn)
	
	json.NewEncoder(w).Encode(gatewayResp)
}
