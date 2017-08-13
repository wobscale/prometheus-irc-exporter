GO_SOURCES := $(shell find . -type f -name "*.go")

./bin/prometheus-irc-exporter: $(GO_SOURCES)
	@mkdir -p ./bin
	go build -o ./bin/prometheus-irc-exporter .

