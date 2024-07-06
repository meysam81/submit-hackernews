# Submit HackerNews

[![GitHub repo size](https://img.shields.io/github/repo-size/meysam81/submit-hackernews)](https://github.com/meysam81/submit-hackernews)
[![GitHub commit activity](https://img.shields.io/github/commit-activity/m/meysam81/submit-hackernews)](https://github.com/meysam81/submit-hackernews/commits/main/)
[![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/meysam81/submit-hackernews)](https://github.com/meysam81/submit-hackernews)
[![pre-commit.ci status](https://results.pre-commit.ci/badge/github/meysam81/submit-hackernews/main.svg)](https://results.pre-commit.ci/latest/github/meysam81/submit-hackernews/main)

This small Shell script will submit links to HackerNews. It is hugely
beneficial to automate the task of submission at desired schedules.

To run the script, run the following command:

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

## Usage: GitHub Actions

You can run this in a GitHub Actions workflow. Here is an example:

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

## Star History

<a href="https://star-history.com/#meysam81/submit-hackernews&Timeline">
 <picture>
   <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=meysam81/submit-hackernews&type=Timeline&theme=dark" />
   <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=meysam81/submit-hackernews&type=Timeline" />
   <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=meysam81/submit-hackernews&type=Timeline" />
 </picture>
</a>
