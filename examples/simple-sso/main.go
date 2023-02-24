package main

import (
	"github.com/BenStokmans/go-somtoday/auth"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()

	// create a new context using the cache directory of the current user
	ctx := auth.NewContext(logger).WithDefaultCache()

	// set the organisation of the current auth context (this can also be done by name)
	// for a list of organisations see https://servers.somtoday.nl/organisaties.json
	err := ctx.SetOrganisationByID("abcce4ce-6f96-40b3-9185-7ab043ac3a94")
	if err != nil {
		logrus.Fatal(err)
	}

	// pick an OID url for our SSO flow
	err = ctx.SetOIDCurl(ctx.Organisation.OID[0].URL)
	if err != nil {
		logrus.Fatal(err)
	}

	// load our token from the cache or proceed to the SSO flow
	err = ctx.LoadOrDoSSOAuth()
	if err != nil {
		logrus.Fatal(err)
	}

	// check if the token is valid
	err = ctx.TokenValid()
	if err != nil {
		logrus.Fatal(err)
	}

	// do anything with this auth context
}
