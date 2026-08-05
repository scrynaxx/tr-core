FROM golang:1.26.5-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
RUN go build -o /out/app ./services/order/cmd/

FROM alpine:3.23.3

WORKDIR /usr/local/bin

COPY --from=builder /out/app app
COPY ./config ./config

CMD ["app"]
