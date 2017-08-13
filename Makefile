GO_SOURCES := $(shell find . -type f -name "*.go")

REPO:=github.com/wobscale/prometheus-irc-exporter

./bin/prometheus-irc-exporter: $(GO_SOURCES)
	@mkdir -p ./gopath/src/github.com/wobscale
	@[ -L "./gopath/src/${REPO}" ] || ln -s ../../../.. ./gopath/src/github.com/wobscale/prometheus-irc-exporter
	@mkdir -p ./bin
	GOPATH=${CURDIR}/gopath/ go build -o ./bin/prometheus-irc-exporter github.com/wobscale/prometheus-irc-exporter 

.PHONY: docker-push
docker-push:
	docker build -t wobscale/prometheus-irc-exporter:$(git rev-parse --short HEAD) .
	docker push wobscale/prometheus-irc-exporter:$(git rev-parse --short HEAD)

.PHONY: clean
clean:
	rm -f ./bin/prometheus-irc-exporter
