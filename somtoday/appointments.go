package somtoday

import (
	"encoding/json"
	"fmt"
	"github.com/BenStokmans/go-somtoday/auth"
	"time"
)

type Appointment struct {
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
		Subject              interface{} `json:"vak"`
		TeacherAbbreviations string      `json:"docentAfkortingen"`
		Students             struct {
			Type  string `json:"$type"`
			Items []struct {
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
				UUID          string `json:"UUID"`
				StudentNumber int    `json:"leerlingnummer"`
				FirstName     string `json:"roepnaam"`
				LastName      string `json:"achternaam"`
			} `json:"items"`
		} `json:"leerlingen"`
		OnlineParticipation struct {
			Type  string        `json:"$type"`
			Items []interface{} `json:"items"`
		} `json:"onlineDeelnames"`
	} `json:"additionalObjects"`
	AppointmentType struct {
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
		Name                       string `json:"naam"`
		Description                string `json:"omschrijving"`
		StandardColor              int    `json:"standaardKleur"`
		Category                   string `json:"categorie"`
		Activity                   string `json:"activiteit"`
		PercentageIIVO             int    `json:"percentageIIVO"`
		PresentRegistrationDefault bool   `json:"presentieRegistratieDefault"`
		Active                     bool   `json:"actief"`
		Establishment              `json:"vestiging"`
	} `json:"afspraakType"`
	Location                     string    `json:"locatie"`
	StartTime                    time.Time `json:"beginDatumTijd"`
	EndTime                      time.Time `json:"eindDatumTijd"`
	Title                        string    `json:"titel"`
	Description                  string    `json:"omschrijving"`
	PresentRegistrationRequired  bool      `json:"presentieRegistratieVerplicht"`
	PresentRegistrationProcessed bool      `json:"presentieRegistratieVerwerkt"`
	AppointmentStatus            string    `json:"afspraakStatus"`
	Establishment                `json:"vestiging"`
	Attachments                  []interface{} `json:"bijlagen"`
	StartLessonHour              int           `json:"beginLesuur,omitempty"`
	EndLessonHour                int           `json:"eindLesuur,omitempty"`
}

const appointmentUrl = "/rest/v1/afspraken?begindatum=%s&einddatum=%s&additional=leerlingen&additional=vak&additional=docentAfkortingen&additional=onlineDeelnames&sort=asc-id"

// GetAppointments uses a from and to data formatted as YY-MM-DD
func GetAppointments(from string, to string, ctx *auth.Context) ([]Appointment, error) {
	rangeL, rangeH := 0, 100
	var appointments []Appointment
	for {
		s, err := getAppointmentsRange(from, to, ctx, rangeL, rangeH)
		if err != nil {
			return nil, err
		}
		appointments = append(appointments, s...)
		if len(s) < 100 {
			break
		}
		rangeL, rangeH = rangeH, rangeH+100
	}
	return appointments, nil
}

func getAppointmentsRange(from string, to string, ctx *auth.Context, rangeL, rangeH int) ([]Appointment, error) {
	reqUrl := fmt.Sprintf(ctx.Token.SomtodayAPIURL+appointmentUrl, from, to)
	headers := map[string][]string{
		"Range": {fmt.Sprintf("items=%d-%d", rangeL, rangeH)},
	}
	body, err := ctx.DoAuthedGet(reqUrl, headers)

	var appointments struct {
		Appointments []Appointment `json:"items"`
	}
	err = json.Unmarshal(body, &appointments)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling appointments: %v", err)
	}
	return appointments.Appointments, nil
}
