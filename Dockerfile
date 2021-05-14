FROM golang:1.16.0 as go

WORKDIR /app

COPY . .