export GO111MODULE=on

SHELL := /bin/bash

AGENT_DIR = cmd/agent
AGENT_NAME = stzagent
AGENT_CODE = ${AGENT_DIR:=/*.go}
AGENT_LNX32 = ${AGENT_NAME:=-linux32}
AGENT_LNX64 = ${AGENT_NAME:=-linux64}
AGENT_OSX = ${AGENT_NAME:=-osx}
AGENT_WIN32 = ${AGENT_NAME:=-win32.exe}
AGENT_WIN64 = ${AGENT_NAME:=-win64.exe}
AGENT_BSD32 = ${AGENT_NAME:=-freebsd32}
AGENT_BSD64 = ${AGENT_NAME:=-freebsd64}

ADMINSERVER_DIR = cmd/admin
ADMINSERVER_NAME = stzadmin
ADMINSERVER_CODE = ${ADMINSERVER_DIR:=/*.go}

HTTPSERVER_DIR = cmd/httpserver
HTTPSERVER_NAME = stzhttp
HTTPSERVER_CODE = ${HTTPSERVER_DIR:=/*.go}

TCPSERVER_DIR = cmd/tcpserver
TCPSERVER_NAME = stztcp
TCPSERVER_CODE = ${TCPSERVER_DIR:=/*.go}

UDPSERVER_DIR = cmd/udpserver
UDPSERVER_NAME = stzudp
UDPSERVER_CODE = ${UDPSERVER_DIR:=/*.go}

DEST ?= /opt/stanza

OUTPUT = bin

STATIC_ARGS = -ldflags "-linkmode external -extldflags -static"

.PHONY: all build client agents http clean

all: build

# Build code according to caller OS and architecture
build:
	make agent
	make http
	make tcp
	make udp
	make admin

# Build just the agent/client in the current platform, for testing
agent:
	go build -o $(OUTPUT)/$(AGENT_NAME) $(AGENT_CODE)

# Build admin server
admin:
	go build -o $(OUTPUT)/$(ADMINSERVER_NAME) $(ADMINSERVER_CODE)

# Build HTTP server
http:
	go build -o $(OUTPUT)/$(HTTPSERVER_NAME) $(HTTPSERVER_CODE)

# Build TCP server
tcp:
	go build -o $(OUTPUT)/$(TCPSERVER_NAME) $(TCPSERVER_CODE)

# Build UDP server
udp:
	go build -o $(OUTPUT)/$(UDPSERVER_NAME) $(UDPSERVER_CODE)

# Build all agents
agents:
	make agent_osx_intel
	make agent_osx_arm
	make agent_freebsd64
	make agent_freebsd32
	make agent_linux64
	make agent_linux32
	make agent_windows64
	make agent_windows32

# Build the agent for OSX Intel 64-bit
agent_osx_intel:
	GOOS=darwin GOARCH=amd64 go build -o $(OUTPUT)/$(AGENT_OSX) $(AGENT_CODE)

# Build the agent for OSX ARM 64-bit
agent_osx_arm:
	GOOS=darwin GOARCH=arm64 go build -o $(OUTPUT)/$(AGENT_OSX) $(AGENT_CODE)

# Build the agent for FreeBSD 64-bit
agent_freebsd64:
	GOOS=freebsd GOARCH=amd64 go build -o $(OUTPUT)/$(AGENT_BSD64) $(AGENT_CODE)

# Build the agent for FreeBSD 32-bit
agent_freebsd32:
	GOOS=freebsd GOARCH=386 go build -o $(OUTPUT)/$(AGENT_BSD32) $(AGENT_CODE)

# Build the agent for Linux 64-bit
agent_linux64:
	GOOS=linux GOARCH=amd64 go build -o $(OUTPUT)/$(AGENT_LNX64) $(AGENT_CODE)

# Build the agent for Linux 32-bit
agent_linux32:
	GOOS=linux GOARCH=386 go build -o $(OUTPUT)/$(AGENT_LNX32) $(AGENT_CODE)

# Build the agent for Windows 64-bit
agent_windows64:
	GOOS=windows GOARCH=amd64 go build -o $(OUTPUT)/$(AGENT_WIN64) $(AGENT_CODE)

# Build the agent for Windows 32-bit
agent_windows32:
	GOOS=windows GOARCH=386 go build -o $(OUTPUT)/$(AGENT_WIN32) $(AGENT_CODE)

# Build dev docker containers and run them (also generates new certificates)
docker_dev_build:
ifeq (,$(wildcard ./.env))
	$(error Missing .env file)
endif
	docker-compose -f docker-compose-dev.yml build

# Build and run dev docker containers
make docker_dev:
	make docker_dev_build
	make docker_dev_up

# Run docker containers
docker_dev_up:
	docker-compose -f docker-compose-dev.yml up

# Takes down docker containers
docker_dev_down:
	docker-compose -f docker-compose-dev.yml down

# Deletes all osctrl docker images
docker_dev_clean:
	docker images | grep stanza | awk '{print $$3}' | xargs -rI {} docker rmi -f {}

# Rebuilds the admin server
docker_dev_rebuild_admin:
	docker-compose -f docker-compose-dev.yml up --force-recreate --no-deps -d --build stanza-admin

# Rebuilds the http server
docker_dev_rebuild_http:
	docker-compose -f docker-compose-dev.yml up --force-recreate --no-deps -d --build stanza-http

# Rebuilds the agent container
docker_dev_rebuild_agent:
	docker-compose -f docker-compose-dev.yml up --force-recreate --no-deps -d --build stanza-agent
