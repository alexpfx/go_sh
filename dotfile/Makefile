VERS = $(shell git rev-parse HEAD)
DATE = $(shell date)
FLAGS = -X "main.version=${VERS}" -X "main.buildTime=${DATE}"


test_version:
	echo "${FLAGS}"
build:
	go build -ldflags '${FLAGS}' -o bin/initrepo cmd/repo/main.go


