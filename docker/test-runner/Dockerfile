FROM golang:1.26-bookworm@sha256:785a3f7d17276de845996f27f0231da3f61e8cc05ce0dbdaac757f3b4b937ac7

RUN go install gotest.tools/gotestsum@latest
RUN go install github.com/pact-foundation/pact-go/v2@v2.2.0
## This is needed to install a dependency not included in the main package
RUN pact-go -l DEBUG install
