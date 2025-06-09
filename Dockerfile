# Build stage
FROM golang:latest AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS="linux" GOARCH="amd64" go build -o main cmd/main.go

# Final stage
FROM gcr.io/distroless/static-debian12
WORKDIR /src
COPY --from=builder /app/main .
COPY --from=builder /app/.env .
CMD ["/src/main"]