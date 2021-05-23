FROM docker as docker

FROM golang:1.16.4-alpine as go

WORKDIR /testing

COPY --from=docker /usr/local/bin/docker /usr/local/bin/docker 

RUN apk update && apk add make wget curl bash python3 git jq

RUN go get github.com/client9/misspell/cmd/misspell && \
    go install github.com/fzipp/gocyclo/cmd/gocyclo@latest && \
    go get github.com/securego/gosec/cmd/gosec && \
    go get github.com/google/addlicense && \
    go get github.com/github/hub

RUN curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sh -s -- -b /usr/local/bin v0.18.1

RUN curl -o codecov.sh https://codecov.io/bash && chmod +x codecov.sh

RUN wget https://github.com/mikefarah/yq/releases/download/v4.9.0/yq_linux_amd64.tar.gz && \
    tar -xvf yq_linux_amd64.tar.gz && \
    mv yq_linux_amd64 /usr/local/bin/yq

ENV PATH $PATH:/root/google-cloud-sdk/bin
RUN curl -sSL https://sdk.cloud.google.com | bash

RUN wget https://github.com/dominikh/go-tools/releases/download/2020.2.4/staticcheck_linux_amd64.tar.gz && \
    tar -xvf staticcheck_linux_amd64.tar.gz && \
    mv ./staticcheck/staticcheck /usr/local/bin

RUN curl -sfL https://install.goreleaser.com/github.com/goreleaser/goreleaser.sh | sh
RUN mv bin/goreleaser /usr/local/bin/goreleaser

RUN  apk add gcc build-base