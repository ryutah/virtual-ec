FROM ryutah/develop-environment

RUN apt-get update \
 && apt-get install -y build-essential \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/*

ENV GO_VERSION=1.14.6 \
    GOLANGCI_LINT=1.30.0

# Install Golang
RUN wget -q https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz \
 && tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz \
 && echo 'export PATH=${PATH}:/usr/local/go/bin:/root/go/bin' >> .bashrc \
 && rm -f go${GO_VERSION}.linux-amd64.tar.gz

# Install Golang tools
RUN GO111MODULE=on /usr/local/go/bin/go get -u \
      golang.org/x/tools/cmd/goimports \
      github.com/rjeczalik/interfaces/cmd/interfacer \
      github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.3.11
