FROM golang:1.16-alpine

WORKDIR /app

RUN apk update && apk add make wget git

ARG SHA

# clone the current sha because gorelaser will not pass the full directory
# as context, only the build artefact.
# the extra_files does not support wildcard
RUN git clone https://github.com/open-integration/oi /app && \
    git checkout ${SHA}

RUN go mod download