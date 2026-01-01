FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY ./backend .
RUN go mod download
RUN go build -o arb-backend

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/arb-backend .
CMD ["./arb-backend"]
