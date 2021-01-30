package http_test

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/zerogvt/f3c"
	"github.com/zerogvt/f3c/http"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func randomID(length int) string {
	//Only lowercase
	charSet := "0123456789abcdedf"
	var output strings.Builder
	for i := 0; i < length; i++ {
		random := rand.Intn(len(charSet))
		randomChar := charSet[random]
		output.WriteString(string(randomChar))
	}
	return output.String()
}

func TestAccountSvc_Create(t *testing.T) {
	t.Run("A new account should be created without errors",
		func(t *testing.T) {
			svc := http.AccountSvc{
				Base: "http://localhost:8080",
			}
			// get a fresh uid to avoid conflicts with past accounts
			uid := "ad27e265-9605-4b4b-a0e5-" + randomID(12)
			oid := "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c"
			attr := f3c.Attributes{
				Country:               "GB",
				BaseCurrency:          "GBP",
				BankID:                "400300",
				BankIDCode:            "GBDSC",
				Bic:                   "NWBKGB22",
				AccountClassification: "Personal",
			}
			act := f3c.Account{
				Type:           "accounts",
				ID:             uid,
				OrganisationID: oid,
				Attributes:     attr,
			}
			res := f3c.AccountXL{}
			var err error
			if res, err = svc.Create(act); err != nil {
				t.Fatal(err)
			}
			f3c.Pprint(res)
			if err = isEqual(act, res); err != nil {
				t.Fatal(err)
			}
		})
}

func TestAccountSvc_CreateDuplicate(t *testing.T) {
	t.Run("We should catch an HTTP error such as duplicate account creation",
		func(t *testing.T) {
			svc := http.AccountSvc{
				Base: "http://localhost:8080",
			}
			uid := "ad27e265-9605-4b4b-a0e5-" + randomID(12)
			oid := "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c"
			attr := f3c.Attributes{
				Country:               "GB",
				BaseCurrency:          "GBP",
				BankID:                "400300",
				BankIDCode:            "GBDSC",
				Bic:                   "NWBKGB22",
				AccountClassification: "Personal",
			}
			act_1 := f3c.Account{
				Type:           "accounts",
				ID:             uid,
				OrganisationID: oid,
				Attributes:     attr,
			}
			act_2 := act_1
			if _, err := svc.Create(act_1); err != nil {
				t.Fatal(err)
			}
			if _, err := svc.Create(act_2); err == nil {
				t.Fatal("no error was produced")
			} else {
				fmt.Print(err.Error())
			}
		})
}

func TestAccountSvc_Fetch(t *testing.T) {
	t.Run("An existing account should be fetched",
		func(t *testing.T) {
			svc := http.AccountSvc{
				Base: "http://localhost:8080",
			}
			// make sure the test account is in server
			uid := "ad27e265-9605-4b4b-a0e5-000000000000"
			oid := "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c"
			attr := f3c.Attributes{
				Country:               "GB",
				BaseCurrency:          "GBP",
				BankID:                "400300",
				BankIDCode:            "GBDSC",
				Bic:                   "NWBKGB22",
				AccountClassification: "Personal",
			}
			act := f3c.Account{
				Type:           "accounts",
				ID:             uid,
				OrganisationID: oid,
				Attributes:     attr,
			}
			if _, err := svc.Create(act); err != nil {
				//it's ok to get a conflict error
				if !strings.Contains(err.Error(), "Conflict") {
					t.Fatal(err)
				}
			}
			// now try to fetch
			res := f3c.AccountXL{}
			var err error
			if res, err = svc.Fetch(uid); err != nil {
				t.Fatal(err)
			}
			if err = isEqual(act, res); err != nil {
				t.Fatal(err)
			}
		})
}

func TestAccountSvc_List(t *testing.T) {
	t.Run("An existing account should be fetched",
		func(t *testing.T) {
			svc := http.AccountSvc{
				Base: "http://localhost:8080",
			}
			acts := []f3c.AccountXL{}
			var err error
			if acts, err = svc.List(0, 100); err != nil {
				t.Fatal(err)
			}
			f3c.Pprint(acts)
		})
}

func actErr(name string, field1 interface{}, field2 interface{}) error {
	return errors.New(fmt.Sprintf("ERROR %s: %v != %v", name, field1, field2))
}

func isEqual(act f3c.Account, res f3c.AccountXL) error {
	if act.ID != res.ID {
		return actErr("UID", act.ID, res.ID)
	}
	if act.OrganisationID != res.OrganisationID {
		return actErr("Org ID", act.OrganisationID,
			res.OrganisationID)
	}
	if act.Attributes.Country != res.Attributes.Country {
		return actErr("Country", act.Attributes.Country,
			res.Attributes.Country)
	}
	if act.Attributes.BaseCurrency != res.Attributes.BaseCurrency {
		return actErr("BaseCurrency", act.Attributes.BaseCurrency,
			res.Attributes.BaseCurrency)
	}
	if act.Attributes.BankID != res.Attributes.BankID {
		return actErr("BankID", act.Attributes.BankID,
			res.Attributes.BankID)
	}
	if act.Attributes.BankIDCode != res.Attributes.BankIDCode {
		return actErr("BankIDCode", act.Attributes.BankIDCode,
			res.Attributes.BankIDCode)
	}
	if act.Attributes.Bic != res.Attributes.Bic {
		return actErr("BIC", act.Attributes.Bic,
			res.Attributes.Bic)
	}
	return nil
}
