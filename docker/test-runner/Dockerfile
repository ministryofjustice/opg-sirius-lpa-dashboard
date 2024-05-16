FROM golang:1.22-bookworm

RUN go install gotest.tools/gotestsum@latest
RUN go install github.com/pact-foundation/pact-go/v2@v2.0.4
## This is needed to install a dependency not included in the main package
RUN pact-go -l DEBUG install

WORKDIR /opt
RUN curl -fsSL https://raw.githubusercontent.com/pact-foundation/pact-ruby-standalone/master/install.sh > ./install.sh && bash ./install.sh tag=v1.92.0
ENV PATH="$PATH:/opt/pact/bin"