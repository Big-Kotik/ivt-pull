# syntax=docker/dockerfile:1

FROM golang:1.19.1-alpine

WORKDIR /app

COPY go.mod go.sum ./
COPY ./ivt-pull-api ./
RUN go mod download
RUN go mod tidy

COPY . ./

RUN go build -o /pull-server ./cmd/pull-server

ENTRYPOINT [ "/pull-server" ]