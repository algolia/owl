# OS-specific variables
UNAME=$(shell uname)
ifeq ($(UNAME),Darwin)
  SORT_BIN=gsort
else
  SORT_BIN=sort
endif

# Build variables
BUILD_DIR=_build
LAST_TAG=$(shell git tag | $(SORT_BIN) -V | tail -n 1)
HASH_HEAD=$(shell git show-ref --head | head -1 | awk '{print $$1}')
HASH_LAST_TAG=$(shell git rev-list -n 1 $(LAST_TAG))
GIT_TAG=$(shell if [ $(HASH_HEAD) = $(HASH_LAST_TAG) ]; then echo $(LAST_TAG); else echo "head"; fi)
BIN=owl-${GIT_TAG}
LDFLAGS="-X github.com/algolia/owl/info.GitTag=${GIT_TAG}"

NEED_FORMATTING=$(shell for f in $(shell find . -name '*.go' | grep -v vendor); do gofmt -l $$f; done)

all: install

# Build and install

install:
	go install

build: clean-build linux-amd64 mac-amd64

linux-amd64:
	GOOS=linux GOARCH=amd64 go build -ldflags ${LDFLAGS} -o $(BUILD_DIR)/$@/$(BIN)

mac-amd64:
	GOOS=darwin GOARCH=amd64 go build -ldflags ${LDFLAGS} -o $(BUILD_DIR)/$@/$(BIN)

deps:
	glide install

check-formatting:
ifneq ($(NEED_FORMATTING),)
	@echo $(NEED_FORMATTING)
	@exit 1
endif

test: test-unit

test-unit:
	@echo '> Run unit tests'
	go test -v $(shell glide novendor)


clean: clean-build

clean-build:
	go clean
	rm -rf $(BUILD_DIR)

.PHONY: install build test clean
