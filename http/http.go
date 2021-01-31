/*
Package http is implementing f3client.AccountSvc interface
over HTTP communication protocol.

Atm this is the only protocol supported by target API (that being a REST API).
But having this implementation separately lays the basis of implementing the
interface over other protocols as well (e.g. graphQL, protobuffs, etc)
*/
package http

import (
	"fmt"
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
func (svc *AccountSvc) Create(act f3c.Account) (f3c.AccountXL, error) {
	res := f3c.AccountXL{}
	r, err := svc.Cli.Post(svc.Base+"/v1/organisation/accounts",
		"application/vnd.api+json",
		act.ToPayload(),
	)
	if err != nil {
		return res, err
	}
	if err = failed(r); err != nil {
		return res, err
	}
	res, err = f3c.FromPayload(r)
	return res, err
}

// Fetch gets an account id and fetches the account data.
func (svc *AccountSvc) Fetch(id string) (f3c.AccountXL, error) {
	res := f3c.AccountXL{}
	r, err := svc.Cli.Get(svc.Base + "/v1/organisation/accounts/" + id)
	if err != nil {
		return res, err
	}
	if err = failed(r); err != nil {
		return res, err
	}
	res, err = f3c.FromPayload(r)
	return res, err
}

// Delete deletes an account based on account id and version.
func (svc *AccountSvc) Delete(id string, version int) error {
	ep := fmt.Sprintf("%s%s%s%s%d",
		svc.Base,
		"/v1/organisation/accounts/",
		id,
		"?version=", version)
	req, err := http.NewRequest("DELETE", ep, nil)
	if err != nil {
		return err
	}
	r, err := svc.Cli.Do(req)
	if err != nil {
		return err
	}
	if err := failed(r); err != nil {
		return err
	}
	return nil
}

// List gets a list of all accounts. It supports paging.
func (svc *AccountSvc) List(page int, pagesize int) ([]f3c.AccountXL, error) {
	res := []f3c.AccountXL{}
	ep := fmt.Sprintf("%s%s%s%d%s%d",
		svc.Base,
		"/v1/organisation/accounts/",
		"?page[number]=", page,
		"&page[size]=", pagesize)
	r, err := svc.Cli.Get(ep)
	if err != nil {
		return res, err
	}
	if err = failed(r); err != nil {
		return res, err
	}
	res, err = f3c.FromPayloadArr(r)
	return res, err
}

// Err represents an HTTP client or server error (i.e. status code >= 400)
type Err struct {
	Code int
	Text string
}

func (e Err) Error() string {
	return fmt.Sprintf("HTTP Error: %d, %s", e.Code, e.Text)
}

func failed(r *http.Response) error {
	fmt.Println("HTTP Response Status:",
		r.StatusCode,
		http.StatusText(r.StatusCode))
	if !(r.StatusCode >= 200 && r.StatusCode < 300) {
		return Err{r.StatusCode, http.StatusText(r.StatusCode)}
	}
	return nil
}
