VERS = $(shell git rev-parse HEAD)
DATE = $(shell date)
FLAGS = -X "main.version=${VERS}" -X "main.buildTime=${DATE}" -w -s
test_version:
	echo "${REPO_OUTPUT}"
install_cbfix:
	go build -o $(GOPATH)/bin/cbfix -ldflags '${FLAGS}' cmd/main.go
install: install_cbfix

