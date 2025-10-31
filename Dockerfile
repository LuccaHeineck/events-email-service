# ---------- Stage 1: Build ----------
FROM golang:1.25.3 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o events-email-service .

# ---------- Stage 2: Run ----------
FROM alpine:3.20

WORKDIR /app
COPY --from=builder /app/events-email-service .

# The app loads environment variables at runtime
EXPOSE 8082
ENTRYPOINT ["./events-email-service"]
