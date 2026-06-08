package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/meysam81/submit-hackernews/internal/hackernews"
	"github.com/meysam81/submit-hackernews/internal/logger"
	"github.com/urfave/cli/v3"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

func main() {
	var (
		title      string
		link       string
		username   string
		password   string
		verboseRaw string
	)

	cmd := &cli.Command{
		Name:    "submit-hackernews",
		Usage:   "Submit a link to Hacker News",
		Version: fmt.Sprintf("%s (commit %s, built %s by %s)", version, commit, date, builtBy),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "title",
				Aliases:     []string{"t"},
				Usage:       "title of the submission",
				Sources:     cli.EnvVars("HACKERNEWS_TITLE"),
				Destination: &title,
			},
			&cli.StringFlag{
				Name:        "url",
				Aliases:     []string{"u"},
				Usage:       "URL of the submission",
				Sources:     cli.EnvVars("HACKERNEWS_URL"),
				Destination: &link,
			},
			&cli.StringFlag{
				Name:        "username",
				Aliases:     []string{"U"},
				Usage:       "Hacker News username",
				Sources:     cli.EnvVars("HACKERNEWS_USERNAME"),
				Destination: &username,
			},
			&cli.StringFlag{
				Name:        "password",
				Aliases:     []string{"p"},
				Usage:       "Hacker News password",
				Sources:     cli.EnvVars("HACKERNEWS_PASSWORD"),
				Destination: &password,
			},
			&cli.StringFlag{
				Name:        "verbose",
				Usage:       "enable verbose (debug) logging; any non-empty value enables it",
				Sources:     cli.EnvVars("VERBOSE"),
				Destination: &verboseRaw,
			},
		},
		Action: func(ctx context.Context, _ *cli.Command) error {
			log := logger.New(truthy(verboseRaw))

			var errs []error
			if title == "" {
				errs = append(errs, errors.New("--title is required"))
			}
			if link == "" {
				errs = append(errs, errors.New("--url is required"))
			}
			if username == "" {
				errs = append(errs, errors.New("--username is required"))
			}
			if password == "" {
				errs = append(errs, errors.New("--password is required"))
			}
			if err := errors.Join(errs...); err != nil {
				return err
			}

			client, err := hackernews.New(&log)
			if err != nil {
				return fmt.Errorf("init client: %w", err)
			}
			if err := client.Login(ctx, username, password); err != nil {
				return fmt.Errorf("login: %w", err)
			}
			if err := client.Submit(ctx, title, link); err != nil {
				return fmt.Errorf("submit: %w", err)
			}
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

// truthy reports whether a string sourced from an environment variable should
// be treated as enabled. It mirrors the original shell behaviour (any non-empty
// value is true) while still treating the obvious negatives as off.
func truthy(s string) bool {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "", "0", "false", "no", "off":
		return false
	default:
		return true
	}
}
