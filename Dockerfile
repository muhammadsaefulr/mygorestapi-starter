# Stage 1: Build Go binary
FROM golang:1.24.2 AS build

WORKDIR /app

COPY . .

RUN go clean --modcache
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/main.go

# Stage 2: Runtime (dengan Google Chrome)
FROM debian:bullseye-slim

# Install dependencies untuk Google Chrome
RUN apt-get update && apt-get install -y \
    wget gnupg ca-certificates fonts-liberation \
    libappindicator3-1 libasound2 libatk-bridge2.0-0 libatk1.0-0 \
    libcups2 libdbus-1-3 libgdk-pixbuf2.0-0 libnspr4 libnss3 \
    libx11-xcb1 libxcomposite1 libxdamage1 libxrandr2 xdg-utils \
    curl unzip --no-install-recommends && \
    rm -rf /var/lib/apt/lists/*

# Install Google Chrome
RUN wget -O /tmp/chrome.deb https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb && \
    apt-get update && apt-get install -y /tmp/chrome.deb || apt-get -fy install && \
    rm /tmp/chrome.deb && \
    ln -s /usr/bin/google-chrome-stable /usr/bin/google-chrome

WORKDIR /root

COPY --from=build /app/main .

EXPOSE 3000

CMD ["./main"]
