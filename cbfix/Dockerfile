FROM golang:1.18-rc-alpine as builder
WORKDIR /go/src/go_sh
COPY . .
RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/app -ldflags="-s -w" cmd/main.go

FROM scratch as app
WORKDIR /go/src/go_sh
COPY --from=builder /go/bin/app /go/bin/app
ENTRYPOINT [ "/go/bin/app" ]