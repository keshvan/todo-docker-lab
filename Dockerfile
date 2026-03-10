FROM golang:alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN  CGO_ENABLED=0 GOOS=linux go build -o todo-api ./cmd/todo
RUN CGO_ENABLED=0 GOOS=linux go build -o migrate ./cmd/migrate

FROM alpine:latest
WORKDIR /app
COPY --from=builder /build/todo-api .
COPY --from=builder /build/migrate .
COPY --from=builder /build/migrations ./migrations
CMD ["./todo-api"]