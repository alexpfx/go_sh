FROM golang:1.18-rc as builder
WORKDIR /go/src/go_sh/passwrapper
COPY . .

RUN GOOS=linux go build -o ./passwrapper -ldflags="-s -w" cmd/main.go

FROM alpine:3.15
WORKDIR /go/src/go_sh/passwrapper
COPY --from=builder /go/src/go_sh/passwrapper/passwrapper .
# CMD ["sh", "-c", "./passwrapper ${params}" ]
ENTRYPOINT ["./passwrapper" ]
