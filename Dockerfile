# Builder
FROM golang:1.8.3-alpine3.6 as builder

RUN apk --update add make
RUN mkdir -p /go/src/github.com/wobscale/prometheus-irc-exporter
WORKDIR /go/src/github.com/wobscale/prometheus-irc-exporter
COPY . .

RUN make

# Final container
FROM alpine:3.6
MAINTAINER Wobscale (github.com/wobscale)
RUN apk --no-cache add ca-certificates

COPY --from=builder /go/src/github.com/wobscale/prometheus-irc-exporter/bin/prometheus-irc-exporter /prometheus-irc-exporter

ENTRYPOINT ["/prometheus-irc-exporter"]
