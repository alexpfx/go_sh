VERS = $(shell git rev-parse HEAD)
DATE = $(shell date)
FLAGS = -X "main.version=${VERS}" -X "main.buildTime=${DATE}" -w -s
test_version:
	echo "${REPO_OUTPUT}"
install_elapsed:
	go build -o $(GOPATH)/bin/elapsed -ldflags '${FLAGS}' cmd/main.go
install: install_elapsed

