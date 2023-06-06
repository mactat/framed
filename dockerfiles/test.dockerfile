ARG GO_VERSION=1.20.4
ARG ALPINE_VERSION=3.18

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o framed ./framed.go

FROM bats/bats:latest as tester

WORKDIR /framed
COPY --from=builder /app/framed /bin/framed
COPY tests/ /tests

RUN LANG=go ./test/bats/bin/bats /tests/test.bats
