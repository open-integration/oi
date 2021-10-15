FROM golang:1.16-alpine as build

WORKDIR /app

RUN apk update && apk add make wget git

ARG SERVICE

COPY . .

RUN go mod download

RUN go build -o service cmd/catalog/${SERVICE}/main.go

FROM alpine
RUN apk --no-cache add ca-certificates
COPY --from=build /app/service /app/service

WORKDIR /app

ENV PORT=8080

CMD [ "/app/service" ]