// Package hackernews submits links to news.ycombinator.com by driving the same
// login and submit-form flow a browser would: authenticate to obtain a session
// cookie, scrape the anti-CSRF fields from the submit form, then post the story.
package hackernews

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/meysam81/submit-hackernews/internal/logger"
)

const (
	defaultBaseURL   = "https://news.ycombinator.com"
	defaultUserAgent = "Mozilla/5.0 (X11; Linux x86_64; rv:127.0) Gecko/20100101 Firefox/127.0"
	defaultTimeout   = 30 * time.Second
)

// ErrInvalidCredentials is returned when Hacker News rejects the login.
var ErrInvalidCredentials = errors.New("invalid Hacker News credentials")

// ErrFormFieldMissing is returned when a required hidden field is absent from
// the submit form (typically a sign the session is not authenticated).
var ErrFormFieldMissing = errors.New("submit form field not found")

// Client drives the Hacker News web flow. The embedded cookie jar carries the
// login session across to the submit requests.
type Client struct {
	http      *http.Client
	baseURL   string
	userAgent string
	log       *logger.Logger
}

// Option configures a Client.
type Option func(*Client)

// WithHTTPClient overrides the underlying HTTP client (mainly for tests).
func WithHTTPClient(h *http.Client) Option { return func(c *Client) { c.http = h } }

// WithBaseURL overrides the Hacker News base URL (mainly for tests).
func WithBaseURL(u string) Option { return func(c *Client) { c.baseURL = u } }

// WithUserAgent overrides the User-Agent header sent with every request.
func WithUserAgent(ua string) Option { return func(c *Client) { c.userAgent = ua } }

// New builds a Client with a cookie jar installed so that the login session is
// reused by subsequent requests.
func New(log *logger.Logger, opts ...Option) (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("create cookie jar: %w", err)
	}

	c := &Client{
		http:      &http.Client{Timeout: defaultTimeout, Jar: jar},
		baseURL:   defaultBaseURL,
		userAgent: defaultUserAgent,
		log:       log,
	}
	for _, opt := range opts {
		opt(c)
	}

	// A custom HTTP client (e.g. from a test) may not carry a jar; the flow
	// depends on cookies, so install one when missing.
	if c.http.Jar == nil {
		c.http.Jar = jar
	}

	return c, nil
}

// newRequest builds a request against the base URL with the browser-like
// headers Hacker News expects.
func (c *Client) newRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US")
	req.Header.Set("Origin", c.baseURL)
	req.Header.Set("Referer", c.baseURL+"/")

	return req, nil
}

// logClose logs a response-body close error instead of discarding it. It is
// meant to be deferred: defer c.logClose(resp.Body.Close()).
func (c *Client) logClose(err error) {
	if err != nil {
		c.log.Warn().Err(err).Msg("close response body")
	}
}
