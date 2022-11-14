FROM node:17.3.0-alpine3.12 as asset-env

WORKDIR /app

RUN mkdir -p web/static

COPY package.json .
COPY yarn.lock .
RUN yarn

COPY web/assets web/assets
RUN yarn build

FROM golang:1.19.3 as build-env

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY main.go main.go
COPY internal internal

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/opg-sirius-lpa-dashboard

FROM alpine:3.16.3

WORKDIR /go/bin

RUN apk --update --no-cache add \
    ca-certificates \
    && rm -rf /var/cache/apk/*
RUN apk --no-cache add tzdata

COPY --from=build-env /go/bin/opg-sirius-lpa-dashboard opg-sirius-lpa-dashboard
COPY --from=asset-env /app/web/static web/static
COPY web/template web/template

ENTRYPOINT ["./opg-sirius-lpa-dashboard"]
