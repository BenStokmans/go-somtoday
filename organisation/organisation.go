package organisation

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Organization struct {
	UUID string `json:"uuid"`
	Name string `json:"naam"`
	City string `json:"plaats"`
	OID  []struct {
		Description string `json:"omschrijving"`
		URL         string `json:"url"`
		DomainHint  string `json:"domain_hint"`
	} `json:"oidcurls"`
}

const orgUrl = "https://servers.somtoday.nl/organisaties.json"

var orgCache []Organization

func getOrganizations() ([]Organization, error) {
	if orgCache != nil {
		return orgCache, nil
	}

	resp, err := http.Get(orgUrl)
	if err != nil {
		return nil, fmt.Errorf("error retrieving all organizations: %v", err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading body: %v", err)
	}
	var organizations []struct {
		Settings []Organization `json:"instellingen"`
	}
	err = json.Unmarshal(body, &organizations)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling organizations: %v", err)
	}
	return organizations[0].Settings, nil
}

func GetOrganizationByID(tenantId string) (Organization, error) {
	orgs, err := getOrganizations()
	if err != nil {
		return Organization{}, nil
	}
	for _, o := range orgs {
		if o.UUID == tenantId {
			return o, nil
		}
	}
	return Organization{}, errors.New("organization not found")
}

func GetOrganizationByName(name string) (Organization, error) {
	orgs, err := getOrganizations()
	if err != nil {
		return Organization{}, nil
	}
	for _, o := range orgs {
		if o.Name == name {
			return o, nil
		}
	}
	return Organization{}, errors.New("organization not found")
}
