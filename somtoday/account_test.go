package somtoday

import (
	"github.com/BenStokmans/go-somtoday/auth"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestGetAccount(t *testing.T) {
	logger := logrus.New()

	ctx := auth.NewContext(logger).WithDefaultCache()
	err := ctx.SetOrganisationByID("abcce4ce-6f96-40b3-9185-7ab043ac3a94")
	if err != nil {
		t.Fatal(err)
	}
	err = ctx.SetOIDCurl(ctx.Organisation.OID[0].URL)
	if err != nil {
		t.Fatal(err)
	}
	err = ctx.LoadOrDoAuth()
	if err != nil {
		t.Fatal(err)
	}
	choices, err := GetAppointments("2023-02-13", "2023-03-18", ctx)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(choices[3].Description)
}
