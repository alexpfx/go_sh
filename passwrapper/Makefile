VERS = $(shell git rev-parse HEAD)
DATE = $(shell date)
FLAGS = -X "main.version=${VERS}" -X "main.buildTime=${DATE}" -w -s
test_version:
	echo "${REPO_OUTPUT}"
install_passwrapper:
	go build -o $(GOPATH)/bin/pass-wrapper -ldflags '${FLAGS}' cmd/main.go
install: install_passwrapper

