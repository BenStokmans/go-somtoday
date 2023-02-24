package somtoday

import (
	"encoding/json"
	"fmt"
	"github.com/BenStokmans/go-somtoday/auth"
	"strconv"
)

type SchoolYear struct {
	ID    int    `json:"-"`
	Type  string `json:"$type"`
	Links []struct {
		Id   int    `json:"id"`
		Rel  string `json:"rel"`
		Type string `json:"type"`
		Href string `json:"href"`
	} `json:"links"`
	Permissions []struct {
		Full       string   `json:"full"`
		Type       string   `json:"type"`
		Operations []string `json:"operations"`
		Instances  []string `json:"instances"`
	} `json:"permissions"`
	AdditionalObjects struct {
	} `json:"additionalObjects"`
	Name      string `json:"naam"`
	FromDate  string `json:"vanafDatum"`
	UntilDate string `json:"totDatum"`
	IsCurrent bool   `json:"isHuidig"`
}

const schoolYearUrl = "/rest/v1/schooljaren"

func GetCurrentSchoolYear(ctx *auth.Context) (SchoolYear, error) {
	var year SchoolYear

	body, err := ctx.DoAuthedGet(ctx.Token.SomtodayAPIURL+schoolYearUrl+"/huidig", nil)
	err = json.Unmarshal(body, &year)
	if err != nil {
		return year, fmt.Errorf("error unmarshalling school year: %v", err)
	}

	for _, link := range year.Links {
		if link.Rel == "self" {
			year.ID = link.Id
		}
	}
	return year, nil
}

func GetSchoolYearById(id int, ctx *auth.Context) (SchoolYear, error) {
	var year SchoolYear

	body, err := ctx.DoAuthedGet(ctx.Token.SomtodayAPIURL+schoolYearUrl+"/"+strconv.Itoa(id), nil)
	err = json.Unmarshal(body, &year)
	if err != nil {
		return year, fmt.Errorf("error unmarshalling school year: %v", err)
	}

	year.ID = id
	return year, nil
}
