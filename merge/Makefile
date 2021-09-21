VERS = $(shell git rev-parse HEAD)
DATE = $(shell date)
FLAGS = -X "main.version=${VERS}" -X "main.buildTime=${DATE}" -w -s
test_version:
	echo "${REPO_OUTPUT}"
merge:
	go build -o $(GOPATH)/bin/merge -ldflags '${FLAGS}' cmd/main.go
install: merge

