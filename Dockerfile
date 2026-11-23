FROM golang:1.27.4@sha256:205c7636d74e0ffbe3d007162df7d74f859f5a8c8322dfb9f907977686c126cf AS mod

WORKDIR /app

RUN --mount=type=bind,source=go.mod,target=go.mod \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=cache,target=/go/pkg/mod \
    go mod download

FROM golang:1.27.4@sha256:205c7636d74e0ffbe3d007162df7d74f859f5a8c8322dfb9f907977686c126cf AS builder

ARG VERSION=dev
ARG COMMIT=none
ARG DATE=unknown
ARG BUILT_BY=docker

WORKDIR /app

ENV CGO_ENABLED=0

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    go build \
    -ldflags="-s -w -extldflags '-static' -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE} -X main.builtBy=${BUILT_BY}" \
    -trimpath \
    -o entrypoint

FROM scratch AS final

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY LICENSE .
COPY --from=builder /app/entrypoint /entrypoint

ENTRYPOINT ["/entrypoint"]
