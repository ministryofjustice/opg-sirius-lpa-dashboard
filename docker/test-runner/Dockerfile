FROM golang:1.25-bookworm@sha256:35bfe6cf367b9e6756007dcbf388c0a74498d6af54f5bbe0f95cd7f2daf519d7

RUN go install gotest.tools/gotestsum@latest
RUN go install github.com/pact-foundation/pact-go/v2@v2.2.0
## This is needed to install a dependency not included in the main package
RUN pact-go -l DEBUG install
