FROM node:22.12.0-alpine3.19 AS asset-env

WORKDIR /app

RUN mkdir -p web/static

COPY package.json .
COPY package-lock.json .
RUN npm ci --ignore-scripts

COPY web/assets web/assets
RUN npm run build

FROM golang:1.25 AS build-env

ARG TARGETARCH
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY main.go main.go
COPY internal internal

RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -a -installsuffix cgo -o /go/bin/opg-sirius-lpa-dashboard

FROM build-env AS healthcheck-build
WORKDIR /app

COPY healthcheck healthcheck

WORKDIR /app/healthcheck
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/healthcheck

FROM scratch

WORKDIR /go/bin

COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build-env /usr/share/zoneinfo /usr/share/zoneinfo

COPY --from=build-env /go/bin/opg-sirius-lpa-dashboard opg-sirius-lpa-dashboard
COPY --from=healthcheck-build /go/bin/healthcheck healthcheck
COPY --from=asset-env /app/web/static web/static
COPY web/template web/template

USER 65532

HEALTHCHECK --interval=5s --timeout=5s --start-period=5s --retries=3 CMD [ "/go/bin/healthcheck" ]

ENTRYPOINT ["./opg-sirius-lpa-dashboard"]
