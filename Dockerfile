FROM golang:1.11-alpine AS builder

RUN apk update && \
  apk add --no-cache curl git && \
  curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.5.4/dep-linux-amd64 && \
  chmod +x /usr/local/bin/dep

WORKDIR /go/src/github.com/thewizardplusplus/go-link-shortener
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure -vendor-only

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install -a -ldflags='-w -s -extldflags "-static"' ./...

FROM scratch
COPY --from=builder /go/bin/go-link-shortener /usr/local/bin/go-link-shortener
CMD ["/usr/local/bin/go-link-shortener"]
