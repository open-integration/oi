FROM golang:1.16.0 as go

WORKDIR /testing

RUN apt-get update && apt-get -y upgrade  && apt-get install build-essential make wget -y

RUN go get github.com/client9/misspell/cmd/misspell && \
    go install github.com/fzipp/gocyclo/cmd/gocyclo@latest && \
    go get golang.org/x/lint/golint && \
    go get github.com/securego/gosec/cmd/gosec && \
    go get github.com/google/addlicense && \
    go get github.com/github/hub

RUN curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sh -s -- -b /usr/local/bin v0.18.1

RUN curl -o codecov.sh https://codecov.io/bash && chmod +x codecov.sh

RUN wget https://github.com/mikefarah/yq/releases/download/v4.9.0/yq_linux_amd64.tar.gz && \
    tar -xvf yq_linux_amd64.tar.gz && \
    mv yq_linux_amd64 /usr/local/bin/yq