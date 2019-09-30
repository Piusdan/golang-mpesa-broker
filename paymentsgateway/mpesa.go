package paymentsgateway

import (
	"encoding/base64"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type MpesaClient struct {
	BaseURL               string
	Password              string
	BussinessShortCode    string
	Username              string
	SecurityCredential    string
	OrganisationShortcode string
	AccessToken           string
	ResultCallbackURL     string
	QueueTimeOutURL       string
	C2bConsumerKey        string
	C2bConsumerSecret     string
	LipaNaMpesaOnlineCallbackURL       string
}

// define result that will be returned by mpesa's endpoint
type b2cRequestResult struct {
	RequestID                string `json:"requestId"`
	ErrorCode                string `json:"errorCode"`
	ErrorMessage             string `json:"errorMessage"`
	ConversationID           string `json:"ConversationID"`
	OriginatorConversationID string `json:"OriginatorConversationID"`
	ResponseCode             string `json:"ResponseCode"`
	ResponseDescription      string `json:"ResponseDescription"`
}

type lipaNaMpesaOnlineResult struct {
	MerchantRequestID	string
	CheckoutRequestID	string
	ResponseCode	string
	ResultDesc	string
	ResponseDescription	 string
	ResultCode string
}

// result send to the mpesa callback url
type b2cResult struct {
	Result struct {
		ResultType               int
		ResultCode               int
		ResultDesc               string
		OriginatorConversationID string
		ConversationID           string
		TransactionID            string
		ResultParameters         struct {
			ResultParameter []struct {
				Key   string
				Value interface{}
			}
		}
		ReferenceData struct {
			ReferenceItem struct {
				Key   string
				Value string
			}
		}
	}
}

func addHeaders(req *http.Request, headers map[string]string) {
	for key, value := range headers {
		req.Header.Add(key, value)
	}
}

func (mpesa MpesaClient) makeURLString(endpoint string) string {
	return mpesa.BaseURL + endpoint
}

func (mpesa MpesaClient) httpPostRequest(endpoint string, headers map[string]string, payload interface{}) (res http.Response, e error) {
	/* make http request */
	client := &http.Client{
		Timeout: time.Second * 15,
	}
	// parse url
	urlString := mpesa.makeURLString(endpoint)
	// serialize payload to json
	bPayload, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", urlString, bytes.NewReader(bPayload))
	// add content-type headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	// add other headers
	addHeaders(req, headers)
	// make http call
	resp, err := client.Do(req)
	return *resp, err
}

// GetAccessToken returns an accesstoken to be used when calling mpesa APIs
func (mpesa *MpesaClient) GetAccessToken() (error) {
	const TokenEndpoint = "/oauth/v1/generate?grant_type=client_credentials"
	urlString := mpesa.makeURLString(TokenEndpoint)
	req, _ := http.NewRequest("GET", urlString, nil)
	headers := map[string]string{
		"cache-control": "no-cache",
		"Accept":        "application/json",
	}
	addHeaders(req, headers)
	req.SetBasicAuth(mpesa.C2bConsumerKey, mpesa.C2bConsumerSecret)
	client := &http.Client{
		Timeout: time.Second * 120,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to obtain accesstoken %v", err)
	}

	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   string `json:"expires_in"`
	}

	dec := json.NewDecoder(resp.Body)
	dec.Decode(&result)

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.New(fmt.Sprintf("%s", body))
	}
	mpesa.AccessToken = result.AccessToken

	return nil
}

