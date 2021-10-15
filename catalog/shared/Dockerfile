FROM openintegration/oi:base as build

ARG SERVICE

COPY . .

RUN go build -o service cmd/catalog/${SERVICE}/main.go

FROM alpine
RUN apk --no-cache add ca-certificates
COPY --from=build /app/service /app/service

WORKDIR /app

ENV PORT=8080

CMD [ "/app/service" ]