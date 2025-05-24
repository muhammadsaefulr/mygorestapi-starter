FROM golang:1.22 AS build

WORKDIR /app

COPY . .

RUN go clean --modcache
RUN go mod tidy

# Menyusun aplikasi Go dari file main.go yang ada di dalam direktori cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/main.go

# Stage 2: Runtime stage
FROM alpine:latest

# Menginstal dependensi yang diperlukan di runtime (misalnya curl)
RUN apk add --no-cache curl

WORKDIR /root

# Menyalin hasil build dari tahap pertama (file main dan .env)
COPY --from=build /app/main .
COPY --from=build /app/.env .

# Menentukan port yang akan diekspos
EXPOSE 3000

# Menjalankan aplikasi
CMD ["./main"]
