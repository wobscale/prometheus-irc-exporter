FROM golang:1.8.3-alpine as builder

RUN apk --update add make
RUN mkdir -p /go/src/github.com/wobscale/prometheus-irc-exporter

WORKDIR /go/src/github.com/wobscale/prometheus-irc-exporter

COPY . .

RUN make


FROM alpine

RUN apk --no-cache add ca-certificates

COPY --from=builder /go/src/github.com/wobscale/prometheus-irc-exporter/bin/prometheus-irc-exporter /prometheus-irc-exporter

ENTRYPOINT ["/prometheus-irc-exporter"]
