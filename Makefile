# Makefile for the Nagios Incident.io Notification Plugin

# Name of the binary to generate (base name)
BINARY = notify_incident_io

# Version can be set manually or passed in from the GitHub Actions workflow
VERSION ?= development

# Go build flags with version information
BUILD_FLAGS = -ldflags "-s -w -X 'main.version=$(VERSION)'"

# Platforms to build for
PLATFORMS = linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

.PHONY: all build clean release

# Default target: Build the binary for the current OS and architecture
all: build

# Build the binary for the current OS and architecture
build:
	go build $(BUILD_FLAGS) -o $(BINARY) main.go

# Cross-compile binaries for all target platforms
release: clean
	@for platform in $(PLATFORMS); do \
		CGO_ENABLED=0 GOOS=$${platform%/*} GOARCH=$${platform#*/} go build $(BUILD_FLAGS) -o $(BINARY)-$${platform%/*}-$${platform#*/} main.go || exit 1; \
	done

# Clean up build artifacts
clean:
	rm -f $(BINARY) $(BINARY)-*