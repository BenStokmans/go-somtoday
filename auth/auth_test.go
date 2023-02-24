package auth

import (
	"github.com/sirupsen/logrus"
	"testing"
)

func TestAuthFlow(t *testing.T) {
	logger := logrus.New()

	ctx := NewContext(logger).WithDefaultCache()
	err := ctx.SetOrganisationByID("abcce4ce-6f96-40b3-9185-7ab043ac3a94")
	if err != nil {
		t.Fatal(err)
	}
	err = ctx.SetOIDCurl(ctx.Organisation.OID[0].URL)
	if err != nil {
		t.Fatal(err)
	}
	err = ctx.DoAuth()
	if err != nil {
		t.Fatal(err)
	}
	err = ctx.TokenValid()
	if err != nil {
		t.Fatal(err)
	}
}

func TestRandChallenge(t *testing.T) {
	s1, err := getVerifier()
	if err != nil {
		t.Fatal(err)
	}
	s2, err := getVerifier()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("s1: %s, s2: %s", s1, s2)
	if s1 == s2 {
		t.Fatal("strings equal")
	}
}
