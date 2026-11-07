FROM golang:1.27.2@sha256:820153096fda4de6ac1b7194d12e416b7dd8f4e9a33d17e0345b111b945162d8

WORKDIR /app

COPY go.mod .
COPY pact_stub.go .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/pact_stub

ENTRYPOINT /go/bin/pact_stub
