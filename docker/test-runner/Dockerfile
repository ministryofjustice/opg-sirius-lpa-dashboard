FROM golang:1.25-bookworm@sha256:9422886b8f9b52e88344a24e9b05fd4b37d42233b680019fc3cb6b1fb2f2b0a5

RUN go install gotest.tools/gotestsum@latest
RUN go install github.com/pact-foundation/pact-go/v2@v2.2.0
## This is needed to install a dependency not included in the main package
RUN pact-go -l DEBUG install
