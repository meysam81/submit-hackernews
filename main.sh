#!/usr/bin/env sh

set -eu

USERNAME=
PASSWORD=
TITLE=
URL=
COOKIE_JAR=

usage() {
    echo "Usage: $0 {-t|--title TITLE} {-u|--url URL} [-U|--username USERNAME] [-p|--password PASSWORD]"
    echo "Options:"
    echo "  -t, --title     Specify the title of the submission"
    echo "  -u, --url       Specify the URL of the submission"
    echo "  -U, --username  Specify the HackerNews username (defaults to \$HACKERNEWS_USERNAME)"
    echo "  -p, --password  Specify the HackerNews password (defaults to \$HACKERNEWS_PASSWORD)"
    echo ""
    echo "All four arguments MUST be provided."
}

prepare_cookie() {
    COOKIE_JAR=$(mktemp)
}

parse_args() {
    while [ $# -gt 0 ]; do
        key="$1"
        case $key in
            -t|--title)
                TITLE="$2"
                shift
                shift
            ;;
            -u|--url)
                URL="$2"
                shift
                shift
            ;;
            -U|--username)
                USERNAME="$2"
                shift
                shift
            ;;
            -p|--password)
                PASSWORD="$2"
                shift
                shift
            ;;
            -h|--help)
                usage
                exit 0
            ;;
            *)
                usage
                exit 1
            ;;
        esac
    done
}

parse_envs() {
    set +u

    if [ -z "$USERNAME" ]; then
        USERNAME=$HACKERNEWS_USERNAME
    fi

    if [ -z "$PASSWORD" ]; then
        PASSWORD=$HACKERNEWS_PASSWORD
    fi

    set -u
}


validate_args() {
    if [ -z "$TITLE" ] || [ -z "$URL" ] || [ -z "$USERNAME" ] || [ -z "$PASSWORD" ]; then
        usage
        exit 1
    fi
}

login() {
    curl -b $COOKIE_JAR -c $COOKIE_JAR 'https://news.ycombinator.com/login' \
    -X POST \
    -H 'User-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:127.0) Gecko/20100101 Firefox/127.0' \
    -H 'Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8' \
    -H 'Accept-Language: en-US' \
    -H 'Accept-Encoding: gzip, deflate, br, zstd' \
    -H 'Referer: https://news.ycombinator.com/' \
    -H 'Content-Type: application/x-www-form-urlencoded' \
    -H 'Origin: https://news.ycombinator.com' \
    -H 'DNT: 1' \
    -H 'Sec-GPC: 1' \
    -H 'Connection: keep-alive' \
    -H 'Upgrade-Insecure-Requests: 1' \
    -H 'Sec-Fetch-Dest: document' \
    -H 'Sec-Fetch-Mode: navigate' \
    -H 'Sec-Fetch-Site: same-origin' \
    -H 'Sec-Fetch-User: ?1' \
    -H 'Priority: u=1' \
    -H 'Pragma: no-cache' \
    -H 'Cache-Control: no-cache' \
    -H 'TE: trailers' \
    --data-raw "goto=newest&acct=${USERNAME}&pw=${PASSWORD}" \
    -v
}
submit() {
    form_info=$(curl -c $COOKIE_JAR -b $COOKIE_JAR 'https://news.ycombinator.com/submitlink' -sS | pup 'input')


    # <input type="hidden" name="fnid" value="RANDOM_STRING">
    # <input type="hidden" name="fnop" value="submit-page">
    # <input type="text" name="title" size="50" autofocus="t" oninput="tlen(this)" onfocus="tlen(this)">
    # <input type="url" name="url" size="50">
    # <input type="submit" value="submit">

    fnid=$(echo $form_info | grep -oP 'name="fnid" value="\K[^"]+')
    fnop=$(echo $form_info | grep -oP 'name="fnop" value="\K[^"]+')
    submit="submit"

    echo $TITLE
    echo $URL

    curl -L -c $COOKIE_JAR -b $COOKIE_JAR 'https://news.ycombinator.com/r' -X POST  \
    -H 'User-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:127.0) Gecko/20100101 Firefox/127.0' \
    -H 'Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8' \
    -H 'Accept-Language: en-US' \
    -H 'Accept-Encoding: gzip, deflate, br, zstd' \
    -H 'Referer: https://news.ycombinator.com/' \
    -H 'Content-Type: application/x-www-form-urlencoded' \
    -H 'Origin: https://news.ycombinator.com' \
    -H 'DNT: 1' \
    -H 'Sec-GPC: 1' \
    -H 'Connection: keep-alive' \
    -H 'Upgrade-Insecure-Requests: 1' \
    -H 'Sec-Fetch-Dest: document' \
    -H 'Sec-Fetch-Mode: navigate' \
    -H 'Sec-Fetch-Site: same-origin' \
    -H 'Sec-Fetch-User: ?1' \
    -H 'Priority: u=1' \
    -H 'Pragma: no-cache' \
    -H 'Cache-Control: no-cache' \
    -H 'TE: trailers' \
    --data-urlencode "fnid=$fnid"  \
    --data-urlencode "fnop=$fnop"  \
    --data-urlencode "title=$TITLE"  \
    --data-urlencode "url=$URL"  \
    -v

}

main() {
    parse_args "$@"
    parse_envs
    validate_args

    prepare_cookie
    login
    submit
}

main "$@"
