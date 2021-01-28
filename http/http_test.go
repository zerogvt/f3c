package http_test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/zerogvt/f3c"
	"github.com/zerogvt/f3c/http"
)

func randomID(length int) string {
	seed := time.Now().Unix()
	fmt.Printf("seed: %d\n", seed)
	rand.Seed(seed)
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
			fmt.Println("1111")
			svc := http.AccountSvc{
				Base: "http://localhost:8080",
			}
			// get a fresh uid to avoid conflicts with past accounts
			uid := "ad27e265-9605-4b4b-a0e5-" + randomID(12)
			fmt.Println(uid)
			oid := "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c"
			attr := f3c.Attributes{
				Country:      "GB",
				BaseCurrency: "GBP",
				BankID:       "400300",
				BankIDCode:   "GBDSC",
				Bic:          "NWBKGB22",
			}
			act := f3c.Account{
				Type:           "accounts",
				ID:             uid,
				OrganisationID: oid,
				Attributes:     attr,
			}
			if res, err := svc.Create(act); err != nil {
				t.Fatal(err)
			} else {
				pprint(res)
				if uid != res.Data.ID {
					t.Fatal("UID expected", uid, "but got", res.Data.ID)
				} else if oid != res.Data.OrganisationID {
					t.Fatal("Org ID expected", oid, "but got",
						res.Data.OrganisationID)
				} else if attr.Country != res.Data.Attributes.Country {
					t.Fatal("Country expected", attr.Country, "but got",
						res.Data.Attributes.Country)
				} else if attr.BaseCurrency != res.Data.Attributes.BaseCurrency {
					t.Fatal("BaseCurrency expected", attr.BaseCurrency, "but got",
						res.Data.Attributes.BaseCurrency)
				} else if attr.BankID != res.Data.Attributes.BankID {
					t.Fatal("BankID expected", attr.BankID, "but got",
						res.Data.Attributes.BankID)
				} else if attr.BankIDCode != res.Data.Attributes.BankIDCode {
					t.Fatal("BankIDCode expected", attr.BankIDCode, "but got",
						res.Data.Attributes.BankIDCode)
				} else if attr.Bic != res.Data.Attributes.Bic {
					t.Fatal("BankIDCode expected", attr.Bic, "but got",
						res.Data.Attributes.Bic)
				}
			}
		})
}

func TestAccountSvc_CreateDuplicate(t *testing.T) {
	t.Run("We should catch an HTTP error such as duplicate account creation",
		func(t *testing.T) {
			fmt.Println("222222")
			svc := http.AccountSvc{
				Base: "http://localhost:8080",
			}
			uid := "ad27e265-9605-4b4b-a0e5-" + randomID(12)
			fmt.Println(uid)
			oid := "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c"
			attr := f3c.Attributes{
				Country:      "GB",
				BaseCurrency: "GBP",
				BankID:       "400300",
				BankIDCode:   "GBDSC",
				Bic:          "NWBKGB22",
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
