FROM golang:1.24.2 AS build

WORKDIR /app

COPY . .

RUN go clean --modcache
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/main.go

# Stage 2: Runtime stage
FROM alpine:latest
RUN apk add --no-cache curl

WORKDIR /root

COPY --from=build /app/main .

EXPOSE 3000

CMD ["./main"]
