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
	return act, nil
}
