FROM golang:1.24.2 AS build

WORKDIR /app
COPY . .

RUN go clean --modcache
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/main.go

FROM debian:bullseye-slim

WORKDIR /root

COPY --from=build /app/main .

EXPOSE 3000

CMD ["./main"]
