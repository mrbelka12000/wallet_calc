## Build
FROM golang:1.24.4-alpine3.22 AS buildenv

ADD go.mod go.sum /

RUN go mod download

WORKDIR /app

ADD . .

RUN  go build -o main cmd/main.go

## Deploy
FROM alpine

RUN apk add --no-cache poppler-utils

WORKDIR /

COPY --from=buildenv  /app/ /

EXPOSE 8084

CMD ["/main"]