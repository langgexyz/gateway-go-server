pwd := $(shell pwd)

.PHONY: build-linux debug

build-linux:
	@export GOPROXY=https://goproxy.io,direct && go mod tidy && GOOS=linux GOARCH=amd64 go build -o ./bin/gateway-go-server-linux ./src/main/main.go

debug:
	@export GOPROXY=https://goproxy.io,direct && go mod tidy && go build -o ./bin/gateway-go-server ./src/main/main.go; \
	cd bin; \
	if [ -f ./pid ]; then \
		kill -TERM `cat ./pid`; \
	fi; \
	nohup ./gateway-go-server -c ./config.debug.json > ./error.log 2>&1 & echo $$! > pid; \
	cd ..; \
	echo "Gateway debug mode started";
