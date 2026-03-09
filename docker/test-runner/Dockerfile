FROM golang:1.25-bookworm@sha256:7af46e70d2017aef0b4ce2422afbcf39af0511a61993103e948b61011233ec42

RUN go install gotest.tools/gotestsum@latest
RUN go install github.com/pact-foundation/pact-go/v2@v2.2.0
## This is needed to install a dependency not included in the main package
RUN pact-go -l DEBUG install
