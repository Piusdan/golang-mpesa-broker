package main

import (
	"log"
	// "net/http"
	// "time"
	"fmt"

	"github.com/Piusdan/payments-gateway/paymentsgateway"
)

// func main() {
// 	router := paymentsgateway.NewRouter()
// 	server := &http.Server{
// 		Handler: router,
// 		Addr:    "0.0.0.0:8000",
// 		// Good practice: enforce timeouts for servers you create!
// 		WriteTimeout: 15 * time.Second,
// 		ReadTimeout:  15 * time.Second,
// 	}
// 	log.Fatal(server.ListenAndServe())
// }

func main() {
	gateway := paymentsgateway.MpesaClient{}
	gateway.BaseURL = "https://api.safaricom.co.ke"
	gateway.C2bConsumerKey = "W3PdaNeeu6rHGaG2SPX63sBZcR9GUNBq"
	gateway.C2bConsumerSecret = "4dos8cg1F9h9oseQ"
	gateway.Username = "WezaVenturesB2CInit"
	gateway.SecurityCredential = "N9K0K8Yt/4NN1TqF1ZB6oEIu96eybfbbpufKV/xR+J2nxOahG5Mnqir0xBG2wV5B81j82k45h3okbsM938Msfk8HGpMtsCI9SS6MJqgmTeaoTYWt4OEk30ur2ksae+rT6jY2xA4kGIBAue355x8w6c44hWpQA/yYvhHg20Xp/ybZcS3EoZY0WWdadrWKeO+Zl1les3URVipxJjZlz9JmrLg/LS4wCpqv6x1xahtxld8BQJxSPKaKduaT6CRog3WHF/MJbFKk/RIM+5LY/2F7Go/+jPijXVdUJJJRwqGAXLkdlw2AZ/vHEGXBQLXXzzKyTeBxar43VOo8uuPtvoNhVQ=="
	gateway.ResultCallbackURL = "https://mpesa-broker.api.loanbee.tech/api/v1/result"
	gateway.QueueTimeOutURL = "https://mpesa-broker.api.loanbee.tech/api/v1/result"
	gateway.OrganisationShortcode = "327413"
	err := gateway.GetAccessToken()
	if err != nil {
		log.Fatalf("Unable to authorize request %v", err)
	}
	resp := gateway.SendB2C(5.0, "254703554404")
	fmt.Printf("message: %s\nStatus: %s", resp.Message, resp.Status)
}
