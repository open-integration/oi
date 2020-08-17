FROM golang:1.14.0 as go

WORKDIR /testing

RUN apt-get update && apt-get -y upgrade  && apt-get install build-essential make wget -y

RUN go get github.com/client9/misspell/cmd/misspell && \
    go get github.com/fzipp/gocyclo && \
    go get golang.org/x/lint/golint && \
    go get github.com/securego/gosec/cmd/gosec && \
    go get github.com/google/addlicense && \
    go get github.com/github/hub

# Install snyk cli for security testing
RUN curl -o snyk -L https://github.com/snyk/snyk/releases/download/v1.291.0/snyk-linux && chmod +x snyk && mv snyk /usr/local/bin

RUN curl -o codecov.sh https://codecov.io/bash && chmod +x codecov.sh