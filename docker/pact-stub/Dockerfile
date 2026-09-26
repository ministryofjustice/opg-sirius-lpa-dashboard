FROM golang:1.27@sha256:3680233e3204827fbdc66088528ae6d4b3d034f51d03a99d454f6de034888244

WORKDIR /app

COPY go.mod .
COPY pact_stub.go .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/pact_stub

ENTRYPOINT /go/bin/pact_stub
