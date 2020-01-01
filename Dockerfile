FROM golang:alpine as builder

RUN apk update && apk upgrade && \
    apk add --no-cache git

RUN mkdir /app
WORKDIR /app

ENV GO111MODULE=on

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN cd cmd/server && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o smurfin-checkout . && ls && pwd

FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN mkdir /app
WORKDIR /app
COPY --from=builder /app/cmd/server .

CMD ["./smurfin-checkout"]