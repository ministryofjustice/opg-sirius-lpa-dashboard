FROM golang:1.25-bookworm@sha256:aacb1bc47244b9e67128d9e8df5c44eb4dc8064f5c2c12fa6bbc590ee0bd9e53

RUN go install gotest.tools/gotestsum@latest
RUN go install github.com/pact-foundation/pact-go/v2@v2.2.0
## This is needed to install a dependency not included in the main package
RUN pact-go -l DEBUG install
