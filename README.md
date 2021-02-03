[![Build Status](https://travis-ci.com/zerogvt/f3c.svg?token=zDwdCt1iLHUMB3eVQDMy&branch=main)](https://travis-ci.com/github/zerogvt/f3c)

# Client library for Form3 API
as per https://github.com/form3tech-oss/interview-accountapi

# Import
`import "github.com/zerogvt/f3c"`

# Account creation
```
    import (
        "github.com/zerogvt/f3c"
        "github.com/zerogvt/f3c/http"
    )

    // create a local Account instance...
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
    // ...using the NewAccount makes sure that the Account is properly initialised
    act := f3c.NewAccount(uid, oid, attr)
    // and finally use AccountSvc to create the Account in Form3 remote system
    svc := http.AccountSvc{
				Base: "http://form3_api_service",
			}
    svc.Create(act)
```


# Run all-in-one tests
`make tests`
