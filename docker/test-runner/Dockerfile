FROM golang:1.25-bookworm@sha256:298734aec230b5f3e8cee450ce6d7eccc39f1797ba548ee90d57e9803030c6c3

RUN go install gotest.tools/gotestsum@latest
RUN go install github.com/pact-foundation/pact-go/v2@v2.2.0
## This is needed to install a dependency not included in the main package
RUN pact-go -l DEBUG install
