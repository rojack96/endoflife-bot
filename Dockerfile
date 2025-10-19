FROM golang:1.25.3 AS builder

ARG VERSION=v0.1.0
WORKDIR /app

COPY . .
COPY .env .
# Compila per Linux/ARM64 (compatibile col tuo server)
RUN GOOS=linux GOARCH=arm64 go build -o endoflife-${VERSION} .

# --- Runtime stage ---
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

ARG VERSION=v0.1.0
WORKDIR /app

COPY --from=builder /app/endoflife-${VERSION} .
RUN chmod +x endoflife-${VERSION} && \
    useradd -m app && chown app:app endoflife-${VERSION}

USER app
CMD ./endoflife-${VERSION}