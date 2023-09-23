# Get environment variables
include .env

# Name of app
APP_NAME=bintransfer

# Build folder
BUILD_DIR=build

# Build the app
build:
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd

# Run the app
run:
	$(BUILD_DIR)/$(APP_NAME) \
		--path=$(PATH) \
		--out=$(OUTDIR) \
	

.PHONY: build run