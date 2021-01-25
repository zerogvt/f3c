package http_test

import (
	"testing"

	"github.com/zerogvt/f3client"
	"github.com/zerogvt/f3client/http"
)

func TestAccountSvc_Create(t *testing.T) {
	t.Run("Account created OK", func(t *testing.T) {
		var svc http.AccountSvc
		uid := "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"
		oid := "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c"
		attr := f3client.Attributes{
			Country:    "GB",
			BaseCurr:   "GBP",
			BankID:     "400300",
			BankIDCode: "GBDSC",
			BIC:        "NWBKGB22",
		}
		act := f3client.Account{
			ID:         uid,
			OrgID:      oid,
			Attributes: attr,
		}
		if act, err := svc.Create(&act); err != nil {
			t.Fatal(err)
		} else if uid != act.ID {
			t.Fatal("UID expected", uid, "but got", act.ID)
		} else if oid != act.OrgID {
			t.Fatal("UID expected", oid, "but got", act.OrgID)
		} else if attr != act.Attributes {
			t.Fatal("Attributes expected", attr, "but got",
				act.Attributes)
		}
	})
}