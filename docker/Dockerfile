FROM golang:1.25.1-alpine AS builder
WORKDIR /app
# after add go.sum
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/api ./cmd/api

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

FROM scratch
COPY --from=builder /app/bin/api /api
EXPOSE 8080
ENTRYPOINT ["/api"]
