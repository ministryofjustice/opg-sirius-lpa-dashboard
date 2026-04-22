FROM golang:1.25-bookworm@sha256:8df46c6698f5a140170b8806029abb16fe2e381e487f1e120babeaeb0267d2df

RUN go install gotest.tools/gotestsum@latest
RUN go install github.com/pact-foundation/pact-go/v2@v2.2.0
## This is needed to install a dependency not included in the main package
RUN pact-go -l DEBUG install
