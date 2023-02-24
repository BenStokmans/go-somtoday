package somtoday

import (
	"encoding/json"
	"fmt"
	"github.com/BenStokmans/go-somtoday/auth"
	"time"
)

type Result struct {
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
		CurrentOtherSubjectColumn struct {
			Type  string `json:"$type"`
			Items []struct {
				Type  string `json:"$type"`
				Links []struct {
					Id   int64  `json:"id"`
					Rel  string `json:"rel"`
					Type string `json:"type"`
				} `json:"links"`
				Permissions       []interface{} `json:"permissions"`
				AdditionalObjects struct {
				} `json:"additionalObjects"`
				OtherSubject      Subject `json:"anderVak"`
				Weight            int     `json:"weging"`
				ExamWeight        int     `json:"examenWeging"`
				InProgressDossier bool    `json:"inVoortgangsdossier"`
				InExamDossier     bool    `json:"inExamendossier"`
				Batch             `json:"lichting"`
			} `json:"items"`
		} `json:"huidigeAnderVakKolommen"`
		CompositeExamColumnID interface{} `json:"samengesteldeToetskolomId"`
		ResultColumnID        int64       `json:"resultaatkolomId"`
		KnownRapportCardGrade *struct {
			Type          string `json:"$type"`
			RapportCijfer string `json:"rapportCijfer"`
		} `json:"berekendRapportCijfer"`
		ExamTypeName  *string `json:"toetssoortnaam"`
		GradeColumnID *int64  `json:"cijferkolomId"`
	} `json:"additionalObjects"`
	ResitType               string    `json:"herkansingstype"`
	EntryDate               time.Time `json:"datumInvoer"`
	DoesNotCount            bool      `json:"teltNietmee"`
	ExamNotDone             bool      `json:"toetsNietGemaakt"`
	GradeNumber             int       `json:"leerjaar"`
	Period                  int       `json:"periode"`
	ExamWeight              int       `json:"examenWeging,omitempty"`
	IsExamDossierResult     bool      `json:"isExamendossierResultaat"`
	IsProgressDossierResult bool      `json:"isVoortgangsdossierResultaat"`
	Type1                   string    `json:"type"`
	Subject                 `json:"vak"`
	Student                 `json:"leerling"`
	Exemption               bool   `json:"vrijstelling"`
	Result                  string `json:"resultaat,omitempty"`
	ValidResult             string `json:"geldendResultaat,omitempty"`
	Weight                  int    `json:"weging,omitempty"`
	Description             string `json:"omschrijving,omitempty"`
	TrackingNumber          int    `json:"volgnummer,omitempty"`
	ResultLabel             string `json:"resultaatLabel,omitempty"`
	ResultLabelAbbreviation string `json:"resultaatLabelAfkorting,omitempty"`
	Resit                   struct {
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
			HuidigeAnderVakKolommen struct {
				Type  string        `json:"$type"`
				Items []interface{} `json:"items"`
			} `json:"huidigeAnderVakKolommen"`
			SamengesteldeToetskolomId interface{} `json:"samengesteldeToetskolomId"`
			ResultaatkolomId          int64       `json:"resultaatkolomId"`
			BerekendRapportCijfer     interface{} `json:"berekendRapportCijfer"`
			Toetssoortnaam            string      `json:"toetssoortnaam"`
			CijferkolomId             interface{} `json:"cijferkolomId"`
		} `json:"additionalObjects"`
		ResitType               string    `json:"herkansingstype"`
		ResitNumber             int       `json:"herkansingsNummer"`
		Result                  string    `json:"resultaat"`
		EntryDate               time.Time `json:"datumInvoer"`
		DoesNotCount            bool      `json:"teltNietmee"`
		ExamNotDone             bool      `json:"toetsNietGemaakt"`
		GradeNumber             int       `json:"leerjaar"`
		Period                  int       `json:"periode"`
		Weight                  int       `json:"weging"`
		ExamWeight              int       `json:"examenWeging"`
		IsExamDossierResult     bool      `json:"isExamendossierResultaat"`
		IsProgressDossierResult bool      `json:"isVoortgangsdossierResultaat"`
		Type                    string    `json:"type"`
		Description             string    `json:"omschrijving"`
		Subject                 `json:"vak"`
		Student                 `json:"leerling"`
		TrackingNumber          int  `json:"volgnummer"`
		Exemption               bool `json:"vrijstelling"`
	} `json:"herkansing,omitempty"`
	ResitNumber int `json:"herkansingsNummer,omitempty"`
}

const resultUrl = "/rest/v1/resultaten/huidigVoorLeerling/%d?additional=berekendRapportCijfer&additional=samengesteldeToetskolomId&additional=resultaatkolomId&additional=cijferkolomId&additional=toetssoortnaam&additional=huidigeAnderVakKolommen "

func GetResults(studentId int64, ctx *auth.Context) ([]Result, error) {
	rangeL, rangeH := 0, 100
	var results []Result
	for {
		s, err := getResultsRange(studentId, ctx, rangeL, rangeH)
		if err != nil {
			return nil, err
		}
		results = append(results, s...)
		if len(s) < 100 {
			break
		}
		rangeL, rangeH = rangeH, rangeH+100
	}
	return results, nil
}

func getResultsRange(studentId int64, ctx *auth.Context, rangeL, rangeH int) ([]Result, error) {
	reqUrl := fmt.Sprintf(ctx.Token.SomtodayAPIURL+resultUrl, studentId)
	headers := map[string][]string{
		"Range": {fmt.Sprintf("items=%d-%d", rangeL, rangeH)},
	}
	body, err := ctx.DoAuthedGet(reqUrl, headers)

	var results struct {
		Results []Result `json:"items"`
	}
	err = json.Unmarshal(body, &results)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling results: %v", err)
	}
	return results.Results, nil
}
