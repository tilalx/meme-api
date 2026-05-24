# Build Stage
FROM golang:1.26-alpine AS build

WORKDIR /src/meme-api

COPY go.mod go.sum ./
RUN go mod download

# Install swag CLI for doc generation
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.4

COPY . .

# Generate Swagger docs from handler annotations
RUN swag init --parseDependency --parseInternal

# Build a fully static binary with debug info stripped
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/meme-api

# Final Stage — minimal image, no shell needed at runtime
FROM alpine:3.23

# Run as non-root for security
RUN addgroup -S app && adduser -S app -G app
USER app

WORKDIR /app
COPY --from=build /app/meme-api ./

EXPOSE 8080

CMD ["./meme-api"]
