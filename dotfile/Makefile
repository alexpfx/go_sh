VERS = $(shell git rev-parse HEAD)
DATE = $(shell date)
FLAGS = -X "main.version=${VERS}" -X "main.buildTime=${DATE}"
test_version:
	echo "${REPO_OUTPUT}"
install-cfg:
	go build -o $(GOPATH)/bin/cfg -ldflags '${FLAGS}' cmd/cfg/main.go
install-repo:
	go build -o $(GOPATH)/bin/repocfg  -ldflags '${FLAGS}' cmd/repo/main.go
install: install-cfg install-repo







