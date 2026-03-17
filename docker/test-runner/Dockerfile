FROM golang:1.25-bookworm@sha256:1172124903d5f4fbeaa753553652c14215e134aa4bd5788dbb86f76c94103bb4

RUN go install gotest.tools/gotestsum@latest
RUN go install github.com/pact-foundation/pact-go/v2@v2.2.0
## This is needed to install a dependency not included in the main package
RUN pact-go -l DEBUG install
