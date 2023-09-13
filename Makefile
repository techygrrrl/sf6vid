#!/usr/bin/make -f
.PHONY: all test run lint

TARGET="sf6vid"
BUILD_TAGS=-tags ""

# This is the default. Running `make` will run this.
all:
	go build $(BUILD_TAGS) -o $(TARGET)

test:
	go test $(BUILD_TAGS) -v ./... -cover

lint:
	goimports -w .
	gci write .

run: all
	./$(TARGET)

prerelease:
	./scripts/prerelease_version.sh
