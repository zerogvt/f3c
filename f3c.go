/*
Package f3c is a client library for Form3 API.

Design rationale:
	We let the json request and responses as per
	https://api-docs.form3.tech/api.html#organisation-accounts-create
	guide us in creating the domain data types.
	Loosely following design filosophy as worded in
	https://www.gobeyond.dev/standard-package-layout/

Notes:
	Only definitions hosted in this file.
*/
package f3c

// Attributes describe account metadata
type Attributes struct {
	Country    string `json:"country"`
	BaseCurr   string `json:"base_currency"`
	BankID     string `json:"bank_id"`
	BankIDCode string `json:"bank_id_code"`
	BIC        string `json:"bic"`
}

// Account is the basic data model for an F3 account
type Account struct {
	ID    string `json:"id"`
	OrgID string `json:"organisation_id"`
	// composition through embedding
	Attributes
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
