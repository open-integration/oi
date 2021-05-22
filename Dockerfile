FROM golang:1.16.4-alpine

WORKDIR /app

RUN apk update && apk add make wget 

COPY . .

RUN go mod download