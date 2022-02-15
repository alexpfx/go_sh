FROM golang:1.18-rc as builder
WORKDIR /go/src/go_sh
COPY . .

RUN GOOS=linux go build -o ./app -ldflags="-s -w" cmd/main.go

FROM alpine:3.15
WORKDIR /go/src/go_sh
COPY --from=builder /go/src/go_sh .
ENTRYPOINT ["./app" ]
#