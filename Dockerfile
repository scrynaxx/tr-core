FROM golang:1.27-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
RUN go build -o /out ./cmd/

FROM alpine:3.23.3

WORKDIR /usr/local/bin

COPY --from=builder /out main

CMD ["main"]
