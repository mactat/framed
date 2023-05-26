FROM golang:1.20.4-alpine3.18 as builder

WORKDIR /app
COPY go.mod go.sum .
RUN go mod download
COPY . .
RUN go build -o framed ./main.go

FROM alpine:3.14.2 as test

COPY --from=builder /app/framed /bin/framed
COPY ./framed.yaml /app/framed.yaml

WORKDIR /app
CMD ["framed create", "ls"]