func (mpesa MpesaClient) makeB2CPayment(amount float32, recipient string, remarks string) (b2cRequestResult, error) {
	const B2CEndpoint = "/mpesa/b2c/v1/paymentrequest"

	// Mpesa request body
	type Payload struct {
		InitiatorName      string
		SecurityCredential string
		CommandID          string
		Amount             float32
		PartyA             string
		PartyB             string
		Remarks            string
		QueueTimeOutURL    string
		ResultURL          string
		Occasion           string
	}
	// initiate the mpesa request body
	payload := Payload{
		InitiatorName:      mpesa.Username,
		SecurityCredential: mpesa.SecurityCredential,
		CommandID:          "BusinessPayment",
		Amount:             amount,
		PartyA:             mpesa.OrganisationShortcode,
		PartyB:             recipient,
		Remarks:            remarks,
		QueueTimeOutURL:    mpesa.QueueTimeOutURL,
		ResultURL:          mpesa.ResultCallbackURL,
		Occasion:           " ",
	}
	// make auth headers
	authHeader := fmt.Sprintf("Bearer %s", mpesa.AccessToken)
	// add other headers
	headers := map[string]string{
		"Authorization": authHeader,
	}
	resp, err := mpesa.httpPostRequest(B2CEndpoint, headers, payload)
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	result := b2cRequestResult{}
	dec.Decode(&result)

	if err != nil {
		log.Fatalf("Unable to complete request %v", err)
	}
	// return an error if request wasn't successful
	if resp.StatusCode != 200 {
		result.ResponseCode = result.ErrorCode
		return result, errors.New(result.ErrorMessage)
	}
	return result, nil

}

// Parseb2cResult parses data received from the mpesa gateway
// it returns an error if the data can't be parsed
func (mpesa MpesaClient) Parseb2cResult(r io.Reader) (b2cResult, error) {
	result := b2cResult{}
	dec := json.NewDecoder(r)
	err := dec.Decode(&result)
	if err != nil {
		return result, errors.New(fmt.Sprintf("Unable to proccess c2b result payload %v", err))
	}
	return result, nil
}

func (mpesa MpesaClient) ValidateC2BTransaction() {
}

// SendB2C initiates a bussiness to customer transaction
// given amount and recipientMSISDN which is the phonenumber receiveing the funds
// its makes an http req to the mpesa endpoint
func (mpesa MpesaClient) SendB2C(amount float32, recipientMSISDN string) PaymentGatewayResult {
	res, err := mpesa.makeB2CPayment(amount, recipientMSISDN, "BusinessPayment")
	result := PaymentGatewayResult{
		Timestamp:  time.Now(),
		StatusCode: res.ResponseCode,
		Status:     "Success",
		Message:    res.ResponseDescription,
	}
	if err != nil {
		result.Message = res.ErrorMessage
		result.Errors = err.Error()
		result.Status = "Failed"
		result.StatusCode = res.ErrorCode
	}
	return result
}

func (mpesa MpesaClient) InitiateLipaNaMpesaRequest(amount int32, recipient string, transactionDescription string) (lipaNaMpesaOnlineResult){
	const Endpoint = "mpesa/stkpush/v1/processrequest"
	type Payload struct {
        BusinessShortCode string
        Password string
        Timestamp string
        TransactionType string
        Amount float32
        PartyA string
        PartyB string
        PhoneNumber string
        CallBackURL string
        AccountReference string
        TransactionDesc string
	}
	tNow := time.Now()
	timestamp := tNow.Format("20060102030405")
	password := mpesa.BussinessShortCode + mpesa.Password + timestamp
	b64Password := base64.StdEncoding.EncodeToString([]byte(password))
	payload := Payload{
		BusinessShortCode: mpesa.BussinessShortCode,
		Password: b64Password,
		Timestamp: timestamp,
		Amount: float32(amount),
		PartyA: recipient,
		PartyB: mpesa.OrganisationShortcode,
		PhoneNumber: recipient,
		CallBackURL: mpesa.LipaNaMpesaOnlineCallbackURL,
		TransactionDesc: transactionDescription,
	}
	result := lipaNaMpesaOnlineResult{}
	authHeader := fmt.Sprintf("Bearer %s", mpesa.AccessToken)
	headers := map[string]string{
		"Authorization": authHeader,
	}
	resp, err := mpesa.httpPostRequest(Endpoint, headers, payload)
	if err != nil {
		log.Fatalf("Unable to complete request %v", err)
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&result)
	return result
}
