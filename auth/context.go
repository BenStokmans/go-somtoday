package auth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BenStokmans/go-somtoday/organisation"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Context struct {
	clientId     string
	redirectUrl  string
	oidCurl      string
	UserAgent    string
	logger       *logrus.Logger
	Organisation organisation.Organization
	Token        Token

	cacheDir string
}

func NewContext(logger *logrus.Logger) *Context {
	return &Context{
		clientId:    "D50E0C06-32D1-4B41-A137-A9A850C892C2",
		redirectUrl: "somtodayleerling://oauth/callback",
		logger:      logger,
		UserAgent:   "Leerling/134 CFNetwork/1404.0.5",
	}
}

func (ctx *Context) DoAuthedGet(uri string, headers map[string][]string) ([]byte, error) {
	reqUrl, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	req := &http.Request{
		Method: "GET",
		URL:    reqUrl,
		Header: map[string][]string{
			"User-Agent":    {ctx.UserAgent},
			"Authorization": {"Bearer " + ctx.Token.AccessToken},
			"Accept":        {"application/json"},
		},
	}

	// set and override potential headers
	for s, i := range headers {
		req.Header[s] = i
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}
	if res.StatusCode != 200 && res.StatusCode != 206 {
		return nil, fmt.Errorf("error unexpected status code %d: %s", res.StatusCode, body)
	}
	return body, nil
}

func (ctx *Context) TokenValid() error {
	if time.Now().After(ctx.Token.ExpireTime) {
		token, err := RefreshToken(ctx)
		if err != nil {
			return err
		}
		err = ctx.SetToken(token)
		if err != nil {
			return err
		}
	}
	reqUrl, _ := url.Parse(ctx.Token.SomtodayAPIURL + "/rest/v1/account/me")
	req := &http.Request{
		Method: "GET",
		URL:    reqUrl,
		Header: map[string][]string{
			"User-Agent":    {ctx.UserAgent},
			"Authorization": {"Bearer " + ctx.Token.AccessToken},
			"Accept":        {"application/json"},
		},
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error retrieving account: %v", err)
	}
	if res.StatusCode != 200 {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("error unexpected status code %d: %s", res.StatusCode, body)
	}
	return nil
}

func (ctx *Context) LoadOrLogin(username, password string) error {
	if ctx.cacheDir != "" {
		if _, err := os.Stat(filepath.Join(ctx.cacheDir, "token.json")); !os.IsNotExist(err) {
			var token Token
			token, err = ctx.retrieveTokenFromCache()
			if err != nil {
				return err
			}
			return ctx.SetToken(token)
		}
	}
	return ctx.doAuthFlow()
}

func (ctx *Context) LoadOrDoSSOAuth() error {
	if ctx.cacheDir != "" {
		if _, err := os.Stat(filepath.Join(ctx.cacheDir, "token.json")); !os.IsNotExist(err) {
			var token Token
			token, err = ctx.retrieveTokenFromCache()
			if err != nil {
				return err
			}
			return ctx.SetToken(token)
		}
	}
	return ctx.doAuthFlow()
}

func (ctx *Context) DoLogin(username, password string) error {
	token, err := GetTokenPassword(username, password, ctx)
	if err != nil {
		return err
	}
	return ctx.SetToken(token)
}

func (ctx *Context) DoSSOAuth() error {
	return ctx.doAuthFlow()
}

func (ctx *Context) retrieveTokenFromCache() (Token, error) {
	data, err := os.ReadFile(filepath.Join(ctx.cacheDir, "token.json"))
	if err != nil {
		return Token{}, err
	}

	err = json.Unmarshal(data, &ctx.Token)
	if err != nil {
		return Token{}, fmt.Errorf("error unmarshalling token from cache: %v", err)
	}

	// we include this logic here also to properly handle any errors
	jwtParts := strings.Split(ctx.Token.AccessToken, ".")
	if len(jwtParts) < 3 {
		return Token{}, errors.New("invalid token in cache file")
	}
	payloadData, err := base64.RawURLEncoding.DecodeString(jwtParts[1])
	if err != nil {
		return Token{}, err
	}

	var payload map[string]interface{}
	err = json.Unmarshal(payloadData, &payload)
	if err != nil {
		return Token{}, err
	}

	exp, ok := payload["exp"].(float64)
	if !ok {
		return Token{}, errors.New("could not find expiry time in token")
	}

	ctx.Token.ExpireTime = time.Unix(int64(exp), 0)
	if time.Now().After(ctx.Token.ExpireTime) {
		return RefreshToken(ctx)
	}
	return ctx.Token, nil
}

func (ctx *Context) SetToken(token Token) error {
	ctx.Token = token
	if ctx.cacheDir == "" {
		return nil
	}
	data, err := json.Marshal(ctx.Token)
	if err != nil {
		return err
	}

	// overwrite or create token file
	return os.WriteFile(filepath.Join(ctx.cacheDir, "token.json"), data, 0755)
}

func (ctx *Context) doAuthFlow() error {
	browser := NewControlledBrowser(ctx.logger)
	token, err := browser.GetToken(ctx)
	if err != nil {
		return err
	}
	return ctx.SetToken(token)
}

func (ctx *Context) WithDefaultCache() *Context {
	ctx.cacheDir, _ = os.UserCacheDir()
	if _, err := os.Stat(ctx.cacheDir); !os.IsNotExist(err) {
		_ = os.Mkdir(ctx.cacheDir, 0755)
	}
	return ctx
}

func (ctx *Context) WithCache(path string) *Context {
	ctx.cacheDir = path
	if _, err := os.Stat(ctx.cacheDir); !os.IsNotExist(err) {
		_ = os.Mkdir(ctx.cacheDir, 0755)
	}
	return ctx
}

func (ctx *Context) SetOrganisationByID(tenantId string) error {
	org, err := organisation.GetOrganizationByID(tenantId)
	if err != nil {
		return err
	}
	ctx.Organisation = org
	return nil
}

func (ctx *Context) SetOrganisationByName(name string) error {
	org, err := organisation.GetOrganizationByName(name)
	if err != nil {
		return err
	}
	ctx.Organisation = org
	return nil
}

func (ctx *Context) SetOIDCurl(oidCurl string) error {
	if ctx.Organisation.UUID == "" {
		return errors.New("no organisation for current context please call Context::SetOrganisationBy(Id/Name)")
	}
	for _, s := range ctx.Organisation.OID {
		if s.URL == oidCurl {
			ctx.oidCurl = oidCurl
			return nil
		}
	}
	return errors.New("invalid oidCurl for current organisation")
}
