# syntax=docker/dockerfile:1

FROM golang:1.24 AS build
# This needs CGO
RUN <<EOF
go install github.com/k3s-io/kine@v0.13.10
EOF

WORKDIR /go/src/github.com/lxfontes/standalone-api-server

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY . .
RUN go install

FROM ubuntu:jammy AS base
RUN <<EOF
apt-get update
apt-get install -y ca-certificates
EOF
COPY --from=build /go/bin/kine /usr/bin/kine
COPY --from=build /go/bin/standalone-api-server /usr/bin/standalone-api-server
