package auth

import (
	"encoding/base64"
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

type Token struct {
	AccessToken    string    `json:"access_token"`
	RefreshToken   string    `json:"refresh_token"`
	SomtodayAPIURL string    `json:"somtoday_api_url"`
	SomtodayOopURL string    `json:"somtoday_oop_url"`
	Scope          string    `json:"scope"`
	SomtodayTenant string    `json:"somtoday_tenant"`
	IDToken        string    `json:"id_token"`
	TokenType      string    `json:"token_type"`
	ExpiresIn      int       `json:"expires_in"`
	ExpireTime     time.Time `json:"-"`
}

func (t Token) WithExpireTime() Token {
	jwtParts := strings.Split(t.AccessToken, ".")
	if len(jwtParts) < 3 {
		return t
	}
	payloadData, err := base64.StdEncoding.DecodeString(jwtParts[1])
	if err != nil {
		return t
	}

	var payload map[string]string
	err = json.Unmarshal(payloadData, &payload)
	if err != nil {
		return t
	}

	expStr, ok := payload["exp"]
	if !ok {
		return t
	}
	exp, err := strconv.ParseInt(expStr, 10, 64)
	if err != nil {
		return t
	}
	t.ExpireTime = time.Unix(exp, 0)
	return t
}
