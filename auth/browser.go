package auth

import (
	"context"
	"errors"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/sirupsen/logrus"
	"net/url"
	"strings"
)

type ControlledBrowser struct {
	logger *logrus.Logger
	ctx    context.Context
	cancel context.CancelFunc

	code chan string
}

func NewControlledBrowser(logger *logrus.Logger) *ControlledBrowser {
	opts := []chromedp.ExecAllocatorOption{
		chromedp.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0"),
		chromedp.Flag("headless", false),
	}
	ctx, cancel := chromedp.NewExecAllocator(context.Background(),
		append(chromedp.DefaultExecAllocatorOptions[:], opts...)...)

	ctx, cancel = chromedp.NewContext(ctx,
		chromedp.WithDebugf(logger.Debugf),
		chromedp.WithErrorf(logger.Errorf),
	)

	return &ControlledBrowser{
		logger: logger,
		ctx:    ctx,
		cancel: cancel,
		code:   make(chan string),
	}
}

func (c *ControlledBrowser) getAuthCode(authUrl string, ctx *Context) (string, error) {
	if ctx.oidCurl == "" {
		return "", errors.New("no oidCurl set please call Context::SetOIDCurl")
	}
	defer c.cancel()

	eventCtx, cancelEvent := context.WithCancel(c.ctx)
	defer cancelEvent()

	callBackUrl := strings.Replace(ctx.redirectUrl, "oauth", "oauth:443", 1)
	chromedp.ListenTarget(eventCtx, func(ev interface{}) {
		switch ev := ev.(type) {
		case *network.EventRequestWillBeSent:
			if !strings.HasPrefix(ev.Request.URL, callBackUrl) {
				return
			}
			codeUrl, _ := url.Parse(ev.Request.URL)
			code := codeUrl.Query().Get("code")
			if code == "" {
				return
			}
			c.code <- code
		}
	})

	if err := chromedp.Run(c.ctx,
		chromedp.Navigate(authUrl),
	); err != nil {
		return "", err
	}
	return <-c.code, nil
}

func (c *ControlledBrowser) GetToken(ctx *Context) (Token, error) {
	authUrl, verifier, err := GetOAuthUrl(ctx)
	if err != nil {
		return Token{}, err
	}

	code, err := c.getAuthCode(authUrl, ctx)
	if err != nil {
		return Token{}, err
	}

	return GetToken(code, verifier, ctx)
}
