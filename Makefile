VERSION := $(shell git rev-parse --short HEAD)
NAME := $(shell basename $(CURDIR))

all:
	go build -ldflags "-X main.buildVersion=$(VERSION)"

install:
	go install -ldflags "-X main.buildVersion=$(VERSION)"

