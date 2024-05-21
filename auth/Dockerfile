FROM golang:latest AS builder

WORKDIR /app

COPY go.mod ./

RUN go mod tidy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o auth ./cmd

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/auth .

ENTRYPOINT ["./auth"]