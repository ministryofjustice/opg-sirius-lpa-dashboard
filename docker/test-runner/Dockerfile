FROM golang:1.22-bookworm

RUN go install gotest.tools/gotestsum@latest
RUN go install github.com/pact-foundation/pact-go/v2@v2.0.4
## This is needed to install a dependency not included in the main package
RUN pact-go -l DEBUG install

WORKDIR /opt
ENV PACT_CLI_VERSION=v2.4.2
RUN curl -fsSL https://raw.githubusercontent.com/pact-foundation/pact-ruby-standalone/master/install.sh > ./install.sh && bash ./install.sh
ENV PATH="$PATH:/opt/pact/bin"
