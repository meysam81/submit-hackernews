# Submit HackerNews

[![GitHub repo size](https://img.shields.io/github/repo-size/meysam81/submit-hackernews)](https://github.com/meysam81/submit-hackernews)
[![GitHub commit activity](https://img.shields.io/github/commit-activity/m/meysam81/submit-hackernews)](https://github.com/meysam81/submit-hackernews/commits/main/)
[![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/meysam81/submit-hackernews)](https://github.com/meysam81/submit-hackernews)
[![pre-commit.ci status](https://results.pre-commit.ci/badge/github/meysam81/submit-hackernews/main.svg)](https://results.pre-commit.ci/latest/github/meysam81/submit-hackernews/main)
[![Docker Image Version](https://ghcr-badge.egpl.dev/meysam81/submit-hackernews/latest_tag?color=%2344cc11&ignore=latest&label=Docker+Image+Version&trim=)](https://github.com/users/meysam81/packages/container/package/submit-hackernews)
[![Docker Image Size](https://ghcr-badge.egpl.dev/meysam81/submit-hackernews/size?color=%2344cc11&tag=latest&label=Docker+Image+Size&trim=)](https://github.com/users/meysam81/packages/container/package/submit-hackernews)

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Submit HackerNews](#submit-hackernews)
  - [Introduction](#introduction)
  - [Usage: GitHub Actions](#usage-github-actions)
  - [Usage: Docker](#usage-docker)
  - [Usage: CLI](#usage-cli)
  - [Configuration](#configuration)
  - [Development](#development)
  - [Star History](#star-history)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Introduction

A tiny Go CLI that submits links to [Hacker News](https://news.ycombinator.com).
It logs in with your credentials, reads the anti-CSRF fields off the submit
form, and posts the story — handy for automating submissions on a schedule.

It ships three ways: as a GitHub Action, as a container image, and as a static
binary.

## Usage: GitHub Actions

```yaml
name: ci

on:
  workflow_dispatch:
    inputs:
      title:
        description: The title of the link to submit.
        required: true
      url:
        description: The URL of the link to submit.
        required: true
      verbose:
        type: boolean
        description: Verbose?
        default: true

jobs:
  submit-hackernews:
    if: github.event_name == 'workflow_dispatch'
    runs-on: ubuntu-latest
    steps:
      - name: Submit link to Hacker News
        uses: meysam81/submit-hackernews@v1
        with:
          username: ${{ secrets.HACKERNEWS_USERNAME }}
          password: ${{ secrets.HACKERNEWS_PASSWORD }}
          title: ${{ github.event.inputs.title }}
          url: ${{ github.event.inputs.url }}
          verbose: ${{ github.event.inputs.verbose }}
```

## Usage: Docker

```bash
docker run \
  --name submit-hackernews \
  --rm \
  -e "HACKERNEWS_USERNAME=your_username" \
  -e "HACKERNEWS_PASSWORD=your_password" \
  ghcr.io/meysam81/submit-hackernews:v1 \
  -t "This is the title of submission" \
  -u "https://example.com"
```

## Usage: CLI

Download a binary from the [releases page](https://github.com/meysam81/submit-hackernews/releases),
or build from source (see [Development](#development)), then:

```bash
submit-hackernews \
  --username your_username \
  --password your_password \
  --title "This is the title of submission" \
  --url "https://example.com"
```

## Configuration

Every flag has an environment-variable fallback, so the tool works equally well
from a shell, a container, or a GitHub Action.

| Flag                | Alias | Environment variable   | Required | Description                                       |
| ------------------- | ----- | ---------------------- | -------- | ------------------------------------------------- |
| `--title`           | `-t`  | `HACKERNEWS_TITLE`     | yes      | Title of the submission.                          |
| `--url`             | `-u`  | `HACKERNEWS_URL`       | yes      | URL of the submission.                            |
| `--username`        | `-U`  | `HACKERNEWS_USERNAME`  | yes      | Hacker News username.                             |
| `--password`        | `-p`  | `HACKERNEWS_PASSWORD`  | yes      | Hacker News password.                             |
| `--verbose`         |       | `VERBOSE`              | no       | Any non-empty value enables debug logging.        |

## Development

Requires Go (see `go.mod` for the version).

```bash
go build ./...                 # build
go test -race -count=1 ./...   # test
go vet ./...                   # vet
gofmt -l .                     # format check
```

## Star History

<a href="https://star-history.com/#meysam81/submit-hackernews&Timeline">
 <picture>
   <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=meysam81/submit-hackernews&type=Timeline&theme=dark" />
   <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=meysam81/submit-hackernews&type=Timeline" />
   <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=meysam81/submit-hackernews&type=Timeline" />
 </picture>
</a>
