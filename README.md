[![Build Status](https://travis-ci.com/zerogvt/f3c.svg?token=zDwdCt1iLHUMB3eVQDMy&branch=main)](https://travis-ci.com/github/zerogvt/f3c)

# Client library for Form3 API
as per https://github.com/form3tech-oss/interview-accountapi

# Author
Vasileios Gkoumplias

# Run all-in-one tests
This will have the needed docker image built and run the docker-compose with the unit tests. Tested in OSx.

`make test`

In case you have trouble running it you can still see the tests run in the [latest travis build](https://travis-ci.com/github/zerogvt/f3c).


# Import
```
import (
    "github.com/zerogvt/f3c"
    "github.com/zerogvt/f3c/http"
)
```

# Create
To create an account to Form3 system you need to first define the basic account elements locally using an `f3c.Attributes` composite literal with at least the minimum required fields and then have `f3c.NewAccount()` bind them in a local account.

You can then use that account as input to `AccountSvc.Create()` which will create the account in Form3 remote system.

Next snippet should clarify these steps:
```
import (
    "github.com/zerogvt/f3c"
    "github.com/zerogvt/f3c/http"
)

// >> create a local Account instance
uid := "ad27e265-9605-4b4b-a0e5-123456789012"
oid := "eb0bd6f5-c3f5-44b2-b677-123456789012"
attr := f3c.Attributes{
    Country:               "GB",
    BaseCurrency:          "GBP",
    BankID:                "400300",
    BankIDCode:            "GBDSC",
    Bic:                   "NWBKGB22",
    AccountClassification: "Personal",
}
// >> using the NewAccount makes sure that the Account is properly initialised
act := f3c.NewAccount(uid, oid, attr)

// >> finally use AccountSvc to create the Account in Form3 remote system
svc := http.AccountSvc{
            Base: "http://form3_api_service",
        }
svc.Create(act)
```

# Fetch
Fetching an existing account can be done via `AccountSvc.Fetch()` function
```
svc := http.AccountSvc{
    Base: "http://form3_api_service",
}
id := "id_of_target_account"
if act, err := svc.Fetch(id); err != nil {
    t.Fatal(err)
}
```

# Delete
Deleting an existing account can be done via `AccountSvc.Delete()` function
```
svc := http.AccountSvc{
    Base: "http://form3_api_service",
}
id := "id_of_target_account"
version := 0
if act, err := svc.Delete(id, version); err != nil {
    t.Fatal(err)
}
```

# List
Listing existing accounts can be done via `AccountSvc.Delete()` function
```
svc := http.AccountSvc{
    Base: "http://form3_api_service",
}
// get all accounts
acts, res := []f3c.AccountXL{}, []f3c.AccountXL{}
var err error
for pg := 0; true; pg += 1 {
    if res, err = svc.List(pg, 100); err != nil {
        t.Fatal(err)
    }
    if len(res) == 0 {
        break
    }
    acts = append(acts, res...)
}
```

# Full Reference
See [documentation](./doc.md)