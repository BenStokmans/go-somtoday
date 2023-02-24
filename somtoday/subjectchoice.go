package somtoday

import (
	"encoding/json"
	"fmt"
	"github.com/BenStokmans/go-somtoday/auth"
	"time"
)

type SubjectChoice struct {
	Type  string `json:"$type"`
	Links []struct {
		Id   int64  `json:"id"`
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
		SubjectScoring *struct {
			Type         string `json:"$type"`
			SubjectID    int64  `json:"vakId"`
			ExamScoring1 string `json:"toetsnormering1"`
			ExamScoring2 string `json:"toetsnormering2"`
		} `json:"vaknormering"`
	} `json:"additionalObjects"`
	Subject struct {
		Links []struct {
			Id   int64  `json:"id"`
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
		Abbreviation string `json:"afkorting"`
		Name         string `json:"naam"`
	} `json:"vak"`
	Student            `json:"leerling"`
	Exemption          bool `json:"vrijstelling"`
	Batch              `json:"lichting"`
	RelevantGradeBatch Batch `json:"relevanteCijferLichting"`
}

const subjectChoiceUrl = "/rest/v1/vakkeuzes?additional=vaknormering"

func GetActiveSubjectChoices(ctx *auth.Context) ([]SubjectChoice, error) {
	return getSubjectChoices(
		ctx.Token.SomtodayAPIURL+
			subjectChoiceUrl+
			"&actiefOpPeildatum="+
			time.Now().Format("2006-01-02T15:04:05.000"), ctx)
}

func GetSubjectChoices(ctx *auth.Context) ([]SubjectChoice, error) {
	return getSubjectChoices(ctx.Token.SomtodayAPIURL+subjectChoiceUrl, ctx)
}

func getSubjectChoices(url string, ctx *auth.Context) ([]SubjectChoice, error) {
	rangeL, rangeH := 0, 100
	var choices []SubjectChoice
	for {
		s, err := getSubjectChoiceRange(url, ctx, rangeL, rangeH)
		if err != nil {
			return nil, err
		}
		choices = append(choices, s...)
		if len(s) < 100 {
			break
		}
		rangeL, rangeH = rangeH, rangeH+100
	}
	return choices, nil
}

func getSubjectChoiceRange(url string, ctx *auth.Context, rangeL, rangeH int) ([]SubjectChoice, error) {
	headers := map[string][]string{
		"Range": {fmt.Sprintf("items=%d-%d", rangeL, rangeH)},
	}
	body, err := ctx.DoAuthedGet(url, headers)

	var choices struct {
		Choices []SubjectChoice `json:"items"`
	}
	err = json.Unmarshal(body, &choices)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling subject choices: %v", err)
	}
	return choices.Choices, nil
}
