package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const oAuthUrl = "https://somtoday.nl/oauth2/authorize?redirect_uri=%s&client_id=%s&response_type=code&prompt=login&scope=openid&code_challenge=%s&code_challenge_method=S256&tenant_uuid=%s&oidc_iss=%s"
const authorizationUrl = "https://inloggen.somtoday.nl/oauth2/token?grant_type=authorization_code&redirect_uri=%s&code_verifier=%s&code=%s&scope=openid&client_id=%s"
const refreshUrl = "https://inloggen.somtoday.nl/oauth2/token?grant_type=refresh_token&refresh_token=%s&client_id=%s"

func RefreshToken(ctx *Context) (Token, error) {
	reqUrl, _ := url.Parse(fmt.Sprintf(
		refreshUrl,
		ctx.Token.RefreshToken,
		ctx.clientId,
	),
	)
	req := &http.Request{
		Method: "POST",
		URL:    reqUrl,
		Header: map[string][]string{
			"User-Agent":   {ctx.UserAgent},
			"Content-Type": {"application/x-www-form-urlencoded"},
		},
	}
	return tokenFromRequest(req)
}

func GetToken(code string, verifier string, ctx *Context) (Token, error) {
	reqUrl, _ := url.Parse(fmt.Sprintf(
		authorizationUrl,
		url.QueryEscape(ctx.redirectUrl),
		verifier,
		code,
		ctx.clientId,
	),
	)
	req := &http.Request{
		Method: "POST",
		URL:    reqUrl,
		Header: map[string][]string{
			"User-Agent":   {ctx.UserAgent},
			"Content-Type": {"application/x-www-form-urlencoded"},
		},
	}
	return tokenFromRequest(req)
}

func GetOAuthUrl(ctx *Context) (string, string, error) {
	if ctx.Organisation.UUID == "" {
		return "", "", errors.New("no for current organisation please call Context::SetOrganisationBy(Id/Name)")
	}
	if ctx.oidCurl == "" {
		return "", "", errors.New("no oidCurl set please call Context::SetOIDCurl")
	}
	verifier, err := getVerifier()
	if err != nil {
		return "", "", err
	}
	challenge := getChallenge(verifier)
	somUrl := fmt.Sprintf(oAuthUrl,
		url.QueryEscape(ctx.redirectUrl),
		ctx.clientId,
		challenge,
		ctx.Organisation.UUID,
		url.QueryEscape(ctx.oidCurl),
	)

	return somUrl, verifier, nil
}
func tokenFromRequest(req *http.Request) (Token, error) {
	var token Token

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return token, fmt.Errorf("error retrieving token: %v", err)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return token, fmt.Errorf("error reading token: %v", err)
	}
	if res.StatusCode != 200 {
		return token, fmt.Errorf("error unexpected status code %d: %s", res.StatusCode, string(body))
	}
	err = json.Unmarshal(body, &token)
	if err != nil {
		return token, fmt.Errorf("error unmarshalling token: %v", err)
	}
	return token.WithExpireTime(), nil
}
