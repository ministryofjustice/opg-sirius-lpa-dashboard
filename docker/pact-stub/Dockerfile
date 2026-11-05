FROM golang:1.26.10@sha256:672c944cfd07599b5f1bb179869270cdc94c6254e41db4d431259693cb286b82

WORKDIR /app

COPY go.mod .
COPY pact_stub.go .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/pact_stub

ENTRYPOINT /go/bin/pact_stub
