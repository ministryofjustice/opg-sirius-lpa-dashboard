FROM node:20.5.1-alpine3.17@sha256:c8ae1785da54b802347d06e826d5ae127fc7851114803ee4098f782d46ba40a9 as asset-env

WORKDIR /app

RUN mkdir -p web/static

COPY package.json .
COPY yarn.lock .
RUN yarn

COPY web/assets web/assets
RUN yarn build

FROM golang:1.21.0@sha256:b490ae1f0ece153648dd3c5d25be59a63f966b5f9e1311245c947de4506981aa as build-env

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY main.go main.go
COPY internal internal

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/opg-sirius-lpa-dashboard

FROM alpine:3

WORKDIR /go/bin

RUN apk --update --no-cache add \
    ca-certificates \
    tzdata

# Patch vulnerabilities
RUN apk upgrade --no-cache libcrypto3 libssl3

COPY --from=build-env /go/bin/opg-sirius-lpa-dashboard opg-sirius-lpa-dashboard
COPY --from=asset-env /app/web/static web/static
COPY web/template web/template

ENTRYPOINT ["./opg-sirius-lpa-dashboard"]
