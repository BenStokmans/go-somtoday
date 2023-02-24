package somtoday

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BenStokmans/go-somtoday/auth"
	"strconv"
)

type Student struct {
	ID    int64  `json:"-"`
	Type  string `json:"$type,omitempty"`
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
		Photo struct {
			Type  string `json:"$type"`
			Links []struct {
				ID  int64  `json:"id"`
				Rel string `json:"rel"`
			} `json:"links"`
			Permissions       []interface{} `json:"permissions"`
			AdditionalObjects struct {
			} `json:"additionalObjects"`
			DataURI string `json:"datauri"`
		} `json:"pasfoto,omitempty"`
		CurrentBatch  Batch `json:"huidigeLichting,omitempty"`
		Establishment `json:"rVestiging,omitempty"`
		ClassTimes    struct {
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
			} `json:"additionalObjects"`
			Establishment `json:"vestiging"`
			Active        bool `json:"actief"`
			ClassTimes    []struct {
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
				Number    int    `json:"nummer"`
				BeginTime string `json:"begintijd"`
				EndTime   string `json:"eindtijd"`
			} `json:"lesuren"`
		} `json:"lestijden,omitempty"`
	} `json:"additionalObjects"`
	UUID          string `json:"UUID"`
	StudentNumber int    `json:"leerlingnummer"`
	FirstName     string `json:"roepnaam"`
	LastName      string `json:"achternaam"`
	Email         string `json:"email,omitempty"`
	Mobile        string `json:"mobielNummer,omitempty"`
	BirthDate     string `json:"geboortedatum,omitempty"`
	Gender        string `json:"geslacht,omitempty"`
	Prefix        string `json:"voorvoegsel,omitempty"`
}

const studentsUrl = "/rest/v1/leerlingen?additional=lestijden&additional=rVestiging&additional=huidigeLichting"
const singleStudentUrl = "/rest/v1/leerlingen/%d?additional=lestijden&additional=rVestiging&additional=huidigeLichting"

func GetCurrentStudent(ctx *auth.Context) (Student, error) {
	acc, err := GetCurrentAccount(ctx)
	if err != nil {
		return Student{}, err
	}
	if acc.Persoon.Type != "leerling.RLeerlingPrimer" {
		return Student{}, errors.New("authenticated user is not a student")
	}
	return GetStudentForAccount(acc, ctx)
}

func GetStudentForAccount(acc Account, ctx *auth.Context) (Student, error) {
	uri := ctx.Token.SomtodayAPIURL + studentsUrl + "&account=" + strconv.FormatInt(acc.ID, 10)
	body, err := ctx.DoAuthedGet(uri, nil)

	var students struct {
		Students []Student `json:"items"`
	}
	var student Student
	err = json.Unmarshal(body, &students)
	if err != nil {
		return student, fmt.Errorf("error unmarshaling student: %v", err)
	}
	for _, s := range students.Students {
		student = s
	}
	if student.UUID == "" {
		return student, fmt.Errorf("error student not found: %v", err)
	}

	for _, link := range student.Links {
		if link.Rel == "self" {
			student.ID = link.Id
		}
	}
	return student, nil
}

func GetStudentByID(id int64, ctx *auth.Context) (Student, error) {
	uri := ctx.Token.SomtodayAPIURL + fmt.Sprintf(singleStudentUrl, id)
	body, err := ctx.DoAuthedGet(uri, nil)

	var student Student
	err = json.Unmarshal(body, &student)
	if err != nil {
		return student, fmt.Errorf("error unmarshaling student: %v", err)
	}
	student.ID = id
	return student, nil
}
