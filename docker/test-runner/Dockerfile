FROM golang:1.26-bookworm@sha256:252599aeb51ad60b83e4d8821802068127c528c707cb7dd7afd93be057c6011c

RUN go install gotest.tools/gotestsum@latest
RUN go install github.com/pact-foundation/pact-go/v2@v2.4.3
## This is needed to install a dependency not included in the main package
RUN pact-go -l DEBUG install
