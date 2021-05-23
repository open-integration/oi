FROM openintegration/oi

ARG SERVICE

RUN go build -o service cmd/catalog/${SERVICE}/main.go 

ENV PORT=8080

CMD [ "/app/service" ]