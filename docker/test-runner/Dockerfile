FROM golang:1.25-bookworm@sha256:0515b52dd3015b4bc3dfe20e573393dd0846bb3c6c3810b5642da793a236b44a

RUN go install gotest.tools/gotestsum@latest
RUN go install github.com/pact-foundation/pact-go/v2@v2.2.0
## This is needed to install a dependency not included in the main package
RUN pact-go -l DEBUG install
