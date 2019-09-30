package model

import (
	"time"
)

type Credential struct {
	ID string
	BaseURL               string
	Password              string
	BussinessShortCode    string
	Username              string
	SecurityCredential    string
	OrganisationShortcode string
	AccessToken           string
	resultCallbackURL     string
	QueueTimeOutURL       string
	C2bConsumerKey string
	C2bConsumerSecret string
	B2cConsumerKey string
	B2bConsumerKey string
	created_date *time.Time
	last_updated *time.Time
	last_updated_by string
	app_name string
	transaction_id string
}

var credentials []Credential

