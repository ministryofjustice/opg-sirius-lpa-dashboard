FROM golang:1.25-bookworm@sha256:fd0db0c7924675845d82cd0f45bd98376426acf09148392a085dec4b5200ed34

RUN go install gotest.tools/gotestsum@latest
RUN go install github.com/pact-foundation/pact-go/v2@v2.2.0
## This is needed to install a dependency not included in the main package
RUN pact-go -l DEBUG install
