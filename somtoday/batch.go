package somtoday

type Batch struct {
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
	Name            string `json:"naam"`
	BatchSchoolYear []struct {
		Type  string `json:"$type"`
		Links []struct {
			Id   int    `json:"id"`
			Rel  string `json:"rel"`
			Type string `json:"type"`
		} `json:"links"`
		Permissions       []interface{} `json:"permissions"`
		AdditionalObjects struct {
		} `json:"additionalObjects"`
		SchoolYear     `json:"schooljaar"`
		YearNumber     int  `json:"leerjaar"`
		HasExamDossier bool `json:"heeftExamendossier"`
	} `json:"lichtingSchooljaren"`
	EducationType struct {
		Type  string `json:"$type"`
		Links []struct {
			Id   int    `json:"id"`
			Rel  string `json:"rel"`
			Type string `json:"type"`
		} `json:"links"`
		Permissions       []interface{} `json:"permissions"`
		AdditionalObjects struct {
		} `json:"additionalObjects"`
		Afkorting   string `json:"afkorting"`
		IsOnderbouw bool   `json:"isOnderbouw"`
	} `json:"onderwijssoort"`
}
