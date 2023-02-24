package somtoday

type Subject struct {
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
}
