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

func createRandomAcct(t *testing.T) f3c.Account {
	svc := http.AccountSvc{
		Base: "http://localhost:8080",
	}
	// make sure the test account is in server
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
	if _, err := svc.Create(act); err != nil {
		//it's ok to get a conflict error
		if !strings.Contains(err.Error(), "Conflict") {
			t.Fatal(err)
		}
	}
	return act
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
			act := createRandomAcct(t)
			// now try to fetch
			if _, err := svc.Fetch(act.ID); err != nil {
				t.Fatal(err)
			}
		})
}

func TestAccountSvc_Delete(t *testing.T) {
	t.Run("An existing account should be deletable",
		func(t *testing.T) {
			svc := http.AccountSvc{
				Base: "http://localhost:8080",
			}
			act := createRandomAcct(t)
			// now try to delete
			if err := svc.Delete(act.ID, 0); err != nil {
				t.Fatal(err)
			}
		})
}

func TestAccountSvc_List(t *testing.T) {
	t.Run("A list of 5 newly created accounts should be listed",
		func(t *testing.T) {
			svc := http.AccountSvc{
				Base: "http://localhost:8080",
			}
			// create 5 accounts
			want := []string{}
			for i := 1; i <= 5; i++ {
				want = append(want, createRandomAcct(t).ID)
			}
			acts, res := []f3c.AccountXL{}, []f3c.AccountXL{}
			var err error
			// get all accounts
			for pg := 0; true; pg += 1 {
				if res, err = svc.List(pg, 100); err != nil {
					t.Fatal(err)
				}
				if len(res) == 0 {
					break
				}
				acts = append(acts, res...)
			}
			t.Log("I got a total of", len(acts), "accounts")
			// create a map with all uids
			uids := make(map[string]bool)
			for _, ac := range acts {
				uids[ac.ID] = true
			}
			// we must have our 5 newly created ids in that map
			for _, id := range want {
				if _, ok := uids[id]; !ok {
					t.Fatal("uid ", id, "not found in listed accounts.")
				} else {
					t.Log("uid ", id, "found in listed accounts.")
				}
			}
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
