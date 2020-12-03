.PHONY: all
all: build compress

.PHONY: build
build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o docker-credential .

.PHONY: compress
compress:
	../bin/upx ./docker-credential
	
# NOTE docker cli plugins are deprecated, so actually don't do this
.PHONY: install
install:
	mkdir -p ~/.docker/cli-plugins
	cp ./docker-credential ~/.docker/cli-plugins/
