FROM golang:1.22-alpine AS deps

RUN apk add --no-cache git \
  && go install github.com/ericchiang/pup@v0.4.0


FROM curlimages/curl:8.8.0 AS runner

USER root
RUN apk update && apk upgrade
USER curl_user:curl_group

WORKDIR /app

COPY --from=deps /go/bin/pup /usr/local/bin/pup

COPY . .

ENTRYPOINT ["/app/main.sh"]