SHELL = bash
GOTOOLS = \
	github.com/mjibson/esc \
	github.com/mitchellh/gox \

GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)

# Get the git commit
GIT_LATEST_TAG=$(shell git describe --abbrev=0 --tags)
VERSION_IMPORT=github.com/andrexus/cloud-initer/cmd
GOLDFLAGS=-X $(VERSION_IMPORT).Version=$(GIT_LATEST_TAG)

export GOLDFLAGS

# all builds binaries for all targets
all: bin

bin: tools
	go generate
	@echo "==> Building..."
	gox -ldflags "${GOLDFLAGS}" -osarch "darwin/amd64 linux/386 linux/amd64 linux/arm" -output "build/{{.OS}}_{{.Arch}}_{{.Dir}}"

tools:
	go get -u -v $(GOTOOLS)

.PHONY: all bin tools