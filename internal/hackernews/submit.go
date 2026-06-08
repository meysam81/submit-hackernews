package hackernews

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Hidden anti-CSRF fields Hacker News requires on the submit form.
const (
	fieldFnid = "fnid"
	fieldFnop = "fnop"
)

// ErrSubmissionRejected is returned when Hacker News accepts the request but
// does not actually create the submission (expired token, duplicate, captcha,
// rate limit, etc.).
var ErrSubmissionRejected = errors.New("submission rejected by Hacker News")

// hnErrorMessages are phrases Hacker News renders (HTTP 200) when it declines a
// submission instead of redirecting to /newest.
var hnErrorMessages = []string{
	"Unknown or expired link",
	"That link has already been submitted",
	"already been submitted",
	"Please confirm that this is your submission",
	"not able to serve your requests",
	"too fast",
	"Validation required",
	"Bad login",
}

// Submit publishes a link to Hacker News. It first fetches the submit form to
// extract the anti-CSRF hidden fields (fnid, fnop), posts the new story, then
// confirms Hacker News actually accepted it (a 2xx status alone does not mean
// the link was published).
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read submit response: %w", err)
	}

	finalURL := resp.Request.URL
	c.log.Debug().
		Int("status", resp.StatusCode).
		Str("final_url", finalURL.String()).
		Int("body_bytes", len(body)).
		Msg("submit response")

	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("submit returned status %d", resp.StatusCode)
	}

	if msg, ok := matchHNError(body); ok {
		return fmt.Errorf("%w: %q", ErrSubmissionRejected, msg)
	}

	// On success HN redirects the POST to the freshly populated /newest (or to
	// the item page for a duplicate). Landing back on the submit form, or
	// anywhere without our link, means it was not published.
	if submissionAccepted(finalURL, body, link) {
		c.log.Info().Str("title", title).Str("url", link).Str("final_url", finalURL.String()).Msg("submitted link to Hacker News")
		return nil
	}

	c.log.Debug().Str("body", snippet(body, 1024)).Msg("unexpected submit response body")
	return fmt.Errorf("%w: Hacker News did not redirect to /newest; re-run with --debug to inspect the response", ErrSubmissionRejected)
}

// submissionAccepted heuristically decides whether HN published the link, based
// on the final URL we landed on and whether the link appears on the page.
func submissionAccepted(finalURL *url.URL, body []byte, link string) bool {
	path := strings.ToLower(finalURL.Path)
	if strings.Contains(path, "newest") || strings.Contains(path, "item") {
		return true
	}
	return strings.Contains(string(body), link)
}

// matchHNError returns the first known Hacker News rejection phrase present in
// the response body.
func matchHNError(body []byte) (string, bool) {
	page := string(body)
	for _, msg := range hnErrorMessages {
		if strings.Contains(page, msg) {
			return msg, true
		}
	}
	return "", false
}

// snippet returns at most n bytes of s with surrounding whitespace collapsed,
// for compact debug logging.
func snippet(b []byte, n int) string {
	s := strings.Join(strings.Fields(string(b)), " ")
	if len(s) > n {
		return s[:n]
	}
	return s
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

	c.log.Debug().Int("status", resp.StatusCode).Str("final_url", resp.Request.URL.String()).Msg("submit form response")

	fields, err := parseHiddenInputs(resp.Body, fieldFnid, fieldFnop)
	if err != nil {
		return nil, fmt.Errorf("parse submit form: %w", err)
	}
	return fields, nil
}
