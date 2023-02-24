package somtoday

import (
	"encoding/json"
	"fmt"
	"github.com/BenStokmans/go-somtoday/auth"
)

type Account struct {
	ID    int64 `json:"-"`
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
		Restrictions struct {
			Type  string `json:"$type"`
			Items []struct {
				Type                       string        `json:"$type"`
				Links                      []interface{} `json:"links"`
				Permissions                []interface{} `json:"permissions"`
				AdditionalObjects          struct{}      `json:"additionalObjects"`
				LocationId                 int64         `json:"vestigingsId"`
				MobileAppEnabled           bool          `json:"mobieleAppAan"`
				StudyGuideEnabled          bool          `json:"studiewijzerAan"`
				SendMessagesEnabled        bool          `json:"berichtenVerzendenAan"`
				LearningResourcesEnabled   bool          `json:"leermiddelenAan"`
				AdviceTokenEnabled         bool          `json:"adviezenTokenAan"`
				ShowRapportCardGradeRemark bool          `json:"opmerkingRapportCijferTonenAan"`
				ShowPeriodAverage          bool          `json:"periodeGemiddeldeTonenResultaatAan"`
				ShowRapportCardAverage     bool          `json:"rapportGemiddeldeTonenResultaatAan"`
				ShowRapportCardGrade       bool          `json:"rapportCijferTonenResultaatAan"`
				ExamTypeAveragesEnabled    bool          `json:"toetssoortgemiddeldenAan"`
				SchoolExamResultsEnabled   bool          `json:"seResultaatAan"`
				CoreYearGroupEnabled       bool          `json:"stamgroepLeerjaarAan"`
				CanChangeEmail             bool          `json:"emailWijzigenAan"`
				CanChangeMobile            bool          `json:"mobielWijzigenAan"`
				CanChangePassword          bool          `json:"wachtwoordWijzigenAan"`
				SeeAbsences                bool          `json:"absentiesBekijkenAan"`
				SeeAbsenceObservation      bool          `json:"absentieConstateringBekijkenAan"`
				SeeAbsenceConsequence      bool          `json:"absentieMaatregelBekijkenAan"`
				SeeAbsenceNotification     bool          `json:"absentieMeldingBekijkenAan"`
				SeeMessages                bool          `json:"berichtenBekijkenAan"`
				SeeGrades                  bool          `json:"cijfersBekijkenAan"`
				SeeHomework                bool          `json:"huiswerkBekijkenAan"`
				SeeNews                    bool          `json:"nieuwsBekijkenAan"`
				ShowStudentPhoto           bool          `json:"pasfotoLeerlingTonenAan"`
				ShowStaffPhoto             bool          `json:"pasfotoMedewerkerTonenAan"`
				SeeProfile                 bool          `json:"profielBekijkenAan"`
				SeeTimetable               bool          `json:"roosterBekijkenAan"`
				SeeSubjects                bool          `json:"vakkenBekijkenAan"`
				HideClassTimes             bool          `json:"lesurenVerbergenSettingAan"`
			} `json:"items"`
		} `json:"restricties"`
	} `json:"additionalObjects,omitempty"`
	UserName           string        `json:"gebruikersnaam"`
	AccountPermissions []interface{} `json:"accountPermissions"`
	Persoon            struct {
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
		AdditionalObjects struct{} `json:"additionalObjects"`
		UUID              string   `json:"UUID"`
		StudentNumber     int      `json:"leerlingnummer"`
		FirstName         string   `json:"roepnaam"`
		LastName          string   `json:"achternaam"`
	} `json:"persoon"`
}

const accountUrl = "/rest/v1/account/me?additional=restricties"

func GetCurrentAccount(ctx *auth.Context) (Account, error) {
	var account Account

	body, err := ctx.DoAuthedGet(ctx.Token.SomtodayAPIURL+accountUrl, nil)
	err = json.Unmarshal(body, &account)
	if err != nil {
		return account, fmt.Errorf("error unmarshalling account: %v", err)
	}

	for _, link := range account.Links {
		if link.Rel == "self" {
			account.ID = link.Id
		}
	}
	return account, nil
}
