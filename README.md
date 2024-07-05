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

## Star History

<a href="https://star-history.com/#meysam81/submit-hackernews&Timeline">
 <picture>
   <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=meysam81/submit-hackernews&type=Timeline&theme=dark" />
   <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=meysam81/submit-hackernews&type=Timeline" />
   <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=meysam81/submit-hackernews&type=Timeline" />
 </picture>
</a>
