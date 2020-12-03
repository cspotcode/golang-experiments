.PHONY: all
all: build compress

# Goals: fully static compressed binary suitable for committing to git without anything special like GITLFS

.PHONY: build
build:
	CGO_ENABLED=0 go build -ldflags="-s -w" .

.PHONY: compress
compress:
	../bin/upx ./docker-credential
	
# NOTE docker cli plugins are deprecated, so actually don't do this
.PHONY: install
install:
	mkdir -p ~/.docker/cli-plugins
	cp ./docker-credential ~/.docker/cli-plugins/
