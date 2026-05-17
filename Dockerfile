# Build Stage
FROM golang:1.25-alpine as build

WORKDIR /src/meme-api

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o /app/meme-api

# Final Stage
FROM alpine:latest
WORKDIR /app
COPY --from=build /app/meme-api /app/

EXPOSE 8080

CMD ./meme-api
