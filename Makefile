# Makefile for the Nagios Incident.io Notification Plugin

# Name of the binary to generate
BINARY = notify_incident_io

# Version can be set manually or passed in from the GitHub Actions workflow
VERSION ?= development

# Go build flags with version information
BUILD_FLAGS = -ldflags "-s -w -X 'main.version=$(VERSION)'"

.PHONY: all build clean

# Default target: Build the binary
all: build

# Build the binary
build:
	go build $(BUILD_FLAGS) -o $(BINARY) main.go

# Clean up build artifacts
clean:
	rm -f $(BINARY)