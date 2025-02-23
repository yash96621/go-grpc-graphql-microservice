FROM golang:1.16.3-alpine3.13 AS builder
RUN apk add --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/github.com/yash96621/go-grpc-graphql-microservice
COPY go.mod go.sum ./
COPY vendor vendor
COPY catalog catalog
RUN GO111MODULE=on go build -mod verndor  -o /go/bin/app ./catalog/cmd/account

FROM alpine:3.13
WORKDIR /usr/bin
COPY --from=build /go/bin .
EXPOSE 8080
CMD ["app"]