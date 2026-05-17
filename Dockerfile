# Build Stage
FROM golang:1.26-alpine AS build

WORKDIR /src/meme-api

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build a fully static binary with debug info stripped
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/meme-api

# Final Stage — minimal image, no shell needed at runtime
FROM alpine:3.21

# Run as non-root for security
RUN addgroup -S app && adduser -S app -G app
USER app

WORKDIR /app
COPY --from=build /app/meme-api ./

EXPOSE 8080

CMD ["./meme-api"]
