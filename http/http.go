/*
Package http is implementing f3client.AccountSvc interface
over HTTP communication protocol.

Atm this is the only protocol supported by target API (that being a REST API).
But having this implementation separately lays the basis of implementing the
interface over other protocols as well (e.g. graphQL, protobuffs, etc)
*/
package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/zerogvt/f3c"
)

// AccountSvc is an implementation of f3client.AccountSvc
// when the underlying transport is http.
type AccountSvc struct {
	// Allow users to set their own clients if they want to.
	// why? https://youtu.be/cmkKxNN7cs4?t=1509
	// The default zero value is fine as well.
	Cli http.Client
	// base must be set to the base url of the REST server
	Base string
}

// Create creates an account using the REST HTTP API
func (svc *AccountSvc) Create(act f3c.Account) (f3c.AccountCrResp, error) {
	res := f3c.AccountCrResp{}
	r, err := svc.Cli.Post(svc.Base+"/v1/organisation/accounts",
		"application/vnd.api+json",
		act.Payload(),
	)
	if err != nil {
		return res, err
	}
	fmt.Println("HTTP Response Status:",
		r.StatusCode,
		http.StatusText(r.StatusCode))
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, &res)
	return res, err
}
