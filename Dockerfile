FROM golang:1.16.4-alpine as go

WORKDIR /app

COPY . .