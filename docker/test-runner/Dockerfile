FROM golang:1.25-bookworm@sha256:61170847dd0e20c3c1e374aea35c7057c42e1d638d775c360afcd548531c05ab

RUN go install gotest.tools/gotestsum@latest
RUN go install github.com/pact-foundation/pact-go/v2@v2.2.0
## This is needed to install a dependency not included in the main package
RUN pact-go -l DEBUG install
