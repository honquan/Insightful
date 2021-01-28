PROJECT ?= "insightful"
PROJECT_NAME := $(PROJECT)

.PHONY: all prepare build

all: build

prepare:
	@go mod download

build:
	@go build -o ./$(PROJECT_NAME)
