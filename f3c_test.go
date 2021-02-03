package f3c_test

import (
	"testing"

	"github.com/zerogvt/f3c"
)

func TestNewAccount(t *testing.T) {
	t.Run("A new account should be initialised properly",
		func(t *testing.T) {
			uid := "someid"
			oid := "someorgid"
			attr := f3c.Attributes{
				Country:               "GB",
				BaseCurrency:          "GBP",
				BankID:                "400300",
				BankIDCode:            "GBDSC",
				Bic:                   "NWBKGB22",
				AccountClassification: "Personal",
			}
			act := f3c.NewAccount(uid, oid, attr)
			if act.Type != "accounts" {
				t.Fatal("Account type not set to 'Accounts'")
			}
		})
}
