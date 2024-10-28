# Makefile for the Nagios Incident.io Notification Plugin

# Name of the binary to generate
BINARY = notify_incident_io

# Go build flags
BUILD_FLAGS = -ldflags "-s -w"

.PHONY: all build clean

# Default target: Build the binary
all: build

# Build the binary
build:
	go build $(BUILD_FLAGS) -o $(BINARY) nagios-incident-io.go

# Clean up build artifacts
clean:
	rm -f $(BINARY)