package hackernews

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Login authenticates against Hacker News, populating the cookie jar with the
// session cookie used by subsequent requests.
func (c *Client) Login(ctx context.Context, username, password string) error {
	form := url.Values{
		"goto": {"newest"},
		"acct": {username},
		"pw":   {password},
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/login", strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("build login request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	c.log.Debug().Str("username", username).Msg("logging in to Hacker News")

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("login request: %w", err)
	}
	defer func() { c.logClose(resp.Body.Close()) }()

	// Hacker News answers a failed login with a 200 page containing "Bad login"
	// rather than a non-2xx status, so inspect the body to detect failure.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read login response: %w", err)
	}
	if strings.Contains(string(body), "Bad login") {
		return ErrInvalidCredentials
	}

	c.log.Info().Str("username", username).Msg("logged in")
	return nil
}
