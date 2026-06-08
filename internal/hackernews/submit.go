package hackernews

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Hidden anti-CSRF fields Hacker News requires on the submit form.
const (
	fieldFnid = "fnid"
	fieldFnop = "fnop"
)

// Submit publishes a link to Hacker News. It first fetches the submit form to
// extract the anti-CSRF hidden fields (fnid, fnop), then posts the new story.
func (c *Client) Submit(ctx context.Context, title, link string) error {
	fields, err := c.submitFormFields(ctx)
	if err != nil {
		return err
	}

	form := url.Values{
		fieldFnid: {fields[fieldFnid]},
		fieldFnop: {fields[fieldFnop]},
		"title":   {title},
		"url":     {link},
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/r", strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("build submit request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	c.log.Debug().Str("title", title).Str("url", link).Msg("submitting link")

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("submit request: %w", err)
	}
	defer func() { c.logClose(resp.Body.Close()) }()

	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("submit returned status %d", resp.StatusCode)
	}

	c.log.Info().Str("title", title).Str("url", link).Msg("submitted link to Hacker News")
	return nil
}

// submitFormFields fetches the submit form and returns its fnid/fnop hidden
// fields, which Hacker News requires on the subsequent POST.
func (c *Client) submitFormFields(ctx context.Context) (map[string]string, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/submitlink", nil)
	if err != nil {
		return nil, fmt.Errorf("build submit-form request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch submit form: %w", err)
	}
	defer func() { c.logClose(resp.Body.Close()) }()

	fields, err := parseHiddenInputs(resp.Body, fieldFnid, fieldFnop)
	if err != nil {
		return nil, fmt.Errorf("parse submit form: %w", err)
	}
	return fields, nil
}
