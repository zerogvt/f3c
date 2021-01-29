/*
Package f3c is a client library for Form3 API.

Design rationale:
	We let the json request and responses as per
	https://api-docs.form3.tech/api.html#organisation-accounts-create
	guide us in creating the domain data types.
	Loosely following design filosophy as worded in
	https://www.gobeyond.dev/standard-package-layout/

*/
package f3c

import (
	"bytes"
	"encoding/json"
)

// internal use only
type payload struct {
	Data Account `json:"data"`
}

// Account is the core structure representing an account.
type Account struct {
	Type           string     `json:"type"`
	ID             string     `json:"id"`
	OrganisationID string     `json:"organisation_id"`
	Attributes     Attributes `json:"attributes"`
}

// Attributes are account attributes.
type Attributes struct {
	Country       string `json:"country"`
	BaseCurrency  string `json:"base_currency"`
	BankID        string `json:"bank_id"`
	BankIDCode    string `json:"bank_id_code"`
	AccountNumber string `json:"account_number"`
	BIC           string `json:"bic"`
	IBAN          string `json:"iban"`
	CustomerID    string `json:"customer_id"`
}

// AccountSvc encompasses account-related actions.
// As per specs in https://github.com/form3tech-oss/interview-accountapi
// we need to implement Create, Fetch, List and Delete.
type AccountSvc interface {
	Create(ac *Account) (*Account, error)
	Fetch(id string) (Account, error)
	List(page int, filter interface{}) ([]Account, error)
	Delete(id string) error
}

// Payload creates a reader out of an account that can be used as body
// in a POST or GET.
func (act *Account) Payload() *bytes.Buffer {
	var data []byte
	var err error
	payload := payload{*act}
	if data, err = json.Marshal(payload); err != nil {
		return nil
	}
	return bytes.NewBuffer([]byte(data))
}

// AccountCrResp represents the expected response when we create a new account.
type AccountCrResp struct {
	Data struct {
		Type           string `json:"type"`
		ID             string `json:"id"`
		Version        int    `json:"version"`
		OrganisationID string `json:"organisation_id"`
		Attributes     struct {
			Country       string `json:"country"`
			BaseCurrency  string `json:"base_currency"`
			AccountNumber string `json:"account_number"`
			BankID        string `json:"bank_id"`
			BankIDCode    string `json:"bank_id_code"`
			BIC           string `json:"bic"`
			IBAN          string `json:"iban"`
			Status        string `json:"status"`
		} `json:"attributes"`
		Relationships struct {
			AccountEvents struct {
				Data []struct {
					Type string `json:"type"`
					ID   string `json:"id"`
				} `json:"data"`
			} `json:"account_events"`
		} `json:"relationships"`
	} `json:"data"`
}
