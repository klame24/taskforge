FROM golang:1.25.1-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/api ./cmd/api

FROM scratch
COPY --from=builder /app/bin/api /api
EXPOSE 8080
ENTRYPOINT ["/api"]
