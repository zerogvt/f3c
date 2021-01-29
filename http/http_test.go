package http_test

import (
	"encoding/json"
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

func pprint(data interface{}) {
	json, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return
	}
	fmt.Println(string(json))
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
				Country:      "GB",
				BaseCurrency: "GBP",
				BankID:       "400300",
				BankIDCode:   "GBDSC",
				BIC:          "NWBKGB22",
			}
			act := f3c.Account{
				Type:           "accounts",
				ID:             uid,
				OrganisationID: oid,
				Attributes:     attr,
			}
			res := f3c.AccountCrResp{}
			var err error
			if res, err = svc.Create(act); err != nil {
				t.Fatal(err)
			}
			pprint(res)
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
				Country:      "GB",
				BaseCurrency: "GBP",
				BankID:       "400300",
				BankIDCode:   "GBDSC",
				BIC:          "NWBKGB22",
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
				Country:      "GB",
				BaseCurrency: "GBP",
				BankID:       "400300",
				BankIDCode:   "GBDSC",
				BIC:          "NWBKGB22",
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
			res := f3c.AccountCrResp{}
			var err error
			if res, err = svc.Fetch(uid); err != nil {
				t.Fatal(err)
			}
			if err = isEqual(act, res); err != nil {
				t.Fatal(err)
			}
		})
}

func actErr(name string, field1 interface{}, field2 interface{}) error {
	return errors.New(fmt.Sprintf("%s: %v != %v", name, field1, field2))
}

func isEqual(act f3c.Account, res f3c.AccountCrResp) error {
	if act.ID != res.Data.ID {
		return actErr("UID", act.ID, res.Data.ID)
	}
	if act.OrganisationID != res.Data.OrganisationID {
		return actErr("Org ID", act.OrganisationID,
			res.Data.OrganisationID)
	}
	if act.Attributes.Country != res.Data.Attributes.Country {
		return actErr("Country", act.Attributes.Country,
			res.Data.Attributes.Country)
	}
	if act.Attributes.BaseCurrency != res.Data.Attributes.BaseCurrency {
		return actErr("BaseCurrency", act.Attributes.BaseCurrency,
			res.Data.Attributes.BaseCurrency)
	}
	if act.Attributes.BankID != res.Data.Attributes.BankID {
		return actErr("BankID", act.Attributes.BankID,
			res.Data.Attributes.BankID)
	}
	if act.Attributes.BankIDCode != res.Data.Attributes.BankIDCode {
		return actErr("BankIDCode", act.Attributes.BankIDCode,
			res.Data.Attributes.BankIDCode)
	}
	if act.Attributes.BIC != res.Data.Attributes.BIC {
		return actErr("BIC", act.Attributes.BIC,
			res.Data.Attributes.BIC)
	}
	return nil
}
