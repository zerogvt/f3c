package http_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	nethttp "net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/zerogvt/f3c"
	"github.com/zerogvt/f3c/http"
)

var Server = "http://localhost:8080"
var ServerEnvKey = "TEST_SERVER"

func init() {
	rand.Seed(time.Now().Unix())
	if _, ok := os.LookupEnv(ServerEnvKey); ok {
		Server = os.Getenv(ServerEnvKey)
	}
	fmt.Println("Tests will be run against server: ", Server)
}

func TestCreate(t *testing.T) {
	t.Run("A new account should be created without errors",
		func(t *testing.T) {
			svc := http.AccountSvc{
				Base: Server,
			}
			// Get a fresh uid to avoid conflicts with past accounts.
			// We don't use the createRandomAcct() helper func here because we
			// want fine grained control on the creation process for this test.
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
			act := f3c.NewAccount(uid, oid, attr)
			res := f3c.AccountXL{}
			var err error
			if res, err = svc.Create(act); err != nil {
				t.Fatal(err)
			}
			if err = isEqual(act, res.Account); err != nil {
				t.Fatal(err)
			}
		})
}

func TestCreateDuplicate(t *testing.T) {
	t.Run("We should catch an HTTP error such as duplicate account creation",
		func(t *testing.T) {
			svc := http.AccountSvc{
				Base: Server,
			}
			// Create first account.
			act_1 := createRandomAcct(t, &svc)
			// When trying to recreate the same account
			// we should get an error.
			act_2 := act_1
			if _, err := svc.Create(act_2); err == nil {
				t.Fatal("no error was produced")
			}
		})
}

func TestCreateWithNonStandardClient(t *testing.T) {
	t.Run("We should catch an HTTP error such as duplicate account creation",
		func(t *testing.T) {
			svc := http.AccountSvc{
				Cli: nethttp.Client{
					Timeout: time.Duration(10) * time.Second,
				},
				Base: Server,
			}
			// Create first account.
			createRandomAcct(t, &svc)
		})
}

func TestFetch(t *testing.T) {
	t.Run("An existing account should be fetched",
		func(t *testing.T) {
			svc := http.AccountSvc{
				Base: Server,
			}
			act := createRandomAcct(t, &svc)
			// now try to fetch
			if _, err := svc.Fetch(act.ID); err != nil {
				t.Fatal(err)
			}
		})
}

func TestDelete(t *testing.T) {
	t.Run("An existing account should be deletable",
		func(t *testing.T) {
			svc := http.AccountSvc{
				Base: Server,
			}
			act := createRandomAcct(t, &svc)
			// now try to delete
			if err := svc.Delete(act.ID, 0); err != nil {
				t.Fatal(err)
			}
		})
}

func TestList(t *testing.T) {
	t.Run("A list of 5 newly created accounts should be listed",
		func(t *testing.T) {
			svc := http.AccountSvc{
				Base: Server,
			}
			// create 5 accounts
			want := []string{}
			for i := 1; i <= 5; i++ {
				want = append(want, createRandomAcct(t, &svc).ID)
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
			// create a map with all uids
			uids := make(map[string]bool)
			for _, ac := range acts {
				uids[ac.ID] = true
			}
			// we must have our 5 newly created ids in that map
			for _, id := range want {
				if _, ok := uids[id]; !ok {
					t.Fatal("uid ", id, "not found in listed accounts.")
				}
			}
		})
}

func TestToPayload(t *testing.T) {
	t.Run("An Account should be marshalled in a payload",
		func(t *testing.T) {
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
			want := f3c.NewAccount(uid, oid, attr)
			// convert our account to a payload
			payload := http.ToPayload(want)
			// and unmarshall the payload
			have := f3c.PayloadIn{}
			if err := json.Unmarshal(payload.Bytes(), &have); err != nil {
				t.Fatal(err)
			}
			// the account in unmarshalled payload should match the initial one
			if err := isEqual(want, have.Account.Account); err != nil {
				t.Fatal(err)
			}
		})
}

func TestFromPayload(t *testing.T) {
	t.Run("A valid reply should be unmarshalled in a valid account",
		func(t *testing.T) {
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
			act := f3c.NewAccount(uid, oid, attr)
			// simulate a reply payload
			want := f3c.PayloadIn{
				Account: f3c.AccountXL{
					Account:       act,
					Version:       0,
					Relationships: f3c.Relationships{},
				},
			}
			// we're going to need the body as a io.Reader
			var body *bytes.Buffer
			if data, err := json.Marshal(want); err != nil {
				t.Fatal(err)
			} else {
				body = bytes.NewBuffer([]byte(data))
			}
			// so that we can wrap it in Nopcloser and pass it for a io.ReadCloser
			mresp := nethttp.Response{
				Status:     "200 OK",
				StatusCode: 200,
				Body:       ioutil.NopCloser(body),
			}
			// now run this mock reply through FromPayload()
			have := f3c.AccountXL{}
			var err error
			if have, err = http.FromPayload(&mresp); err != nil {
				t.Fatal(err)
			}
			// the account in the unmarshalled payload should match the initial one
			if err := isEqual(want.Account.Account, have.Account); err != nil {
				t.Fatal(err)
			}
		})
}

func actErr(name string, field1 interface{}, field2 interface{}) error {
	return errors.New(fmt.Sprintf("ERROR %s: %v != %v", name, field1, field2))
}

func isEqual(act f3c.Account, res f3c.Account) error {
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

// randomID returns a random string of length length
// helper func
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

// createRandomAcct creates a random account
// helper func
func createRandomAcct(t *testing.T, svc *http.AccountSvc) f3c.Account {
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
	act := f3c.NewAccount(uid, oid, attr)
	if _, err := svc.Create(act); err != nil {
		t.Fatal(err)
	}
	return act
}
