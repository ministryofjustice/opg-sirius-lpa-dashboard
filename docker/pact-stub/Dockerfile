FROM golang:1.26@sha256:11fd8f7f63db3b6fb198797042ba4c40a4a34dc83325d3328ca3bc4bb7726786

WORKDIR /app

COPY go.mod .
COPY pact_stub.go .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/pact_stub

ENTRYPOINT /go/bin/pact_stub
