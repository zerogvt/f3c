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
	"fmt"
	"io/ioutil"
	"net/http"
)

type Account struct {
	Type           string     `json:"type"`
	ID             string     `json:"id"`
	OrganisationID string     `json:"organisation_id"`
	Attributes     Attributes `json:"attributes"`
}

type AccountXL struct {
	Account
	Version       int           `json:"version"`
	Relationships Relationships `json:"relationships"`
}

type PrivateIdentification struct {
	BirthDate      string   `json:"birth_date"`
	BirthCountry   string   `json:"birth_country"`
	Identification string   `json:"identification"`
	Address        []string `json:"address"`
	City           string   `json:"city"`
	Country        string   `json:"country"`
}
type Actors struct {
	Name      []string `json:"name"`
	BirthDate string   `json:"birth_date"`
	Residency string   `json:"residency"`
}
type OrganisationIdentification struct {
	Identification string   `json:"identification"`
	Actors         []Actors `json:"actors"`
	Address        []string `json:"address"`
	City           string   `json:"city"`
	Country        string   `json:"country"`
}
type Attributes struct {
	Country                    string                     `json:"country"`
	BaseCurrency               string                     `json:"base_currency"`
	AccountNumber              string                     `json:"account_number"`
	BankID                     string                     `json:"bank_id"`
	BankIDCode                 string                     `json:"bank_id_code"`
	Bic                        string                     `json:"bic"`
	Iban                       string                     `json:"iban"`
	Name                       []string                   `json:"name"`
	AlternativeNames           []string                   `json:"alternative_names"`
	AccountClassification      string                     `json:"account_classification"`
	JointAccount               bool                       `json:"joint_account"`
	AccountMatchingOptOut      bool                       `json:"account_matching_opt_out"`
	SecondaryIdentification    string                     `json:"secondary_identification"`
	Switched                   bool                       `json:"switched"`
	PrivateIdentification      PrivateIdentification      `json:"private_identification"`
	OrganisationIdentification OrganisationIdentification `json:"organisation_identification"`
	Status                     string                     `json:"status"`
}

type Relationships struct {
	MasterAccount MasterAccount `json:"master_account"`
	AccountEvents AccountEvents `json:"account_events"`
}

type Rel struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type MasterAccount struct {
	Rels []Rel `json:"data"`
}

type AccountEvents struct {
	Evts []Rel `json:"data"`
}

// internal use to help unmarshalling responses
type payloadOut struct {
	Account Account `json:"data"`
}

type payloadIn struct {
	Account AccountXL `json:"data"`
}

type payloadInArr struct {
	Accounts []AccountXL `json:"data"`
}

// AccountSvc encompasses account-related actions.
// As per specs in https://github.com/form3tech-oss/interview-accountapi
// we need to implement Create, Fetch, List and Delete.
type AccountSvc interface {
	Create(ac Account) (AccountXL, error)
	Fetch(id string) (AccountXL, error)
	List(page int, pagesize int) ([]AccountXL, error)
	Delete(id string) error
}

func Pprint(data interface{}) {
	json, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return
	}
	fmt.Println(string(json))
}

// Payload creates a reader out of an account that can be used as body
// in a POST or GET.
func (act *Account) ToPayload() *bytes.Buffer {
	var data []byte
	var err error
	payload := payloadOut{*act}
	//Pprint(payload)
	if data, err = json.Marshal(payload); err != nil {
		return nil
	}
	return bytes.NewBuffer([]byte(data))
}

// Payload creates a reader out of an account that can be used as body
// in a POST or GET.
func FromPayload(r *http.Response) (AccountXL, error) {
	res := AccountXL{}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return res, err
	}
	data := payloadIn{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return res, err
	}
	res = data.Account
	return res, nil
}

// Payload creates a reader out of an account that can be used as body
// in a POST or GET.
func FromPayloadArr(r *http.Response) ([]AccountXL, error) {
	res := []AccountXL{}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return res, err
	}
	data := payloadInArr{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return res, err
	}
	res = data.Accounts
	return res, nil
}
