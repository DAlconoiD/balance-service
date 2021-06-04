FROM golang:1.13-alpine AS build

WORKDIR /go/src/balance_microservice

COPY . .
COPY ./config.yaml /go/bin

RUN go install ./...

FROM alpine:3.12
WORKDIR /usr/bin
COPY --from=build /go/bin .

#docker build . -t balance_srv
#docker run --link pg_balance --rm -p 8081:8081 -d --name balance balance_srv balance-service
#docker kill balance
