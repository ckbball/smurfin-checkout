FROM golang:alpine as builder

RUN apk update && apk upgrade && \
    apk add --no-cache git

RUN mkdir /app
WORKDIR /app

ENV GO111MODULE=on

COPY . .

RUN cd cmd/server && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o smurfin-checkout && ls && pwd

FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN mkdir /app
WORKDIR /app
COPY --from=builder /app/cmd/server .

CMD ["./smurfin-checkout -grpc-port=9091 -http-port=8080 -db-host=blah -db-user=dev -db-password=dev-user5 -db-schema=checkout -catalog-service-address=localhost:9090"]