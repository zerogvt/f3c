/*
Package f3c is a client library for Form3 API.

Design rationale:
	We let the json request and responses as per
	https://api-docs.form3.tech/api.html#organisation-accounts-create
	dictate the domain data types.
*/
package f3c

// Account models a user account.
type Account struct {
	Type           string     `json:"type"`
	ID             string     `json:"id"`
	OrganisationID string     `json:"organisation_id"`
	Attributes     Attributes `json:"attributes"`
}

// NewAccount creates an account that can be used as input to CRUD operations.
func NewAccount(uid string, orgid string, attr Attributes) Account {
	return Account{
		Type:           "accounts",
		ID:             uid,
		OrganisationID: orgid,
		Attributes:     attr,
	}
}

// AccountXL is a user account with added metadata.
type AccountXL struct {
	Account
	Version       int           `json:"version"`
	Relationships Relationships `json:"relationships"`
}

// PrivateIdentification models user personal identification data.
type PrivateIdentification struct {
	BirthDate      string   `json:"birth_date"`
	BirthCountry   string   `json:"birth_country"`
	Identification string   `json:"identification"`
	Address        []string `json:"address"`
	City           string   `json:"city"`
	Country        string   `json:"country"`
}

// Actors models actors. (Fixme: no idea what this is in the real world).
type Actors struct {
	Name      []string `json:"name"`
	BirthDate string   `json:"birth_date"`
	Residency string   `json:"residency"`
}

// OrganisationIdentification models org id data.
type OrganisationIdentification struct {
	Identification string   `json:"identification"`
	Actors         []Actors `json:"actors"`
	Address        []string `json:"address"`
	City           string   `json:"city"`
	Country        string   `json:"country"`
}

// Attributes model user account attributes.
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

// Relationships models relationships of this account.
type Relationships struct {
	MasterAccount MasterAccount `json:"master_account"`
	AccountEvents AccountEvents `json:"account_events"`
}

// Rel models a relationship.
type Rel struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// MasterAccount models relationships to master account.
type MasterAccount struct {
	Rels []Rel `json:"data"`
}

// AccountEvents models events happened on this account.
type AccountEvents struct {
	Evts []Rel `json:"data"`
}

// PayloadOut is a helper struct to model the json data that client sends.
type PayloadOut struct {
	Account Account `json:"data"`
}

// PayloadIn is a helper struct to model the json data that client receives.
type PayloadIn struct {
	Account AccountXL `json:"data"`
}

// PayloadInArr is the array equivalent for PayloadIn.
type PayloadInArr struct {
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
