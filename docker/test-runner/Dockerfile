FROM golang:1.25-bookworm@sha256:7fb09d8804035fbde8a84ed59ca9f46dd68c6f160f9d193e98d795d8d9e002ec

RUN go install gotest.tools/gotestsum@latest
RUN go install github.com/pact-foundation/pact-go/v2@v2.2.0
## This is needed to install a dependency not included in the main package
RUN pact-go -l DEBUG install
