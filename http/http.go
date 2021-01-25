/*
Package http is implementing f3client.AccountSvc interface
over HTTP communication protocol.

Atm this is the only protocol supported by target API (that being a REST API).
But having this implementation separately lays the basis of implementing the
interface over other protocols as well (e.g. graphQL, protobuffs, etc)
*/
package http

import (
	f3c "github.com/zerogvt/f3client"
)

// AccountSvc is an implementation of f3client.AccountSvc
// when the underlying transport is http.
type AccountSvc struct {
}

// Create creates an account using the REST HTTP API
func (*AccountSvc) Create(act *f3c.Account) (*f3c.Account, error) {
	// TODO
	return nil, nil
}
